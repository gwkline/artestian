package agent

import (
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/gwkline/artestian/types"
)

var MAX_TOKENS = 8192

type AnthropicProvider struct {
	client *anthropic.Client
	logger types.IPromptLogger
}

func NewAnthropicProvider(logger types.IPromptLogger) (*AnthropicProvider, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	return &AnthropicProvider{
		client: anthropic.NewClient(
			option.WithAPIKey(apiKey),
		),
		logger: logger,
	}, nil
}
