package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/types"
)

func (p *AnthropicProvider) FixTestErrors(params types.IterateTestParams) (string, error) {
	// Format context files section
	var contextSection strings.Builder
	for _, cf := range params.ContextFiles {
		contextSection.WriteString(fmt.Sprintf("\n=== %s (%s) ===\n%s\n", cf.Description, cf.Type, cf.Content))
	}

	prompt := fmt.Sprintf(`Fix the test errors in this code. The errors are:

%s

Here's the current test code:

%s

Here are the relevant context files:

%s

Return ONLY the fixed test code, no explanations.`, params.Errors, params.TestCode, contextSection.String())

	slog.Info("sending request to Anthropic API",
		"promptId", "fix_test_errors",
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
		slog.Error("failed to fix test errors with Anthropic API", "error", err)
		return "", fmt.Errorf("failed to fix test errors: %w", err)
	}

	response := removeBackticks(msg.Content[0].Text)

	if err := p.logger.Log("fix_test_errors", prompt, response); err != nil {
		slog.Warn("failed to log prompt", "error", err)
	}

	slog.Debug("received response from Anthropic API",
		"responseLength", len(response))

	return response, nil
}
