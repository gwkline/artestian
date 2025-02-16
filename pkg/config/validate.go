package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gwkline/artestian/types"
)

// validateConfig performs comprehensive validation of the configuration
func (c *Config) validate() error {
	if c.Version == "" {
		return fmt.Errorf("config version is required")
	}

	seenNames := make(map[string]bool)
	for i, example := range c.Examples {
		// Validate example name
		if example.Name == "" {
			return fmt.Errorf("example #%d: name is required", i+1)
		}
		if seenNames[example.Name] {
			return fmt.Errorf("duplicate example name found: %s", example.Name)
		}
		seenNames[example.Name] = true

		// Validate example type
		if example.Type == "" {
			return fmt.Errorf("example #%d (%s): type is required", i+1, example.Name)
		}

		// Validate test type is one of the allowed values
		switch example.Type {
		case string(types.TestTypeUnit),
			string(types.TestTypeIntegration),
			string(types.TestTypeWorker),
			string(types.TestTypePrompt):
			// Valid type
		default:
			return fmt.Errorf("example #%d (%s): invalid test type %q. Must be one of: unit, integration, worker, prompt", i+1, example.Name, example.Type)
		}

		// Validate and check file path
		if example.FilePath == "" {
			return fmt.Errorf("example #%d (%s): file_path is required", i+1, example.Name)
		}

		fullPath := c.resolveFilePath(example.FilePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("example #%d (%s): file not found at path: %s", i+1, example.Name, fullPath)
		}

		// Validate description
		if example.Description == "" {
			return fmt.Errorf("example #%d (%s): description is required", i+1, example.Name)
		}
	}

	// Validate context files if present
	for i, file := range c.Context.Files {
		if file.Path == "" {
			return fmt.Errorf("context file #%d: path is required", i+1)
		}

		fullPath := c.resolveFilePath(file.Path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("context file #%d: file not found at path: %s", i+1, fullPath)
		}

		if file.Description == "" {
			return fmt.Errorf("context file #%d: description is required", i+1)
		}

		if file.Type == "" {
			return fmt.Errorf("context file #%d: type is required", i+1)
		}
	}

	contextDirs := make([]string, len(c.Context.Files))
	for i, file := range c.Context.Files {
		contextDirs[i] = filepath.Dir(c.resolveFilePath(file.Path))
	}
	slog.Debug("Validated config", "config", c, "num_context_files", len(c.Context.Files), "context_dirs", contextDirs)

	// Validate settings
	if c.Settings.DefaultTestDirectory == "" {
		slog.Warn("no default test directory specified, using current directory")
	}

	// Validate excluded directories
	for i, dir := range c.Settings.ExcludedDirs {
		if dir == "" {
			return fmt.Errorf("excluded directory at index %d cannot be empty", i)
		}

		// Strip ./ prefix if present
		dir = strings.TrimPrefix(dir, "./")

		// Remove trailing /** if present
		dir = strings.TrimSuffix(dir, "/**")

		// Remove trailing slash if present
		dir = strings.TrimSuffix(dir, "/")

		fullPath := c.resolveFilePath(dir)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			slog.Warn("excluded directory does not exist", "path", fullPath)
		}

		// Update the normalized path in the config
		c.Settings.ExcludedDirs[i] = dir
	}

	// Validate excluded files
	for i, file := range c.Settings.ExcludedFiles {
		if file == "" {
			return fmt.Errorf("excluded file at index %d cannot be empty", i)
		}

		// Remove leading ./ if present
		c.Settings.ExcludedFiles[i] = strings.TrimPrefix(file, "./")
	}

	// Language validation
	if c.Settings.Language != "" {
		if !isValidLanguage(c.Settings.Language) {
			var supportedLangs []string
			for lang := range languageRunnerMap {
				supportedLangs = append(supportedLangs, string(lang))
			}
			return fmt.Errorf("unsupported language: %s. Must be one of: %s",
				c.Settings.Language, strings.Join(supportedLangs, ", "))
		}
	}

	// Test runner validation
	if c.Settings.TestRunner != "" {
		lang := Language(c.Settings.Language)
		if c.Settings.Language == "" {
			lang = Language(c.GetLanguage())
		}

		if !isValidTestRunner(lang, TestRunner(c.Settings.TestRunner)) {
			validRunners := languageRunnerMap[lang]
			var runnerStrings []string
			for _, r := range validRunners {
				runnerStrings = append(runnerStrings, string(r))
			}
			return fmt.Errorf("unsupported test runner: %s for language %s. Must be one of: %s",
				c.Settings.TestRunner, lang, strings.Join(runnerStrings, ", "))
		}
	}

	return nil
}

// IsValidLanguage checks if the given language is supported
func isValidLanguage(lang string) bool {
	slog.Debug("checking language validity", "language", lang, "supported_languages", languageRunnerMap)
	_, exists := languageRunnerMap[Language(lang)]
	return exists
}

// IsValidTestRunner checks if the given test runner is valid for a language
func isValidTestRunner(lang Language, runner TestRunner) bool {
	validRunners, exists := languageRunnerMap[lang]
	if !exists {
		return false
	}
	for _, r := range validRunners {
		if r == runner {
			return true
		}
	}
	return false
}
