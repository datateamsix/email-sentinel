package filter

// Filter represents an email filter rule
type Filter struct {
	Name    string   `yaml:"name"`
	From    []string `yaml:"from"`
	Subject []string `yaml:"subject"`
	Match   string   `yaml:"match"` // "any" or "all"
}

// Config represents the application configuration
type Config struct {
	PollingInterval int      `yaml:"polling_interval"`
	Filters         []Filter `yaml:"filters"`
	Notifications   struct {
		Desktop bool `yaml:"desktop"`
		Mobile  struct {
			Enabled   bool   `yaml:"enabled"`
			NtfyTopic string `yaml:"ntfy_topic"`
		} `yaml:"mobile"`
	} `yaml:"notifications"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	cfg := &Config{
		PollingInterval: 45,
		Filters:         []Filter{},
	}
	cfg.Notifications.Desktop = true
	cfg.Notifications.Mobile.Enabled = false
	return cfg
}
