package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"path/filepath"
	"strings"

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

func (g *GoSupport) GetFunctions(sourceCode string) ([]types.Function, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceCode, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var functions []types.Function
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Get the function source code
			startPos := fset.Position(funcDecl.Pos())
			endPos := fset.Position(funcDecl.End())
			lines := strings.Split(sourceCode, "\n")
			funcSource := strings.Join(lines[startPos.Line-1:endPos.Line], "\n")

			// Create Function struct
			function := types.Function{
				Name:       funcDecl.Name.Name,
				SourceCode: funcSource,
				IsExported: funcDecl.Name.IsExported(),
			}
			functions = append(functions, function)
		}
		return true
	})

	return functions, nil
}
