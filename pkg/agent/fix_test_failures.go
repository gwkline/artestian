package agent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/pkg/prompt_utils"
	"github.com/gwkline/artestian/types"
)

func (p *AnthropicProvider) FixTestFailures(params types.IterateTestParams) (string, error) {
	xmlString, err := prompt_utils.StructToXMLString(params)
	if err != nil {
		return "", fmt.Errorf("failed to format params: %w", err)
	}

	prompt := fmt.Sprintf(`Fix the test failures in this code.

Here are some reminders:
- Use the conventions and types of the language you're writing the test in
- Use the provided context and examples to help you fix the failures

%s

Return ONLY the fixed test code, no explanations.`, xmlString)

	slog.Info("completion started",
		"promptId", "fix_test_failures",
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
		if err := p.logger.Log("fix_test_failures", prompt, ""); err != nil {
			slog.Warn("failed to log prompt", "error", err)
		}
		slog.Error("failed to fix test errors with Anthropic API", "error", err)
		return "", fmt.Errorf("failed to fix test errors: %w", err)
	}

	response := removeBackticks(fmt.Sprintf("```%s%s", params.Language.GetName(), msg.Content[0].Text))

	if err := p.logger.Log("fix_test_errors", prompt, response); err != nil {
		slog.Warn("failed to log prompt", "error", err)
	}

	slog.Debug("received response from Anthropic API",
		"responseLength", len(response))

	return response, nil
}
