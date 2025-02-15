package typescript

import (
	"os/exec"
)

type JestRunner struct{}

func (r *JestRunner) RunTests(rootDir, testFilePath string) (bool, string, error) {
	cmd := exec.Command("npx", "jest", testFilePath, "--no-cache")
	cmd.Dir = rootDir

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
