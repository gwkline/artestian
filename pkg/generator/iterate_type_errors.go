package generator

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gwkline/artestian/types"
)

func (g *TestGenerator) iterateTypeErrors(params types.GenerateTestParams, testCode string) (string, error) {
	var attempts []types.ErrorAttempt
	maxTypeAttempts := 5

	for i := 0; i < maxTypeAttempts; i++ {
		slog.Debug("checking types", "attempt", i+1, "path", params.TestPath)
		ok, typeErrors, err := g.language.CheckTypes(params.TestPath)
		if err != nil {
			return "", fmt.Errorf("error checking types: %w", err)
		}
		if ok {
			slog.Info("type check passed")
			return testCode, nil
		}

		slog.Debug("type errors", "errors", typeErrors)
		attempts = append(attempts, types.ErrorAttempt{
			Code:   testCode,
			Errors: strings.Split(typeErrors, "\n"),
		})

		slog.Info("fixing type errors", "attempt", i+1)
		fixedCode, err := g.ai.FixTypeErrors(types.IterateTestParams{
			GenerateTestParams: params,
			TestCode:           testCode,
			Errors:             strings.Split(typeErrors, "\n"),
		})
		if err != nil {
			return "", fmt.Errorf("error fixing type errors: %w", err)
		}

		testCode = fixedCode
		if err := os.WriteFile(params.TestPath, []byte(testCode), 0644); err != nil {
			return "", fmt.Errorf("error writing fixed test file: %w", err)
		}
	}

	return "", fmt.Errorf("failed to fix type errors after %d attempts", maxTypeAttempts)
}
