package types

type IAgent interface {
	GenerateTest(params GenerateTestParams) (string, error)
	FixTypeErrors(params IterateTestParams) (string, error)
	FixTestErrors(params IterateTestParams) (string, error)
	PickExample(sourceCode string, testExamples []TestExample) (TestExample, error)
}

// TestRunner interface for different test frameworks
type ITestRunner interface {
	GetName() string
	RunTests(rootDir, testFilePath string) (bool, string, error)
}

// FileFinder interface for finding files that need tests
type IFileFinder interface {
	FindNextFile(rootDir string, excludeDirs []string) (string, error)
	GetTestPath(sourcePath string) string
}

// LanguageSupport interface for language-specific operations
type ILanguage interface {
	GetName() string
	GetTestRunner() ITestRunner
	GetFileExtension() string
	GetTestFilePattern() string
	CheckTypes(testFilePath string) (bool, string, error)
}

type IPromptLogger interface {
	Log(operation string, prompt string, response string) error
}
