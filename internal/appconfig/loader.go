/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package appconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/datateamsix/email-sentinel/internal/config"
	"gopkg.in/yaml.v3"
)

// Load loads the unified app configuration
// It first tries to load app-config.yaml, and if not found, attempts to migrate
// from the old separate config files (ai-config.yaml, rules.yaml, otp_rules.yaml)
func Load() (*AppConfig, error) {
	// Try loading unified config first
	appConfig, err := loadUnifiedConfig()
	if err == nil {
		return appConfig, nil
	}

	// If unified config doesn't exist, try migration from old configs
	if os.IsNotExist(err) {
		fmt.Println("📦 Migrating from separate config files to unified app-config.yaml...")
		appConfig, migErr := migrateFromLegacyConfigs()
		if migErr != nil {
			// If migration fails, return default config
			fmt.Printf("⚠️  Migration failed: %v\n", migErr)
			fmt.Println("📝 Creating default configuration...")
			return DefaultConfig(), nil
		}

		// Save migrated config
		if saveErr := Save(appConfig); saveErr != nil {
			fmt.Printf("⚠️  Failed to save migrated config: %v\n", saveErr)
		} else {
			fmt.Println("✅ Successfully migrated to app-config.yaml")
		}

		return appConfig, nil
	}

	// Some other error occurred
	return nil, err
}

// loadUnifiedConfig loads the app-config.yaml file
func loadUnifiedConfig() (*AppConfig, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read app-config.yaml: %w", err)
	}

	// Parse YAML
	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse app-config.yaml: %w", err)
	}

	return &cfg, nil
}

// Save saves the app configuration to app-config.yaml
func Save(cfg *AppConfig) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if _, err := config.EnsureConfigDir(); err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write file with secure permissions (0600 - owner read/write only)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write app-config.yaml: %w", err)
	}

	return nil
}

// ConfigPath returns the path to the unified app-config.yaml file
func ConfigPath() (string, error) {
	configDir, err := config.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "app-config.yaml"), nil
}

