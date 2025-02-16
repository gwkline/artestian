package config

// GetRootDir returns the configured default test directory, defaulting to current directory
func (c *Config) GetRootDir() string {
	if c.Settings.DefaultTestDirectory == "" {
		return "."
	}
	return c.resolveFilePath(c.Settings.DefaultTestDirectory)
}
