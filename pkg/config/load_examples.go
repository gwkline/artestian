package config

import (
	"fmt"
	"os"

	"github.com/gwkline/artestian/types"
)

// LoadExamples loads the test examples from the configuration
func (c *Config) LoadExamples() ([]types.TestExample, error) {
	examples := make([]types.TestExample, 0, len(c.Examples))

	for _, ex := range c.Examples {
		fullPath := c.resolveFilePath(ex.FilePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read example file %s: %w", fullPath, err)
		}

		examples = append(examples, types.TestExample{
			Name:        ex.Name,
			Type:        types.TestType(ex.Type),
			SourceCode:  string(content),
			Description: ex.Description,
		})
	}

	return examples, nil
}