// ConfigExists checks if app-config.yaml exists
func ConfigExists() bool {
	path, err := ConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// migrateFromLegacyConfigs attempts to migrate from old separate config files
// to the new unified app-config.yaml format
func migrateFromLegacyConfigs() (*AppConfig, error) {
	configDir, err := config.ConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Start with default config
	appConfig := DefaultConfig()

	// Track which files were successfully migrated
	migratedFiles := []string{}

	// 1. Try to load and migrate ai-config.yaml
	aiConfigPath := filepath.Join(configDir, "ai-config.yaml")
	if _, err := os.Stat(aiConfigPath); err == nil {
		if err := migrateAIConfig(aiConfigPath, appConfig); err != nil {
			fmt.Printf("⚠️  Failed to migrate ai-config.yaml: %v\n", err)
		} else {
			migratedFiles = append(migratedFiles, "ai-config.yaml")
		}
	}

	// 2. Try to load and migrate rules.yaml
	rulesPath := filepath.Join(configDir, "rules.yaml")
	if _, err := os.Stat(rulesPath); err == nil {
		if err := migrateRules(rulesPath, appConfig); err != nil {
			fmt.Printf("⚠️  Failed to migrate rules.yaml: %v\n", err)
		} else {
			migratedFiles = append(migratedFiles, "rules.yaml")
		}
	}

	// 3. Try to load and migrate otp_rules.yaml
	otpRulesPath := filepath.Join(configDir, "otp_rules.yaml")
	if _, err := os.Stat(otpRulesPath); err == nil {
		if err := migrateOTPRules(otpRulesPath, appConfig); err != nil {
			fmt.Printf("⚠️  Failed to migrate otp_rules.yaml: %v\n", err)
		} else {
			migratedFiles = append(migratedFiles, "otp_rules.yaml")
		}
	}

	// If no files were found to migrate, return an error
	if len(migratedFiles) == 0 {
		return nil, fmt.Errorf("no legacy config files found to migrate")
	}

	fmt.Printf("📦 Migrated from: %v\n", migratedFiles)
	return appConfig, nil
}

// migrateAIConfig migrates ai-config.yaml to the new format
func migrateAIConfig(path string, appConfig *AppConfig) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Use the old AI config structure
	var oldConfig struct {
		AISummary struct {
			Enabled  bool   `yaml:"enabled"`
			Provider string `yaml:"provider"`
			API      struct {
				Claude struct {
					Model       string  `yaml:"model"`
					MaxTokens   int     `yaml:"max_tokens"`
					Temperature float64 `yaml:"temperature"`
				} `yaml:"claude"`
				OpenAI struct {
					Model       string  `yaml:"model"`
					MaxTokens   int     `yaml:"max_tokens"`
					Temperature float64 `yaml:"temperature"`
				} `yaml:"openai"`
				Gemini struct {
					Model       string  `yaml:"model"`
					MaxTokens   int     `yaml:"max_tokens"`
					Temperature float64 `yaml:"temperature"`
				} `yaml:"gemini"`
			} `yaml:"api"`
			Behavior struct {
				EnableCache bool `yaml:"enable_cache"`
			} `yaml:"behavior"`
			RateLimit struct {
				MaxPerHour int `yaml:"max_per_hour"`
				MaxPerDay  int `yaml:"max_per_day"`
			} `yaml:"rate_limit"`
			Prompt struct {
				System       string `yaml:"system"`
				UserTemplate string `yaml:"user_template"`
			} `yaml:"prompt"`
		} `yaml:"ai_summary"`
	}

	if err := yaml.Unmarshal(data, &oldConfig); err != nil {
		return err
	}

	// Migrate to new structure
	appConfig.AISummary.Enabled = oldConfig.AISummary.Enabled
	appConfig.AISummary.Provider = oldConfig.AISummary.Provider

	// Migrate Claude config
	if oldConfig.AISummary.API.Claude.Model != "" {
		appConfig.AISummary.Providers.Claude.Model = oldConfig.AISummary.API.Claude.Model
		appConfig.AISummary.Providers.Claude.MaxTokens = oldConfig.AISummary.API.Claude.MaxTokens
		appConfig.AISummary.Providers.Claude.Temperature = oldConfig.AISummary.API.Claude.Temperature
	}

	// Migrate OpenAI config
	if oldConfig.AISummary.API.OpenAI.Model != "" {
		appConfig.AISummary.Providers.OpenAI.Model = oldConfig.AISummary.API.OpenAI.Model
		appConfig.AISummary.Providers.OpenAI.MaxTokens = oldConfig.AISummary.API.OpenAI.MaxTokens
		appConfig.AISummary.Providers.OpenAI.Temperature = oldConfig.AISummary.API.OpenAI.Temperature
	}

	// Migrate Gemini config
	if oldConfig.AISummary.API.Gemini.Model != "" {
		appConfig.AISummary.Providers.Gemini.Model = oldConfig.AISummary.API.Gemini.Model
		appConfig.AISummary.Providers.Gemini.MaxTokens = oldConfig.AISummary.API.Gemini.MaxTokens
		appConfig.AISummary.Providers.Gemini.Temperature = oldConfig.AISummary.API.Gemini.Temperature
	}

	// Migrate cache settings
	appConfig.AISummary.Cache.Enabled = oldConfig.AISummary.Behavior.EnableCache

	// Migrate rate limits (convert from hour/day to per-minute/per-day)
	if oldConfig.AISummary.RateLimit.MaxPerHour > 0 {
		// Convert per-hour to per-minute (rough approximation)
		perMinute := oldConfig.AISummary.RateLimit.MaxPerHour / 60
		if perMinute > 0 {
			appConfig.AISummary.Providers.Gemini.RateLimit.RequestsPerMinute = perMinute
			appConfig.AISummary.Providers.Claude.RateLimit.RequestsPerMinute = perMinute
			appConfig.AISummary.Providers.OpenAI.RateLimit.RequestsPerMinute = perMinute
		}
	}
	if oldConfig.AISummary.RateLimit.MaxPerDay > 0 {
		appConfig.AISummary.Providers.Gemini.RateLimit.RequestsPerDay = oldConfig.AISummary.RateLimit.MaxPerDay
		appConfig.AISummary.Providers.Claude.RateLimit.RequestsPerDay = oldConfig.AISummary.RateLimit.MaxPerDay
		appConfig.AISummary.Providers.OpenAI.RateLimit.RequestsPerDay = oldConfig.AISummary.RateLimit.MaxPerDay
	}

	// Migrate prompt settings
	if oldConfig.AISummary.Prompt.System != "" {
		appConfig.AISummary.Prompt.System = oldConfig.AISummary.Prompt.System
	}

	return nil
}

