package agent

import (
	"log"
	"testing"

	"github.com/gwkline/artestian/pkg/languages"
	"github.com/gwkline/artestian/pkg/prompt_logger"
	"github.com/gwkline/artestian/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTest(t *testing.T) {
	tests := []struct {
		name          string
		params        types.GenerateTestParams
		expected      string
		expectedError string
	}{
		{
			name: "successful test generation",
			params: types.GenerateTestParams{
				SourceCode: `
func Multiply(a int, b int) int {
    return a * b
}`,
				Example: types.TestExample{
					Name:    "simple unit test",
					Type:    types.TestTypeUnit,
					Content: "test content",
				},
				ContextFiles: []types.ContextFile{
					{
						Description: "test context",
						Type:        "test",
						Content:     "test content",
					},
				},
				Language:   languages.NewGoSupport(),
				TestRunner: &languages.GoTestRunner{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := godotenv.Load("../../.env"); err != nil {
				log.Fatal("Error loading .env file")
			}

			logger, err := prompt_logger.Init(false)
			assert.NoError(t, err)

			provider, err := NewAnthropicProvider(logger)
			assert.NoError(t, err)

			result, err := provider.GenerateTest(tt.params)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveBackticks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "code with language tag",
			input:    "```go\npackage main\n\nfunc main() {}\n```",
			expected: "package main\n\nfunc main() {}",
		},
		{
			name:     "code without language tag",
			input:    "```\nsome code\n```",
			expected: "some code",
		},
		{
			name:     "plain text",
			input:    "plain text",
			expected: "plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeBackticks(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
