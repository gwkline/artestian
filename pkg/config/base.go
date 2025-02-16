package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gwkline/artestian/types"
)

type ContextFile struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Type        string `json:"type"` // e.g., "types", "utils", "constants", etc.
}
type Config struct {
	Version  string          `json:"version"`
	Examples []types.Example `json:"examples"`
	Settings types.Settings  `json:"settings"`
	Context  types.Context   `json:"context"`
	basePath string
}

type Language string

const (
	TypeScript Language = "typescript"
	Go         Language = "go"
)

type TestRunner string

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

// LoadConfig loads and validates the configuration from a file
func Init(configPath string) (types.IConfig, error) {
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
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// resolveFilePath resolves a path relative to the config file location
func (c *Config) resolveFilePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(c.basePath, path)
}