// migrateRules migrates rules.yaml to the new format
func migrateRules(path string, appConfig *AppConfig) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var oldRules struct {
		PriorityRules struct {
			UrgentKeywords []string `yaml:"urgent_keywords"`
			VIPSenders     []string `yaml:"vip_senders"`
			VIPDomains     []string `yaml:"vip_domains"`
		} `yaml:"priority_rules"`
		NotificationSettings struct {
			QuietHoursStart string `yaml:"quiet_hours_start"`
			QuietHoursEnd   string `yaml:"quiet_hours_end"`
			WeekendMode     string `yaml:"weekend_mode"`
		} `yaml:"notification_settings"`
	}

	if err := yaml.Unmarshal(data, &oldRules); err != nil {
		return err
	}

	// Migrate priority rules
	if len(oldRules.PriorityRules.UrgentKeywords) > 0 {
		appConfig.Priority.UrgentKeywords = oldRules.PriorityRules.UrgentKeywords
	}
	if len(oldRules.PriorityRules.VIPSenders) > 0 {
		appConfig.Priority.VIPSenders = oldRules.PriorityRules.VIPSenders
	}
	if len(oldRules.PriorityRules.VIPDomains) > 0 {
		appConfig.Priority.VIPDomains = oldRules.PriorityRules.VIPDomains
	}

	// Migrate notification settings
	if oldRules.NotificationSettings.QuietHoursStart != "" {
		appConfig.Notifications.QuietHours.Start = oldRules.NotificationSettings.QuietHoursStart
	}
	if oldRules.NotificationSettings.QuietHoursEnd != "" {
		appConfig.Notifications.QuietHours.End = oldRules.NotificationSettings.QuietHoursEnd
	}
	if oldRules.NotificationSettings.WeekendMode != "" {
		appConfig.Notifications.WeekendMode = oldRules.NotificationSettings.WeekendMode
	}

	return nil
}

// migrateOTPRules migrates otp_rules.yaml to the new format
func migrateOTPRules(path string, appConfig *AppConfig) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var oldOTPRules struct {
		Enabled           bool   `yaml:"enabled"`
		ExpiryDuration    string `yaml:"expiry_duration"`
		AutoCopy          bool   `yaml:"auto_copy_to_clipboard"`
		AutoClearDuration string `yaml:"clipboard_auto_clear"`
		CustomPatterns    []struct {
			Name       string  `yaml:"name"`
			Regex      string  `yaml:"regex"`
			Confidence float64 `yaml:"confidence"`
		} `yaml:"custom_patterns"`
		TrustedSenders []string `yaml:"trusted_otp_senders"`
	}

	if err := yaml.Unmarshal(data, &oldOTPRules); err != nil {
		return err
	}

	// Migrate OTP settings
	appConfig.OTP.Enabled = oldOTPRules.Enabled

	if oldOTPRules.ExpiryDuration != "" {
		appConfig.OTP.ExpiryDuration = oldOTPRules.ExpiryDuration
	}

	// Migrate clipboard settings
	appConfig.OTP.Clipboard.AutoCopy = oldOTPRules.AutoCopy
	if oldOTPRules.AutoClearDuration != "" {
		appConfig.OTP.Clipboard.ClearAfter = oldOTPRules.AutoClearDuration
	}

	// Migrate custom patterns
	if len(oldOTPRules.CustomPatterns) > 0 {
		appConfig.OTP.CustomPatterns = []CustomPattern{}
		for _, pattern := range oldOTPRules.CustomPatterns {
			confidence := "medium"
			if pattern.Confidence >= 0.8 {
				confidence = "high"
			} else if pattern.Confidence < 0.5 {
				confidence = "low"
			}

			appConfig.OTP.CustomPatterns = append(appConfig.OTP.CustomPatterns, CustomPattern{
				Pattern:     pattern.Regex,
				Description: pattern.Name,
				Confidence:  confidence,
			})
		}
	}

	// Migrate trusted senders
	if len(oldOTPRules.TrustedSenders) > 0 {
		appConfig.OTP.TrustedSenders = oldOTPRules.TrustedSenders
	}

	return nil
}

