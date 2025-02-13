package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/types"
)

func (p *AnthropicProvider) FixTypeErrors(params types.IterateTestParams) (string, error) {
	// Format context files section
	var contextSection strings.Builder
	for _, cf := range params.ContextFiles {
		contextSection.WriteString(fmt.Sprintf("\n=== %s (%s) ===\n%s\n", cf.Description, cf.Type, cf.Content))
	}

	prompt := fmt.Sprintf(`Fix the type errors in this test code. The errors are:

%s

Here's the current test code:

%s

Which are found in the directory:

%s

Here are the relevant context files:

%s

Return ONLY the fixed test code, no explanations.`, params.Errors, params.TestCode, params.TestDir, contextSection.String())

	slog.Info("sending request to Anthropic API",
		"promptId", "fix_type_errors",
		"model", anthropic.ModelClaude3_5SonnetLatest,
		"maxTokens", MAX_TOKENS)

	msg, err := p.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(2000)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})

	if err != nil {
		if err := p.logger.Log("fix_type_errors", prompt, ""); err != nil {
			slog.Warn("failed to log prompt", "error", err)
		}
		slog.Error("failed to fix type errors with Anthropic API", "error", err)
		return "", fmt.Errorf("failed to fix type errors: %w", err)
	}

	response := removeBackticks(msg.Content[0].Text)

	// Log the prompt and response
	if err := p.logger.Log("fix_type_errors", prompt, response); err != nil {
		slog.Warn("failed to log prompt", "error", err)
	}

	slog.Debug("received response from Anthropic API",
		"responseLength", len(response))

	return response, nil
}
