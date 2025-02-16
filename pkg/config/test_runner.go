package config

// GetTestRunner returns the configured test runner, defaulting to the default for the language
func (c *Config) GetTestRunner() string {
	if c.Settings.TestRunner == "" {
		lang := Language(c.GetLanguage())
		return string(defaultTestRunner[lang])
	}
	return c.Settings.TestRunner
}
