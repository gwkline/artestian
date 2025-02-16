package config

// GetExcludedFiles returns the list of excluded file patterns
func (c *Config) GetExcludedFiles() []string {
	if len(c.Settings.ExcludedFiles) == 0 {
		return nil
	}
	return c.Settings.ExcludedFiles
}
