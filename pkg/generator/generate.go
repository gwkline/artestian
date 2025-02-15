package generator

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gwkline/artestian/types"
)

func (g *TestGenerator) GenerateNextTest(projectDir string, cfg types.IConfig) error {
	slog.Debug("finding next file that needs tests", "rootDir", cfg.GetRootDir())

	sourcePath, err := g.finder.FindNextFile(cfg)
	if err != nil {
		slog.Error("failed to find next file", "error", err)
		return fmt.Errorf("error finding file: %w", err)
	}
	if sourcePath == "" {
		slog.Info("no files found needing tests")
		return fmt.Errorf("no files found needing tests")
	}

	sourceCode, err := os.ReadFile(sourcePath)
	if err != nil {
		slog.Error("failed to read source file", "path", sourcePath, "error", err)
		return fmt.Errorf("error reading source file: %w", err)
	}

	functions, err := g.language.GetFunctions(string(sourceCode))
	if err != nil {
		slog.Error("failed to get functions", "path", sourcePath, "error", err)
		return fmt.Errorf("error getting functions: %w", err)
	}

	if len(functions) == 0 {
		slog.Error("no functions found in source file", "path", sourcePath)
		return fmt.Errorf("no functions found in source file")
	}

	slog.Debug("finding best example for source code")
	example := g.findBestExample(string(sourceCode))

	var allTestCode string
	testPath := g.finder.GetTestPath(sourcePath)

	for _, function := range functions {
		slog.Info("generating test for function", "function", function.Name)
		testCode, err := g.ai.GenerateTest(types.GenerateTestParams{
			SourceCode:   string(sourceCode),
			Function:     function,
			Example:      example,
			Language:     g.language,
			TestRunner:   g.language.GetTestRunner(),
			TestDir:      testPath,
			ContextFiles: g.contextFiles,
		})
		if err != nil {
			slog.Error("failed to generate test", "function", function.Name, "error", err)
			continue
		}

		// Create temp file in the test directory
		tempFile, err := os.CreateTemp(filepath.Dir(testPath), fmt.Sprintf("%s*%s", function.Name, g.language.GetTestFilePattern()))
		if err != nil {
			slog.Error("failed to create temp file", "function", function.Name, "error", err)
			continue
		}
		tempPath := tempFile.Name()
		tempFile.Close()

		if err := os.WriteFile(tempPath, []byte(testCode), 0644); err != nil {
			slog.Error("failed to write temp test file", "path", tempPath, "error", err)
			continue
		}

		testCode, err = g.iterateTypeErrors(string(sourceCode), tempPath, testCode)
		if err != nil {
			slog.Error("error fixing type errors", "function", function.Name, "error", err)
			continue
		}

		testCode, err = g.iterateTestFailures(projectDir, tempPath, testCode)
		if err != nil {
			slog.Error("error fixing test errors", "function", function.Name, "error", err)
			continue
		}

		// Clean up the temp file only if the test code is valid
		// otherwise, keep it around for further iteration
		defer os.Remove(tempPath)

		allTestCode += testCode + "\n"
	}

	if allTestCode != "" {
		slog.Info("writing final test file", "path", testPath)

		if err := os.MkdirAll(filepath.Dir(testPath), 0755); err != nil {
			slog.Error("failed to create test directory", "path", filepath.Dir(testPath), "error", err)
			return fmt.Errorf("error creating test directory: %w", err)
		}

		if err := os.WriteFile(testPath, []byte(allTestCode), 0644); err != nil {
			slog.Error("failed to write test file", "path", testPath, "error", err)
			return fmt.Errorf("error writing test file: %w", err)
		}
	}

	return nil
}
