package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/pkg/prompt_utils"
	"github.com/gwkline/artestian/types"
)

func (p *AnthropicProvider) GenerateTest(params types.GenerateTestParams) (string, error) {
	slog.Debug("preparing test generation prompt",
		"exampleName", params.Example.Name,
		"exampleType", params.Example.Type,
		"sourceCodeLength", len(params.SourceCode))

	xmlParams, err := prompt_utils.StructToXMLString(params)
	if err != nil {
		return "", fmt.Errorf("failed to format params: %w", err)
	}

	prompt := fmt.Sprintf(`Generate a test for the function provided. 
The test should:
	1. Follow the same/similar patterns as the example
	2. Focus on testing the core functionality and happy path. Don't waste time testing unlikely edge cases or invalid inputs
	3. You should basically never use mocks, except for external API calls
	4. Use the supplied language and test runner to write the test
	5. Include all necessary imports, the current working directory is %s

Return ONLY the test code, no explanations.

%s`, params.TestPath, xmlParams)
	slog.Info("completion started",
		"promptId", "generate_test",
		"model", anthropic.ModelClaude3_5SonnetLatest,
		"maxTokens", MAX_TOKENS)

	msg, err := p.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(MAX_TOKENS)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
			anthropic.NewAssistantMessage(anthropic.NewTextBlock(fmt.Sprintf("```%s", params.Language.GetName()))),
		}),
	})

	if err != nil {
		if err := p.logger.Log("generate_test", prompt, ""); err != nil {
			slog.Warn("failed to log prompt", "error", err)
		}
		slog.Error("failed to generate test with Anthropic API", "error", err)
		return "", fmt.Errorf("failed to generate test: %w", err)
	}

	response := removeBackticks(fmt.Sprintf("```%s%s", params.Language.GetName(), msg.Content[0].Text))

	// Log the prompt and response
	if err := p.logger.Log("generate_test", prompt, response); err != nil {
		slog.Warn("failed to log prompt", "error", err)
	}

	slog.Debug("received response from Anthropic API",
		"responseLength", len(response))

	return response, nil
}

func removeBackticks(text string) string {
	// Check if text starts with code fence
	if strings.HasPrefix(text, "```") {
		// Find the end of the first line to isolate potential language tag
		firstLineEnd := strings.Index(text, "\n")
		if firstLineEnd != -1 {
			// Remove opening fence with any language tag
			text = text[firstLineEnd+1:]
			// Remove closing fence
			text = strings.TrimSuffix(text, "```")
			text = strings.TrimSpace(text)
		}
	}
	return text
}
