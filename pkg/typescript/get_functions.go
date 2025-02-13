package typescript

import (
	"regexp"
	"strings"

	"github.com/gwkline/artestian/types"
)

func (ts *TypeScriptSupport) GetFunctions(sourceCode string) ([]types.Function, error) {
	var functions []types.Function

	// Build regex pattern from components
	patterns := []string{
		`(?m)`,                                 // Multiline mode
		`(?:^|\n)\s*`,                          // Start of line or newline, followed by whitespace
		`(?:(export)\s+)?`,                     // Optional export keyword (captured in group 1)
		`(?:async\s+)?`,                        // Optional async keyword
		`(?:`,                                  // Start of main function pattern group
		`function\s+(\w+)`,                     // Named function declaration (captured in group 2)
		`|`,                                    // OR
		`const\s+(\w+)`,                        // Const declaration (captured in group 3)
		`\s*(?::\s*(?:[^=;]|=>|\{[^}]*\})*?)?`, // Optional type annotation
		`\s*=\s*`,                              // Assignment
		`(?:async\s+)?`,                        // Optional async keyword for arrow function
		`(?:`,                                  // Start of function implementation group
		`(?:<[^>]+>\s*)?`,                      // Optional generic type parameters
		`(?:function|\([^)]*\)\s*(?::\s*[^{]*?)?\s*=>)`, // Function keyword or arrow function with return type
		`)`,
		`)`,
		`\s*[^{]*{`, // Any non-brace chars until opening brace
	}

	functionPattern := strings.Join(patterns, "")
	re := regexp.MustCompile(functionPattern)

	// First find all potential function starts
	matches := re.FindAllStringIndex(sourceCode, -1)

	for _, match := range matches {
		startPos := match[0]

		// Look for export keyword before the function
		preContext := sourceCode[max(0, startPos-50):startPos]
		isExported := strings.Contains(preContext, "export") ||
			re.FindStringSubmatch(sourceCode[startPos:match[1]])[1] != ""

		// Find matching closing brace
		braceCount := 1
		endPos := match[1]

		for i := endPos; i < len(sourceCode); i++ {
			if sourceCode[i] == '{' {
				braceCount++
			} else if sourceCode[i] == '}' {
				braceCount--
				if braceCount == 0 {
					endPos = i + 1
					break
				}
			}
		}

		// Get the full function text including any preceding export
		fullMatch := sourceCode[startPos:endPos]
		if isExported && !strings.HasPrefix(strings.TrimSpace(fullMatch), "export") {
			if idx := strings.LastIndex(preContext, "export"); idx != -1 {
				fullMatch = "export " + strings.TrimSpace(sourceCode[startPos:endPos])
			}
		}

		// Extract function name
		nameMatch := re.FindStringSubmatch(fullMatch)

		if len(nameMatch) < 4 {
			continue
		}

		name := nameMatch[2]
		if name == "" {
			name = nameMatch[3] // Use const name if it's an arrow function
		}

		if name == "" {
			continue
		}

		functions = append(functions, types.Function{
			Name:       name,
			SourceCode: strings.TrimSpace(fullMatch),
			IsExported: isExported,
		})
	}

	return functions, nil
}

// Helper function for logging
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// max returns the larger of x or y
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
