package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/datateamsix/email-sentinel/internal/config"
	"gopkg.in/yaml.v3"
)

// Config represents the AI summary configuration
type Config struct {
	AISummary AISummaryConfig `yaml:"ai_summary"`
}

// AISummaryConfig holds all AI summarization settings
type AISummaryConfig struct {
	Enabled  bool              `yaml:"enabled"`
	Provider string            `yaml:"provider"` // "claude", "openai", "gemini"
	API      APIConfig         `yaml:"api"`
	Behavior BehaviorConfig    `yaml:"behavior"`
	RateLimit RateLimitConfig  `yaml:"rate_limit"`
	Prompt   PromptConfig      `yaml:"prompt"`
}

// APIConfig holds API settings for all providers
type APIConfig struct {
	Claude ClaudeConfig `yaml:"claude"`
	OpenAI OpenAIConfig `yaml:"openai"`
	Gemini GeminiConfig `yaml:"gemini"`
}

// ClaudeConfig holds Claude (Anthropic) API settings
type ClaudeConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// OpenAIConfig holds OpenAI API settings
type OpenAIConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// GeminiConfig holds Google Gemini API settings
type GeminiConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// BehaviorConfig controls summary generation behavior
type BehaviorConfig struct {
	MaxSummaryLength       int  `yaml:"max_summary_length"`
	PriorityOnly           bool `yaml:"priority_only"`
	EnableCache            bool `yaml:"enable_cache"`
	TimeoutSeconds         int  `yaml:"timeout_seconds"`
	RetryAttempts          int  `yaml:"retry_attempts"`
	IncludeInNotifications bool `yaml:"include_in_notifications"`
	ShowAIIcon             bool `yaml:"show_ai_icon"`
}

// RateLimitConfig controls API usage limits
type RateLimitConfig struct {
	MaxPerHour int `yaml:"max_per_hour"`
	MaxPerDay  int `yaml:"max_per_day"`
}

// PromptConfig holds customizable prompts
type PromptConfig struct {
	System       string `yaml:"system"`
	UserTemplate string `yaml:"user_template"`
}

// LoadConfig loads AI configuration from ai-config.yaml
func LoadConfig() (*Config, error) {
	// Try current directory first
	configPath := "ai-config.yaml"
	data, err := os.ReadFile(configPath)

	if err != nil {
		// Try config directory
		configDir, err := config.ConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get config directory: %w", err)
		}
		configPath = filepath.Join(configDir, "ai-config.yaml")
		data, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("ai-config.yaml not found: %w", err)
		}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse ai-config.yaml: %w", err)
	}

	// Load API keys from environment variables if not set in config
	cfg.loadAPIKeysFromEnv()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ai-config.yaml: %w", err)
	}

	return &cfg, nil
}

// loadAPIKeysFromEnv loads API keys from environment variables
func (c *Config) loadAPIKeysFromEnv() {
	// Claude API key
	if c.AISummary.API.Claude.APIKey == "" {
		if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
			c.AISummary.API.Claude.APIKey = key
		}
	}

	// OpenAI API key
	if c.AISummary.API.OpenAI.APIKey == "" {
		if key := os.Getenv("OPENAI_API_KEY"); key != "" {
			c.AISummary.API.OpenAI.APIKey = key
		}
	}

	// Gemini API key
	if c.AISummary.API.Gemini.APIKey == "" {
		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			c.AISummary.API.Gemini.APIKey = key
		}
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if !c.AISummary.Enabled {
		return nil // If disabled, no need to validate
	}

	provider := strings.ToLower(c.AISummary.Provider)
	if provider != "claude" && provider != "openai" && provider != "gemini" {
		return fmt.Errorf("invalid provider: %s (must be: claude, openai, or gemini)", c.AISummary.Provider)
	}

	// Check if API key is set for the selected provider
	var apiKey string
	switch provider {
	case "claude":
		apiKey = c.AISummary.API.Claude.APIKey
		if c.AISummary.API.Claude.Model == "" {
			return fmt.Errorf("claude model not specified")
		}
	case "openai":
		apiKey = c.AISummary.API.OpenAI.APIKey
		if c.AISummary.API.OpenAI.Model == "" {
			return fmt.Errorf("openai model not specified")
		}
	case "gemini":
		apiKey = c.AISummary.API.Gemini.APIKey
		if c.AISummary.API.Gemini.Model == "" {
			return fmt.Errorf("gemini model not specified")
		}
	}

	if apiKey == "" {
		return fmt.Errorf("API key not set for provider: %s", provider)
	}

	return nil
}

// GetProviderConfig returns the configuration for the active provider
func (c *Config) GetProviderConfig() interface{} {
	switch strings.ToLower(c.AISummary.Provider) {
	case "claude":
		return c.AISummary.API.Claude
	case "openai":
		return c.AISummary.API.OpenAI
	case "gemini":
		return c.AISummary.API.Gemini
	default:
		return nil
	}
}
