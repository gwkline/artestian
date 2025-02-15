package golang

import (
	"strings"
	"testing"

	"github.com/gwkline/artestian/types"
	"github.com/stretchr/testify/assert"
)

func TestGoSupport_GetFunctions(t *testing.T) {
	g := NewGoSupport()

	tests := []struct {
		name     string
		input    string
		expected []types.Function
	}{
		{
			name: "regular functions",
			input: `package main

				func hello() string { return "hello" }
				func World() string { return "world" }
			`,
			expected: []types.Function{
				{
					Name:       "hello",
					SourceCode: `func hello() string { return "hello" }`,
					IsExported: false,
				},
				{
					Name:       "World",
					SourceCode: `func World() string { return "world" }`,
					IsExported: true,
				},
			},
		},
		{
			name: "method functions",
			input: `package main

				type MyStruct struct{}
				func (m MyStruct) hello() string { return "hello" }
				func (m *MyStruct) World() string { return "world" }
			`,
			expected: []types.Function{
				{
					Name:       "hello",
					SourceCode: `func (m MyStruct) hello() string { return "hello" }`,
					IsExported: false,
				},
				{
					Name:       "World",
					SourceCode: `func (m *MyStruct) World() string { return "world" }`,
					IsExported: true,
				},
			},
		},
		{
			name: "function with comments",
			input: `package main

				// SimpleFunction returns a string
				func SimpleFunction() string {
					return "simple"
				}
			`,
			expected: []types.Function{
				{
					Name: "SimpleFunction",
					SourceCode: `func SimpleFunction() string {
					return "simple"
				}`,
					IsExported: true,
				},
			},
		},
		{
			name: "function with nested braces",
			input: `package main

				func ComplexFunction() interface{} {
					if true {
						return struct{
							nested struct{
								data string
							}
						}{
							nested: struct{
								data string
							}{
								data: "test",
							},
						}
					}
					return nil
				}
			`,
			expected: []types.Function{
				{
					Name: "ComplexFunction",
					SourceCode: `func ComplexFunction() interface{} {
					if true {
						return struct{
							nested struct{
								data string
							}
						}{
							nested: struct{
								data string
							}{
								data: "test",
							},
						}
					}
					return nil
				}`,
					IsExported: true,
				},
			},
		},
		{
			name: "complex mixed example",
			input: `// Package comment
package main

import "fmt"

// Regular function
func simple() {
	fmt.Println("simple")
}

// Exported function with multiple parameters
func ProcessData(input string, count int) (string, error) {
	return input, nil
}

// Method on a struct
type Handler struct{}

func (h *Handler) handle() error {
	return nil
}

// Exported method
func (h *Handler) Process() error {
	return nil
}

// Function with complex return type
func GetConfig() struct{
	Name string
	Value int
} {
	return struct{
		Name string
		Value int
	}{
		Name: "test",
		Value: 42,
	}
}`,
			expected: []types.Function{
				{
					Name: "simple",
					SourceCode: `func simple() {
	fmt.Println("simple")
}`,
					IsExported: false,
				},
				{
					Name: "ProcessData",
					SourceCode: `func ProcessData(input string, count int) (string, error) {
	return input, nil
}`,
					IsExported: true,
				},
				{
					Name: "handle",
					SourceCode: `func (h *Handler) handle() error {
	return nil
}`,
					IsExported: false,
				},
				{
					Name: "Process",
					SourceCode: `func (h *Handler) Process() error {
	return nil
}`,
					IsExported: true,
				},
				{
					Name: "GetConfig",
					SourceCode: `func GetConfig() struct{
	Name string
	Value int
} {
	return struct{
		Name string
		Value int
	}{
		Name: "test",
		Value: 42,
	}
}`,
					IsExported: true,
				},
			},
		},
	}

	// Helper function to normalize whitespace for comparison
	normalizeWhitespace := func(s string) string {
		// Replace all whitespace (tabs, newlines, spaces) with a single space
		s = strings.Join(strings.Fields(s), " ")
		// Remove spaces around braces
		s = strings.ReplaceAll(s, "{ ", "{")
		s = strings.ReplaceAll(s, " }", "}")
		return s
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			functions, err := g.GetFunctions(tt.input)

			assert.NoError(t, err)
			if len(tt.expected) != len(functions) {
				t.Errorf("expected %d functions but got %d", len(tt.expected), len(functions))
				for i, f := range functions {
					t.Logf("function %d: %s", i, f.Name)
				}
			}

			// Check each function matches expected
			for i, expected := range tt.expected {
				actual := functions[i]
				assert.Equal(t, expected.Name, actual.Name)
				assert.Equal(t, expected.IsExported, actual.IsExported)
				assert.Equal(t,
					normalizeWhitespace(expected.SourceCode),
					normalizeWhitespace(actual.SourceCode),
					"Function body should match after normalizing whitespace",
				)
			}
		})
	}
}
