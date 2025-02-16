package config

// GetLanguage returns the configured language, defaulting to typescript
func (c *Config) GetLanguage() string {
	if c.Settings.Language == "" {
		return string(TypeScript)
	}
	return c.Settings.Language
}
