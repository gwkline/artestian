package typescript

import (
	"strings"
	"testing"

	"github.com/gwkline/artestian/types"
	"github.com/stretchr/testify/assert"
)

func TestTypeScriptSupport_GetFunctions(t *testing.T) {
	ts := NewTypeScriptSupport()

	tests := []struct {
		name     string
		input    string
		expected []types.Function
	}{
		{
			name: "regular functions",
			input: `
						function hello() { return "hello" }
						export function world() { return "world" }
					`,
			expected: []types.Function{
				{
					Name:       "hello",
					SourceCode: `function hello() { return "hello" }`,
					IsExported: false,
				},
				{
					Name:       "world",
					SourceCode: `export function world() { return "world" }`,
					IsExported: true,
				},
			},
		},
		{
			name: "arrow functions",
			input: `
						const arrowFn = () => { return "arrow" }
						export const exportedArrow = () => { return "exported" }
					`,
			expected: []types.Function{
				{
					Name:       "arrowFn",
					SourceCode: `const arrowFn = () => { return "arrow" }`,
					IsExported: false,
				},
				{
					Name:       "exportedArrow",
					SourceCode: `export const exportedArrow = () => { return "exported" }`,
					IsExported: true,
				},
			},
		},
		{
			name: "mixed function types",
			input: `
						function regular() { return "regular" }
						const arrow = () => { return "arrow" }
						export function exportedRegular() { return "exported" }
						export const exportedArrow = () => { return "exported arrow" }
					`,
			expected: []types.Function{
				{
					Name:       "regular",
					SourceCode: `function regular() { return "regular" }`,
					IsExported: false,
				},
				{
					Name:       "arrow",
					SourceCode: `const arrow = () => { return "arrow" }`,
					IsExported: false,
				},
				{
					Name:       "exportedRegular",
					SourceCode: `export function exportedRegular() { return "exported" }`,
					IsExported: true,
				},
				{
					Name:       "exportedArrow",
					SourceCode: `export const exportedArrow = () => { return "exported arrow" }`,
					IsExported: true,
				},
			},
		},
		{
			name: "function with comments",
			input: `
						// Regular function with comment
						function simpleFunction() {
							return "simple";
						}
					`,
			expected: []types.Function{
				{
					Name: "simpleFunction",
					SourceCode: `function simpleFunction() {
							return "simple";
						}`,
					IsExported: false,
				},
			},
		},
		{
			name: "function with nested braces",
			input: `
						function complexFunction() {
							if (true) {
								return {
									nested: {
										data: "test"
									}
								};
							}
							return null;
						}
					`,
			expected: []types.Function{
				{
					Name: "complexFunction",
					SourceCode: `function complexFunction() {
							if (true) {
								return {
									nested: {
										data: "test"
									}
								};
							}
							return null;
						}`,
					IsExported: false,
				},
			},
		},
		{
			name: "async function",
			input: `
						async function processData(input: string) {
							return await input.toUpperCase();
						}
					`,
			expected: []types.Function{
				{
					Name: "processData",
					SourceCode: `async function processData(input: string) {
							return await input.toUpperCase();
						}`,
					IsExported: false,
				},
			},
		},
		{
			name: "split export keyword",
			input: `
						export
						function splitExport() {
							return "split";
						}
					`,
			expected: []types.Function{
				{
					Name: "splitExport",
					SourceCode: `export
						function splitExport() {
							return "split";
						}`,
					IsExported: true,
				},
			},
		},
		{
			name: "arrow function with type annotation",
			input: `
						export const typedArrow: () => boolean = () => {
							return true;
						}
					`,
			expected: []types.Function{
				{
					Name: "typedArrow",
					SourceCode: `export const typedArrow: () => boolean = () => {
							return true;
						}`,
					IsExported: true,
				},
			},
		},
		{
			name: "complex mixed example",
			input: `// Regular function with comment
		function simpleFunction() {
			return "simple";
		}

		/* Multiline comment for
			exported function */
		export function exportedFunction() {
			const x = 5;
			return x * 2;
		}

		// Arrow function with parameters
		const arrowWithParams = (x: number, y: string) => {
			console.log(y);
			return x + 1;
		}

		// Exported arrow function with type annotation
		export const typedArrow: () => boolean = () => {
			return true;
		}

		/**
			* JSDoc comment for async function
			* @param {string} input
			* @returns {Promise<string>}
			*/
		async function processData(input: string) {
			return await input.toUpperCase();
		}

		// Function with complex body
		function complexFunction() {
			if (true) {
				return {
					nested: {
						data: "test"
					}
				};
			}
			return null;
		}

		// split export
		export // Split export keyword
		function splitExport() {
			return "split";
		}

		// Exported arrow function with complex types
		export const genericArrow = <T extends object>(data: T) => {
			return {
				...data,
				timestamp: Date.now()
			};
		}`,
			expected: []types.Function{
				{
					Name: "simpleFunction",
					SourceCode: `function simpleFunction() {
							return "simple";
						}`,
					IsExported: false,
				},
				{
					Name: "exportedFunction",
					SourceCode: `export function exportedFunction() {
							const x = 5;
							return x * 2;
						}`,
					IsExported: true,
				},
				{
					Name: "arrowWithParams",
					SourceCode: `const arrowWithParams = (x: number, y: string) => {
							console.log(y);
							return x + 1;
						}`,
					IsExported: false,
				},
				{
					Name: "typedArrow",
					SourceCode: `export const typedArrow: () => boolean = () => {
							return true;
						}`,
					IsExported: true,
				},
				{
					Name: "processData",
					SourceCode: `async function processData(input: string) {
							return await input.toUpperCase();
						}`,
					IsExported: false,
				},
				{
					Name: "complexFunction",
					SourceCode: `function complexFunction() {
							if (true) {
								return {
									nested: {
										data: "test"
									}
								};
							}
							return null;
						}`,
					IsExported: false,
				},
				{
					Name: "splitExport",
					SourceCode: `export function splitExport() {
							return "split";
						}`,
					IsExported: true,
				},
				{
					Name: "genericArrow",
					SourceCode: `export const genericArrow = <T extends object>(data: T) => {
							return {
								...data,
								timestamp: Date.now()
							};
						}`,
					IsExported: true,
				},
			},
		},
		{
			name: "arrow function with complex types",
			input: `import { someTable } from "@/drizzle/schema"
import { drizzleDb } from "@/prisma/drizzle/schema"
import { and, eq } from "drizzle-orm"

export const updateSomeTable = async ({
  teamId,
  someTableId,
  someTableHeader,
  someTableDescription,
}: {
  teamId: string
  someTableId: string
  someTableHeader: string
  someTableDescription: string
}): Promise<void> => {
  await drizzleDb
    .update(someTable)
    .set({
      header: someTableHeader,
      description: someTableDescription,
    })
    .where(and(eq(someTable.id, someTableId), eq(someTable.teamId, teamId)))
}`,
			expected: []types.Function{
				{
					Name: "updateSomeTable",
					SourceCode: `export const updateSomeTable = async ({
  teamId,
  someTableId,
  someTableHeader,
  someTableDescription,
}: {
  teamId: string
  someTableId: string
  someTableHeader: string
  someTableDescription: string
}): Promise<void> => {
  await drizzleDb
    .update(someTable)
    .set({
      header: someTableHeader,
      description: someTableDescription,
    })
    .where(and(eq(someTable.id, someTableId), eq(someTable.teamId, teamId)))
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
			functions, err := ts.GetFunctions(tt.input)

			assert.NoError(t, err)
			if len(tt.expected) != len(functions) {
				t.Errorf("expected %d functions but got %d", len(tt.expected), len(functions))
				for i, f := range functions {
					t.Logf("function %d: %s", i, f.Name)
				}
				return // Early return to avoid index out of bounds
			}

			// Check each function matches expected
			for i, expected := range tt.expected {
				if i >= len(functions) {
					t.Errorf("missing expected function: %s", expected.Name)
					continue
				}
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
