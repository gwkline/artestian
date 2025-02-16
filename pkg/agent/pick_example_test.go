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
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

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
					SourceCode: `
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    assert.Equal(t, 5, result)
}`,
					Description: "Basic unit test for addition function",
				},
				{
					Name: "Complex Integration Test",
					Type: types.TestTypeIntegration,
					SourceCode: `
func TestComplexScenario(t *testing.T) {
    // Complex test setup
}`,
					Description: "Complex integration test",
				},
			},
			expected: types.TestExample{
				Name: "Simple Unit Test",
				Type: types.TestTypeUnit,
				SourceCode: `
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

func TestParseExampleIndex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "simple digit",
			input:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "digit with json closing",
			input:   "1}",
			want:    1,
			wantErr: false,
		},
		{
			name:    "digit with json closing and space",
			input:   "2 }",
			want:    2,
			wantErr: false,
		},
		{
			name:    "full json format",
			input:   `{"exampleIndex":3}`,
			want:    3,
			wantErr: false,
		},
		{
			name:    "invalid input - letters",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "invalid input - empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExampleIndex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExampleIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseExampleIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
