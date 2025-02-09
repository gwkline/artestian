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
