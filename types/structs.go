package types

type TestType string

const (
	TestTypeUnit        TestType = "unit"
	TestTypeIntegration TestType = "integration"
	TestTypeWorker      TestType = "worker"
	TestTypePrompt      TestType = "prompt"
)

type TestExample struct {
	Name        string
	Type        TestType
	Content     string
	Description string
}

type ErrorAttempt struct {
	Code   string
	Errors []string
}

type GenerateTestParams struct {
	SourceCode   string
	Example      TestExample
	Language     ILanguage
	TestRunner   ITestRunner
	TestDir      string
	ContextFiles []ContextFile // Additional context files for test generation
	Function     Function      // Function to generate a test for
}

type IterateTestParams struct {
	SourceCode   string
	TestCode     string
	Errors       []string
	Example      TestExample
	Language     ILanguage
	TestRunner   ITestRunner
	TestDir      string
	ContextFiles []ContextFile // Additional context files for test generation
}

// ContextFile represents a file that provides additional context for test generation
type ContextFile struct {
	Path        string // Path to the file relative to the config file
	Content     string // Content of the file
	Description string // Description of what this file contains/provides
	Type        string // Type of context (e.g., "types", "utils", "constants")
}

// Example represents a single test example configuration
type Example struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
}

// Settings represents global configuration settings
type Settings struct {
	DefaultTestDirectory string   `json:"default_test_directory"`
	Language             string   `json:"language"`
	TestRunner           string   `json:"test_runner"`
	ExcludedDirs         []string `json:"excluded_dirs"`
	ExcludedFiles        []string `json:"excluded_files"`
}

// Context represents additional files to be used as context for test generation
type Context struct {
	Files []ContextFile `json:"files"`
}

type Function struct {
	Name       string
	SourceCode string
	IsExported bool
}
