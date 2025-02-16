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

func TestFixTestErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tests := []struct {
		name          string
		sourceCode    string
		testCode      string
		testErrors    string
		contextFiles  []types.ContextFile
		expected      string
		expectedError string
	}{
		{
			name: "success case",
			sourceCode: `
func Addition(a int, b int) int {
    return a + b
}`,
			testCode: `
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAddition(t *testing.T) {
    result := Addition(2, 3)
    assert.Equal(t, 4, result)
}`,
			testErrors: "--- FAIL: TestAddition (0.00s)\n    addition_test.go:6: \n        	Error Trace:	addition_test.go:6\n        	Error:      	Not equal: \n        	            	expected: 4\n        	            	actual  : 5\n        	Test:       	TestAddition",
			contextFiles: []types.ContextFile{
				{
					Description: "test context",
					Type:        "test",
					Content:     "test content",
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
			// Setup
			logger, err := prompt_logger.Init(false)
			assert.NoError(t, err)

			provider, err := NewAnthropicProvider(logger)
			assert.NoError(t, err)

			// Execute
			result, err := provider.FixTestFailures(types.IterateTestParams{
				GenerateTestParams: types.GenerateTestParams{
					SourceCode:   tt.sourceCode,
					ContextFiles: tt.contextFiles,
				},
				TestCode: tt.testCode,
				Errors:   []string{tt.testErrors},
			})

			// Assert
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
