package agent

import (
	"log"
	"testing"

	"github.com/gwkline/artestian/pkg/prompt_logger"
	"github.com/gwkline/artestian/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestPickExample(t *testing.T) {
	tests := []struct {
		name          string
		sourceCode    string
		testExamples  []types.TestExample
		expected      types.TestExample
		expectedError string
	}{
		{
			name: "success case",
			sourceCode: `
func Add(a, b int) int {
    return a + b
}`,
			testExamples: []types.TestExample{
				{
					Name: "Simple Unit Test",
					Type: types.TestTypeUnit,
					Content: `
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    assert.Equal(t, 5, result)
}`,
					Description: "Basic unit test for addition function",
				},
				{
					Name: "Complex Integration Test",
					Type: types.TestTypeIntegration,
					Content: `
func TestComplexScenario(t *testing.T) {
    // Complex test setup
}`,
					Description: "Complex integration test",
				},
			},
			expected: types.TestExample{
				Name: "Simple Unit Test",
				Type: types.TestTypeUnit,
				Content: `
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    assert.Equal(t, 5, result)
}`,
				Description: "Basic unit test for addition function",
			},
		},
		{
			name:          "error case - empty examples",
			sourceCode:    "func Add(a, b int) int { return a + b }",
			testExamples:  []types.TestExample{},
			expectedError: "no test examples provided",
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

			result, err := provider.PickExample(tt.sourceCode, tt.testExamples)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
