package prompt_logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromptLogger(t *testing.T) {
	// Create temp directory for tests
	tempDir, err := os.MkdirTemp("", "prompt_logger_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change working directory to temp dir for test
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	tests := []struct {
		name          string
		enabled       bool
		operation     string
		prompt        string
		response      string
		expectLogFile bool
		expectError   bool
	}{
		{
			name:          "logging enabled with valid inputs",
			enabled:       true,
			operation:     "test_op",
			prompt:        "test prompt",
			response:      "test response",
			expectLogFile: true,
			expectError:   false,
		},
		{
			name:          "logging disabled",
			enabled:       false,
			operation:     "test_op",
			prompt:        "test prompt",
			response:      "test response",
			expectLogFile: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Remove logs directory if it exists
			os.RemoveAll(filepath.Join(tempDir, "logs"))

			// Initialize logger
			logger, err := Init(tt.enabled)
			assert.NotNil(t, logger)

			// Log prompt
			err = logger.Log(tt.operation, tt.prompt, tt.response)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check if log file exists
			if tt.expectLogFile {
				// Get log directory contents
				files, err := os.ReadDir(filepath.Join(tempDir, "logs"))
				assert.NoError(t, err)
				assert.Len(t, files, 1)

				// Read and verify log file contents
				logFile := files[0]
				content, err := os.ReadFile(filepath.Join(tempDir, "logs", logFile.Name()))
				assert.NoError(t, err)

				var logEntry map[string]interface{}
				err = json.Unmarshal(content, &logEntry)
				assert.NoError(t, err)

				assert.Equal(t, tt.operation, logEntry["operation"])
				assert.Equal(t, tt.prompt, logEntry["prompt"])
				assert.Equal(t, tt.response, logEntry["response"])

				// Verify timestamp format
				timestamp := logEntry["timestamp"].(string)
				_, err = time.Parse("2006-01-02_15-04-05", timestamp)
				assert.NoError(t, err)
			} else {
				// Verify no log directory created when disabled
				_, err := os.Stat(filepath.Join(tempDir, "logs"))
				assert.True(t, os.IsNotExist(err))
			}
		})
	}
}