// DefaultConfig returns a new AppConfig with sensible defaults
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Monitoring: MonitoringConfig{
			PollingInterval: 45,
			Database: DatabaseConfig{
				WALMode:         true,
				CleanupInterval: "1h",
			},
		},
		AISummary: AISummaryConfig{
			Enabled:  false,
			Provider: "gemini",
			Providers: AIProvidersConfig{
				Gemini: GeminiProviderConfig{
					Model:       "gemini-2.0-flash-exp",
					Endpoint:    "https://generativelanguage.googleapis.com/v1beta/models",
					MaxTokens:   1000,
					Temperature: 0.3,
					RateLimit: RateLimitConfig{
						RequestsPerMinute: 15,
						RequestsPerDay:    1500,
					},
				},
				Claude: ClaudeProviderConfig{
					Model:       "claude-3-5-haiku-20241022",
					Endpoint:    "https://api.anthropic.com/v1/messages",
					MaxTokens:   1000,
					Temperature: 0.3,
					RateLimit: RateLimitConfig{
						RequestsPerMinute: 50,
						RequestsPerDay:    10000,
					},
				},
				OpenAI: OpenAIProviderConfig{
					Model:       "gpt-4o-mini",
					Endpoint:    "https://api.openai.com/v1/chat/completions",
					MaxTokens:   1000,
					Temperature: 0.3,
					RateLimit: RateLimitConfig{
						RequestsPerMinute: 60,
						RequestsPerDay:    10000,
					},
				},
			},
			Cache: CacheConfig{
				Enabled: true,
				TTL:     "24h",
				MaxSize: 1000,
			},
			Prompt: PromptConfig{
				System: "You are an expert email assistant. Analyze emails and provide:\n1. A concise 2-3 sentence summary\n2. Key questions that need answers (if any)\n3. Action items required (if any)\n\nBe direct and factual. Focus on what matters.",
				Templates: map[string]string{
					"meeting":      "Focus on time, participants, and agenda items.",
					"task":         "Extract deadlines, deliverables, and dependencies.",
					"notification": "Identify what changed and why it matters.",
				},
			},
		},
		Priority: PriorityConfig{
			UrgentKeywords: []string{
				"urgent", "asap", "immediate", "emergency", "critical",
				"deadline", "eod", "today", "now", "time sensitive",
				"action required", "action needed", "please review",
				"needs attention", "requires action", "response needed",
				"approval required", "issue", "problem", "error",
				"failed", "failure", "down", "outage", "incident",
				"alert", "warning", "invoice", "payment", "overdue",
				"past due", "billing", "important", "high priority",
				"priority 1", "p1", "escalation",
			},
			VIPSenders: []string{
				"boss@company.com",
				"ceo@company.com",
				"manager@company.com",
			},
			VIPDomains: []string{
				"costar.com",
				"costargroup.com",
				"mckinleyinc.com",
			},
		},
		OTP: OTPConfig{
			Enabled:        true,
			ExpiryDuration: "5m",
			MaxCodes:       50,
			TrustedSenders: []string{
				"noreply@accountprotection.microsoft.com",
				"account-security-noreply@accountprotection.microsoft.com",
				"no-reply@accounts.google.com",
				"account-update@account.apple.com",
				"noreply@github.com",
				"team@vercel.com",
				"noreply@gitlab.com",
				"notifications@bitbucket.org",
				"no-reply@aws.amazon.com",
				"azure-noreply@microsoft.com",
				"cloudplatform-noreply@google.com",
				"verify@twilio.com",
				"noreply@slack.com",
				"notifications@discord.com",
			},
			TrustedDomains: []string{
				"amazon.com",
				"paypal.com",
				"stripe.com",
				"linkedin.com",
				"twitter.com",
			},
			CustomPatterns: []CustomPattern{
				{
					Pattern:     `\b\d{6}\b`,
					Description: "6-digit numeric code",
					Confidence:  "high",
				},
				{
					Pattern:     `\b\d{8}\b`,
					Description: "8-digit numeric code",
					Confidence:  "medium",
				},
				{
					Pattern:     `\b[A-Z0-9]{6}\b`,
					Description: "6-character alphanumeric code",
					Confidence:  "medium",
				},
			},
			TriggerPhrases: []string{
				"verification code", "confirm your", "security code",
				"authentication code", "login code", "access code",
				"one-time password", "otp", "2fa", "two-factor",
				"verify your account", "confirm your email",
				"confirm your identity",
			},
			Clipboard: ClipboardConfig{
				AutoCopy:   false,
				ClearAfter: "30s",
			},
		},
		Notifications: NotificationsConfig{
			Desktop: DesktopNotifConfig{
				Enabled:  true,
				Duration: 10,
				Sound:    true,
			},
			Mobile: MobileNotifConfig{
				Enabled:  false,
				Topic:    "",
				Server:   "https://ntfy.sh",
				Priority: 4,
			},
			QuietHours: QuietHoursConfig{
				Start:       "",
				End:         "",
				AllowUrgent: true,
			},
			WeekendMode: "normal",
		},
	}
}
