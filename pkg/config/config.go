package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gwkline/artestian/types"
)

// ContextFile represents a single context file with metadata
type ContextFile struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Type        string `json:"type"` // e.g., "types", "utils", "constants", etc.
}

// Language represents a supported programming language
type Language string

// TestRunner represents a supported test runner
type TestRunner string

type Config struct {
	Version  string          `json:"version"`
	Examples []types.Example `json:"examples"`
	Settings types.Settings  `json:"settings"`
	Context  types.Context   `json:"context"`
	basePath string
}

const (
	TypeScript Language = "typescript"
	Go         Language = "go"
)

const (
	Jest   TestRunner = "jest"
	GoTest TestRunner = "go test"
)

// languageRunnerMap defines which test runners are compatible with each language
var languageRunnerMap = map[Language][]TestRunner{
	TypeScript: {Jest},
	Go:         {GoTest},
}

// defaultTestRunner defines the default test runner for each language
var defaultTestRunner = map[Language]TestRunner{
	TypeScript: Jest,
	Go:         GoTest,
}

// IsValidLanguage checks if the given language is supported
func IsValidLanguage(lang string) bool {
	slog.Debug("checking language validity", "language", lang, "supported_languages", languageRunnerMap)
	_, exists := languageRunnerMap[Language(lang)]
	return exists
}

// IsValidTestRunner checks if the given test runner is valid for a language
func IsValidTestRunner(lang Language, runner TestRunner) bool {
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

// LoadConfig loads and validates the configuration from a file
func LoadConfig(configPath string) (types.IConfig, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path is required")
	}

	// Get the absolute path of the config directory and ensure it has a trailing separator
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	absPath = filepath.Clean(absPath)

	// Check if path exists and is a directory
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("path must be a directory: %s", absPath)
	}

	// Look for artestian config file in directory
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var configFile string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if (filepath.Ext(name) == ".json" || filepath.Ext(name) == ".jsonc") && (strings.Contains(name, "artestian")) {
			configFile = filepath.Join(absPath, name)
			break
		}
	}

	if configFile == "" {
		return nil, fmt.Errorf("no artestian config file found in directory")
	}

	// Read and parse the config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Store the base path
	config.basePath = absPath

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig performs comprehensive validation of the configuration
func validateConfig(config *Config) error {
	if config.Version == "" {
		return fmt.Errorf("config version is required")
	}

	seenNames := make(map[string]bool)
	for i, example := range config.Examples {
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

		fullPath := config.resolveFilePath(example.FilePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("example #%d (%s): file not found at path: %s", i+1, example.Name, fullPath)
		}

		// Validate description
		if example.Description == "" {
			return fmt.Errorf("example #%d (%s): description is required", i+1, example.Name)
		}
	}

	// Validate context files if present
	for i, file := range config.Context.Files {
		if file.Path == "" {
			return fmt.Errorf("context file #%d: path is required", i+1)
		}

		fullPath := config.resolveFilePath(file.Path)
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

	contextDirs := make([]string, len(config.Context.Files))
	for i, file := range config.Context.Files {
		contextDirs[i] = filepath.Dir(config.resolveFilePath(file.Path))
	}
	slog.Debug("Validated config", "config", config, "num_context_files", len(config.Context.Files), "context_dirs", contextDirs)

	// Validate settings
	if config.Settings.DefaultTestDirectory == "" {
		slog.Warn("no default test directory specified, using current directory")
	}

	// Validate excluded directories
	for i, dir := range config.Settings.ExcludedDirs {
		if dir == "" {
			return fmt.Errorf("excluded directory at index %d cannot be empty", i)
		}

		// Strip ./ prefix if present
		dir = strings.TrimPrefix(dir, "./")

		// Remove trailing /** if present
		dir = strings.TrimSuffix(dir, "/**")

		// Remove trailing slash if present
		dir = strings.TrimSuffix(dir, "/")

		fullPath := config.resolveFilePath(dir)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			slog.Warn("excluded directory does not exist", "path", fullPath)
		}

		// Update the normalized path in the config
		config.Settings.ExcludedDirs[i] = dir
	}

	// Validate excluded files
	for i, file := range config.Settings.ExcludedFiles {
		if file == "" {
			return fmt.Errorf("excluded file at index %d cannot be empty", i)
		}

		// Remove leading ./ if present
		config.Settings.ExcludedFiles[i] = strings.TrimPrefix(file, "./")
	}

	// Language validation
	if config.Settings.Language != "" {
		if !IsValidLanguage(config.Settings.Language) {
			var supportedLangs []string
			for lang := range languageRunnerMap {
				supportedLangs = append(supportedLangs, string(lang))
			}
			return fmt.Errorf("unsupported language: %s. Must be one of: %s",
				config.Settings.Language, strings.Join(supportedLangs, ", "))
		}
	}

	// Test runner validation
	if config.Settings.TestRunner != "" {
		lang := Language(config.Settings.Language)
		if config.Settings.Language == "" {
			lang = Language(config.GetLanguage())
		}

		if !IsValidTestRunner(lang, TestRunner(config.Settings.TestRunner)) {
			validRunners := languageRunnerMap[lang]
			var runnerStrings []string
			for _, r := range validRunners {
				runnerStrings = append(runnerStrings, string(r))
			}
			return fmt.Errorf("unsupported test runner: %s for language %s. Must be one of: %s",
				config.Settings.TestRunner, lang, strings.Join(runnerStrings, ", "))
		}
	}

	return nil
}

