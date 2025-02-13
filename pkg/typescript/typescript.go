package typescript

import (
	"os/exec"
	"path/filepath"

	"github.com/gwkline/artestian/types"
)

type TypeScriptSupport struct{}

func NewTypeScriptSupport() *TypeScriptSupport {
	return &TypeScriptSupport{}
}

func (ts *TypeScriptSupport) GetTestRunner() types.ITestRunner {
	return &JestRunner{}
}

func (ts *TypeScriptSupport) GetFileExtension() string {
	return ".ts"
}

func (ts *TypeScriptSupport) GetTestFilePattern() string {
	return ".test.ts"
}

func (r *TypeScriptSupport) CheckTypes(testFilePath string) (bool, string, error) {
	// Get the directory of the test file
	dir := filepath.Dir(testFilePath)

	cmd := exec.Command("npx", "tsc", "--noEmit")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, string(output), nil
	}
	return true, string(output), nil
}

func (ts *TypeScriptSupport) GetName() string {
	return "typescript"
}
