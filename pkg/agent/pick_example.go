package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/types"
)

const assistantMessage = `{"exampleIndex":`

func (p *AnthropicProvider) PickExample(sourceCode string, testExamples []types.TestExample) (types.TestExample, error) {
	if len(testExamples) == 0 {
		return types.TestExample{}, fmt.Errorf("no test examples provided")
	}

	type TruncatedExample struct {
		Name        string
		Type        types.TestType
		Description string
	}

	truncatedExamples := make([]TruncatedExample, len(testExamples))
	for i, example := range testExamples {
		truncatedExamples[i] = TruncatedExample{
			Name:        example.Name,
			Type:        example.Type,
			Description: example.Description,
		}
	}

	prompt := fmt.Sprintf(`Given this source code:

%s

And these test examples:

%v

Which test example would be the best match for testing this code? Consider:
1. The complexity and structure of the code
2. The testing patterns demonstrated in each example
3. The similarity between the example and what needs to be tested

Return a JSON object with the key "exampleIndex" and the value being the index number of the best matching example. 
Do not include any other text or explanations in your response.`, sourceCode, truncatedExamples)

	slog.Info("completion started",
		"promptId", "pick_example",
		"model", anthropic.ModelClaude3_5SonnetLatest,
		"maxTokens", MAX_TOKENS)

	msg, err := p.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(MAX_TOKENS)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
			anthropic.NewAssistantMessage(anthropic.NewTextBlock(assistantMessage)),
		}),
	})

	if err != nil {
		slog.Error("failed to pick example with Anthropic API", "error", err)
		return types.TestExample{}, fmt.Errorf("failed to pick example: %w", err)
	}

	response := msg.Content[0].Text

	selectedIndex, err := parseExampleIndex(response)
	if err != nil {
		return types.TestExample{}, fmt.Errorf("failed to parse example index: %w", err)
	}

	if err := p.logger.Log("pick_example", prompt, msg.Content[0].Text); err != nil {
		slog.Warn("failed to log prompt", "error", err)
	}

	slog.Info("selected test example", "name", testExamples[selectedIndex].Name, "type", testExamples[selectedIndex].Type)

	return testExamples[selectedIndex], nil
}

func parseExampleIndex(response string) (int, error) {
	// Try to parse as a simple JSON object first
	var jsonResponse struct {
		ExampleIndex int `json:"exampleIndex"`
	}
	if err := json.Unmarshal([]byte(response), &jsonResponse); err == nil {
		return jsonResponse.ExampleIndex, nil
	}

	// If not JSON, try to extract the first number from the string
	var selectedIndex int
	_, err := fmt.Sscanf(response, "%d", &selectedIndex)
	if err != nil {
		return 0, fmt.Errorf("invalid example index returned")
	}

	return selectedIndex, nil
}
