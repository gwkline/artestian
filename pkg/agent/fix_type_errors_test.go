package agent

import (
	"log"
	"strings"
	"testing"

	"github.com/gwkline/artestian/pkg/prompt_logger"
	"github.com/gwkline/artestian/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestFixTypeErrors(t *testing.T) {
	tests := []struct {
		name          string
		testCode      string
		errors        []string
		contextFiles  []types.ContextFile
		expected      string
		expectedError string
	}{
		{
			name: "success case",
			testCode: `
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAddition(t *testing.T) {
    result := Addition("2", 3)
    assert.Equal(t, 5, result)
}`,
			errors: []string{"cannot use \"2\" (type string) as type int in argument to Addition"},
			contextFiles: []types.ContextFile{
				{
					Description: "source code",
					Type:        "go",
					Content: `
func Addition(a int, b int) int {
    return a + b
}`,
				},
			},
			expected: `
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAddition(t *testing.T) {
    result := Addition(2, 3)
    assert.Equal(t, 5, result)
}`,
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

			result, err := provider.FixTypeErrors(types.IterateTestParams{
				TestCode:     tt.testCode,
				Errors:       tt.errors,
				ContextFiles: tt.contextFiles,
			})

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(tt.expected), strings.TrimSpace(result))
			}
		})
	}
}
