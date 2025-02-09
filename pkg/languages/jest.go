package languages

import (
	"os/exec"
	"path/filepath"
)

type JestRunner struct{}

func (r *JestRunner) RunTests(rootDir, testFilePath string) (bool, string, error) {
	dir := filepath.Dir(rootDir)

	cmd := exec.Command("npx", "jest", testFilePath, "--no-cache")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Jest returns non-zero exit code on test failures
		return false, string(output), nil
	}
	return true, string(output), nil
}

func (r *JestRunner) GetName() string {
	return "jest"
}
