package config

import (
	"fmt"
	"os"

	"github.com/gwkline/artestian/types"
)

// LoadContextFiles loads all context files specified in the configuration
func (c *Config) LoadContextFiles() ([]types.ContextFile, error) {
	if len(c.Context.Files) == 0 {
		return nil, nil
	}

	files := make([]types.ContextFile, 0, len(c.Context.Files))
	for _, file := range c.Context.Files {
		fullPath := c.resolveFilePath(file.Path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read context file %s: %w", fullPath, err)
		}

		files = append(files, types.ContextFile{
			Path:        file.Path,
			Content:     string(content),
			Description: file.Description,
			Type:        file.Type,
		})
	}

	return files, nil
}
