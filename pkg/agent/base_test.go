package agent

import (
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gwkline/artestian/pkg/prompt_logger"
	"github.com/stretchr/testify/assert"
)

func TestNewAnthropicProvider(t *testing.T) {
	tests := []struct {
		name          string
		apiKeyEnvVar  string
		expectedError string
	}{
		{
			name:          "success case with valid API key",
			apiKeyEnvVar:  "test-api-key",
			expectedError: "",
		},
		{
			name:          "error case with empty API key",
			apiKeyEnvVar:  "",
			expectedError: "ANTHROPIC_API_KEY environment variable not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.apiKeyEnvVar != "" {
				oldApiKey := os.Getenv("ANTHROPIC_API_KEY")
				os.Setenv("ANTHROPIC_API_KEY", tt.apiKeyEnvVar)
				defer os.Setenv("ANTHROPIC_API_KEY", oldApiKey)
			} else {
				oldApiKey := os.Getenv("ANTHROPIC_API_KEY")
				os.Unsetenv("ANTHROPIC_API_KEY")
				defer os.Setenv("ANTHROPIC_API_KEY", oldApiKey)
			}

			// Execute
			logger, err := prompt_logger.Init(false)
			if err != nil {
				t.Fatalf("failed to create prompt logger: %v", err)
			}
			provider, err := NewAnthropicProvider(logger)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
				assert.IsType(t, &anthropic.Client{}, provider.client)
				assert.NotNil(t, provider.logger)
			}
		})
	}
}
