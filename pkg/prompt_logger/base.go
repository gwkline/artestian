package prompt_logger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// promptLogger handles saving prompts to files for debugging and analysis
type promptLogger struct {
	logsDir         string
	loggingDisabled bool
}

// Init creates a new prompt logger that saves to a logs directory in the current working directory
func Init(enabled bool) (*promptLogger, error) {
	if !enabled {
		return &promptLogger{
			logsDir:         "",
			loggingDisabled: true,
		}, nil
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Create logs directory in current working directory
	logsDir := filepath.Join(cwd, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create prompt logs directory: %w", err)
	}

	return &promptLogger{
		logsDir: logsDir,
	}, nil
}

// logPrompt saves a prompt and its response to a file as JSON
func (l *promptLogger) Log(operation string, prompt string, response string) error {
	if l.loggingDisabled {
		return nil
	}

	// Create timestamp-based filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.json", timestamp, operation)
	fullPath := filepath.Join(l.logsDir, filename)

	// Create log entry as JSON object
	logEntry := map[string]interface{}{
		"operation": operation,
		"timestamp": timestamp,
		"prompt":    prompt,
		"response":  response,
	}

	// Convert to JSON
	content, err := json.MarshalIndent(logEntry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal log entry to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write prompt log: %w", err)
	}

	slog.Debug("saved prompt log", "path", fullPath)
	return nil
}
