package generator

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gwkline/artestian/types"
)

func (g *TestGenerator) GenerateNextTest(projectDir, rootDir string, excludeDirs []string) error {
	slog.Debug("finding next file that needs tests", "rootDir", rootDir)

	sourcePath, err := g.finder.FindNextFile(rootDir, excludeDirs)
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

	slog.Debug("finding best example for source code")
	example := g.findBestExample(string(sourceCode))

	slog.Info("generating initial test")
	testCode, err := g.ai.GenerateTest(types.GenerateTestParams{
		SourceCode:   string(sourceCode),
		Example:      example,
		Language:     g.language,
		TestRunner:   g.language.GetTestRunner(),
		TestDir:      g.finder.GetTestPath(sourcePath),
		ContextFiles: g.contextFiles,
	})
	if err != nil {
		slog.Error("failed to generate test", "error", err)
		return fmt.Errorf("error generating test: %w", err)
	}

	testPath := g.finder.GetTestPath(sourcePath)
	slog.Info("writing test file", "path", testPath)

	if err := os.MkdirAll(filepath.Dir(testPath), 0755); err != nil {
		slog.Error("failed to create test directory", "path", filepath.Dir(testPath), "error", err)
		return fmt.Errorf("error creating test directory: %w", err)
	}

	if err := os.WriteFile(testPath, []byte(testCode), 0644); err != nil {
		slog.Error("failed to write test file", "path", testPath, "error", err)
		return fmt.Errorf("error writing test file: %w", err)
	}

	testCode, err = g.iterateTypeErrors(testPath, testCode)
	if err != nil {
		return fmt.Errorf("error fixing type errors: %w", err)
	}

	testCode, err = g.iterateTestFailures(projectDir, testPath, testCode)
	if err != nil {
		return fmt.Errorf("error fixing test errors: %w", err)
	}

	return nil
}
