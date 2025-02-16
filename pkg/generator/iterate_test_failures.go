package generator

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gwkline/artestian/types"
)

func (g *TestGenerator) iterateTestFailures(params types.GenerateTestParams, testCode string, projectDir string) (string, error) {
	var attempts []types.ErrorAttempt
	maxTestAttempts := 3
	runner := g.language.GetTestRunner()

	for i := 0; i < maxTestAttempts; i++ {
		slog.Debug("running tests", "attempt", i+1, "path", params.TestPath, "projectDir", projectDir)
		ok, testErrors, err := runner.RunTests(projectDir, params.TestPath)
		if err != nil {
			return "", fmt.Errorf("error running tests: %w", err)
		}
		if ok {
			slog.Info("tests passed")
			return testCode, nil
		}

		slog.Debug("test errors", "errors", testErrors)
		attempts = append(attempts, types.ErrorAttempt{
			Code:   testCode,
			Errors: strings.Split(testErrors, "\n"),
		})

		slog.Info("fixing test errors", "attempt", i+1)
		fixedCode, err := g.ai.FixTestFailures(types.IterateTestParams{
			GenerateTestParams: params,
			TestCode:           testCode,
			Errors:             strings.Split(testErrors, "\n"),
		})
		if err != nil {
			return "", fmt.Errorf("error fixing test errors: %w", err)
		}

		testCode = fixedCode
		if err := os.WriteFile(params.TestPath, []byte(testCode), 0644); err != nil {
			return "", fmt.Errorf("error writing fixed test file: %w", err)
		}
	}

	return "", fmt.Errorf("failed to fix test errors after %d attempts", maxTestAttempts)
}