// LoadExamples loads the test examples from the configuration
func (c *Config) LoadExamples() ([]types.TestExample, error) {
	examples := make([]types.TestExample, 0, len(c.Examples))

	for _, ex := range c.Examples {
		fullPath := c.resolveFilePath(ex.FilePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read example file %s: %w", fullPath, err)
		}

		examples = append(examples, types.TestExample{
			Name:        ex.Name,
			Type:        types.TestType(ex.Type),
			Content:     string(content),
			Description: ex.Description,
		})
	}

	return examples, nil
}

// resolveFilePath resolves a path relative to the config file location
func (c *Config) resolveFilePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(c.basePath, path)
}

// GetAbsolutePath converts a relative file path to an absolute path
func (c *Config) GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

// GetLanguage returns the configured language, defaulting to typescript
func (c *Config) GetLanguage() string {
	if c.Settings.Language == "" {
		return string(TypeScript)
	}
	return c.Settings.Language
}

// GetTestRunner returns the configured test runner, defaulting to the default for the language
func (c *Config) GetTestRunner() string {
	if c.Settings.TestRunner == "" {
		lang := Language(c.GetLanguage())
		return string(defaultTestRunner[lang])
	}
	return c.Settings.TestRunner
}

// GetRootDir returns the configured default test directory, defaulting to current directory
func (c *Config) GetRootDir() string {
	if c.Settings.DefaultTestDirectory == "" {
		return "."
	}
	return c.resolveFilePath(c.Settings.DefaultTestDirectory)
}

// GetExcludedDirs returns the normalized excluded directory paths
func (c *Config) GetExcludedDirs() []string {
	if len(c.Settings.ExcludedDirs) == 0 {
		return nil
	}

	// Get the absolute path of the root directory
	rootDir := c.GetRootDir()
	excludedDirs := make([]string, len(c.Settings.ExcludedDirs))

	for i, dir := range c.Settings.ExcludedDirs {
		// If the path starts with ./, use the root directory as base
		if strings.HasPrefix(dir, "./") {
			excludedDirs[i] = filepath.Join(rootDir, strings.TrimPrefix(dir, "./"))
		} else {
			// For paths without ./, resolve them relative to the config base path
			excludedDirs[i] = c.resolveFilePath(dir)
		}
	}

	return excludedDirs
}

// LoadContextFiles loads all context files specified in the configuration
func (c *Config) LoadContextFiles() ([]types.ContextFile, error) {
	if len(c.Context.Files) == 0 {
		return nil, nil
	}

	files := make([]types.ContextFile, 0, len(c.Context.Files))
	for _, file := range c.Context.Files {
		fullPath := c.resolveFilePath(file.Path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read context file %s: %w", fullPath, err)
		}

		files = append(files, types.ContextFile{
			Path:        file.Path,
			Content:     string(content),
			Description: file.Description,
			Type:        file.Type,
		})
	}

	return files, nil
}

// GetExcludedFiles returns the list of excluded file patterns
func (c *Config) GetExcludedFiles() []string {
	if len(c.Settings.ExcludedFiles) == 0 {
		return nil
	}
	return c.Settings.ExcludedFiles
}
