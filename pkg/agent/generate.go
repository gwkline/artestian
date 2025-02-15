package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/types"
)

func (p *AnthropicProvider) GenerateTest(params types.GenerateTestParams) (string, error) {
	slog.Debug("preparing test generation prompt",
		"exampleName", params.Example.Name,
		"exampleType", params.Example.Type,
		"sourceCodeLength", len(params.SourceCode))

	// Format context files section
	var contextSection strings.Builder
	for _, cf := range params.ContextFiles {
		contextSection.WriteString(fmt.Sprintf("\n=== %s (%s) ===\n%s\n", cf.Description, cf.Type, cf.Content))
	}

	prompt := fmt.Sprintf(`Generate a similar test for the source code below. The test should:
1. Follow the same/similar patterns as the example
2. You should almost never use mocks, ideally only external API calls should be mocked
3. Use %s and %s to write the test
4. Include all necessary imports, the current working directory is %s
5. Follow best practices for testing

Return ONLY the test code, no explanations.

You are specifically generating a test for the following function:

%s

Which is part of this source code file:

%s

And this example test as reference:

%s

And these context files:

%s`,
		params.Language.GetName(),
		params.TestRunner.GetName(),
		params.Function.Name,
		params.TestDir,
		params.SourceCode,
		params.Example.Content,
		contextSection.String())

	slog.Info("sending request to Anthropic API",
		"promptId", "generate_test",
		"model", anthropic.ModelClaude3_5SonnetLatest,
		"maxTokens", MAX_TOKENS)

	msg, err := p.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(MAX_TOKENS)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})

	if err != nil {
		slog.Error("failed to generate test with Anthropic API", "error", err)
		return "", fmt.Errorf("failed to generate test: %w", err)
	}

	response := removeBackticks(msg.Content[0].Text)

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
