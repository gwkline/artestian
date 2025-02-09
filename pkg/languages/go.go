package languages

import (
	"os/exec"
	"path/filepath"

	"github.com/gwkline/artestian/types"
)

type GoTestRunner struct{}

func (r *GoTestRunner) RunTests(rootDir, testFilePath string) (bool, string, error) {
	dir := filepath.Dir(rootDir)

	cmd := exec.Command("go", "test", testFilePath)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Go test returns non-zero exit code on test failures
		return false, string(output), nil
	}
	return true, string(output), nil
}

func (r *GoTestRunner) GetName() string {
	return "go test"
}

type GoSupport struct{}

func NewGoSupport() *GoSupport {
	return &GoSupport{}
}

func (g *GoSupport) GetTestRunner() types.ITestRunner {
	return &GoTestRunner{}
}

func (g *GoSupport) GetFileExtension() string {
	return ".go"
}

func (g *GoSupport) GetTestFilePattern() string {
	return "_test.go"
}

func (g *GoSupport) CheckTypes(testFilePath string) (bool, string, error) {
	// Get the directory of the test file
	dir := filepath.Dir(testFilePath)

	cmd := exec.Command("go", "vet")
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, string(output), nil
	}
	return true, string(output), nil
}

func (g *GoSupport) GetName() string {
	return "go"
}
