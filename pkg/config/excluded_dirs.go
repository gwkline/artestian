package config

import (
	"path/filepath"
	"strings"
)

// GetExcludedDirs returns the normalized excluded directory paths
func (c *Config) GetExcludedDirs() []string {
	if len(c.Settings.ExcludedDirs) == 0 {
		return nil
	}

	// Get the absolute path of the root directory
	rootDir := c.GetRootDir()
	excludedDirs := make([]string, len(c.Settings.ExcludedDirs))

	for i, dir := range c.Settings.ExcludedDirs {
		// If the path starts with ./, use the root directory as base
		if strings.HasPrefix(dir, "./") {
			excludedDirs[i] = filepath.Join(rootDir, strings.TrimPrefix(dir, "./"))
		} else {
			// For paths without ./, resolve them relative to the config base path
			excludedDirs[i] = c.resolveFilePath(dir)
		}
	}

	return excludedDirs
}
