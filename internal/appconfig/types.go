/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package appconfig

import (
	"time"
)

// AppConfig represents the unified application configuration
// This replaces the previous separate configs (ai-config.yaml, rules.yaml, otp_rules.yaml)
type AppConfig struct {
	Monitoring    MonitoringConfig    `yaml:"monitoring"`
	AISummary     AISummaryConfig     `yaml:"ai_summary"`
	Priority      PriorityConfig      `yaml:"priority"`
	OTP           OTPConfig           `yaml:"otp"`
	Accounts      AccountsConfig      `yaml:"accounts"`
	Notifications NotificationsConfig `yaml:"notifications"`
}

// MonitoringConfig holds email monitoring settings
type MonitoringConfig struct {
	PollingInterval int              `yaml:"polling_interval"` // seconds
	Database        DatabaseConfig   `yaml:"database"`
}

// DatabaseConfig holds database settings
type DatabaseConfig struct {
	WALMode         bool   `yaml:"wal_mode"`
	CleanupInterval string `yaml:"cleanup_interval"` // duration string like "1h", "0" to disable
}

// ==============================================================================
// AI Summary Configuration
// ==============================================================================

// AISummaryConfig holds AI-powered email summary settings
type AISummaryConfig struct {
	Enabled   bool                       `yaml:"enabled"`
	Provider  string                     `yaml:"provider"` // "gemini", "claude", "openai"
	Providers AIProvidersConfig          `yaml:"providers"`
	Cache     CacheConfig                `yaml:"cache"`
	Prompt    PromptConfig               `yaml:"prompt"`
}

// AIProvidersConfig holds settings for all AI providers
type AIProvidersConfig struct {
	Gemini GeminiProviderConfig `yaml:"gemini"`
	Claude ClaudeProviderConfig `yaml:"claude"`
	OpenAI OpenAIProviderConfig `yaml:"openai"`
}

// GeminiProviderConfig holds Google Gemini settings
type GeminiProviderConfig struct {
	Model       string          `yaml:"model"`
	Endpoint    string          `yaml:"endpoint"`
	MaxTokens   int             `yaml:"max_tokens"`
	Temperature float64         `yaml:"temperature"`
	RateLimit   RateLimitConfig `yaml:"rate_limit"`
}

// ClaudeProviderConfig holds Anthropic Claude settings
type ClaudeProviderConfig struct {
	Model       string          `yaml:"model"`
	Endpoint    string          `yaml:"endpoint"`
	MaxTokens   int             `yaml:"max_tokens"`
	Temperature float64         `yaml:"temperature"`
	RateLimit   RateLimitConfig `yaml:"rate_limit"`
}

// OpenAIProviderConfig holds OpenAI settings
type OpenAIProviderConfig struct {
	Model       string          `yaml:"model"`
	Endpoint    string          `yaml:"endpoint"`
	MaxTokens   int             `yaml:"max_tokens"`
	Temperature float64         `yaml:"temperature"`
	RateLimit   RateLimitConfig `yaml:"rate_limit"`
}

// RateLimitConfig controls API usage limits
type RateLimitConfig struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
	RequestsPerDay    int `yaml:"requests_per_day"`
}

// CacheConfig controls summary caching
type CacheConfig struct {
	Enabled bool   `yaml:"enabled"`
	TTL     string `yaml:"ttl"`      // duration string like "24h"
	MaxSize int    `yaml:"max_size"` // max number of cached summaries
}

// PromptConfig holds customizable AI prompts
type PromptConfig struct {
	System    string                 `yaml:"system"`
	Templates map[string]string      `yaml:"templates"`
}

// ==============================================================================
// Priority Rules Configuration
// ==============================================================================

// PriorityConfig defines rules for marking emails as high priority
type PriorityConfig struct {
	UrgentKeywords []string `yaml:"urgent_keywords"`
	VIPSenders     []string `yaml:"vip_senders"`
	VIPDomains     []string `yaml:"vip_domains"`
}

// ==============================================================================
// OTP Detection Configuration
// ==============================================================================

// OTPConfig holds OTP/2FA detection settings
type OTPConfig struct {
	Enabled          bool             `yaml:"enabled"`
	ExpiryDuration   string           `yaml:"expiry_duration"`   // duration string like "5m"
	MaxCodes         int              `yaml:"max_codes"`
	TrustedSenders   []string         `yaml:"trusted_senders"`
	TrustedDomains   []string         `yaml:"trusted_domains"`
	CustomPatterns   []CustomPattern  `yaml:"custom_patterns"`
	TriggerPhrases   []string         `yaml:"trigger_phrases"`
	Clipboard        ClipboardConfig  `yaml:"clipboard"`
}

// CustomPattern represents a custom OTP detection pattern
type CustomPattern struct {
	Pattern     string `yaml:"pattern"`
	Description string `yaml:"description"`
	Confidence  string `yaml:"confidence"` // "high", "medium", "low"
}

// ClipboardConfig controls clipboard integration
type ClipboardConfig struct {
	AutoCopy   bool   `yaml:"auto_copy"`
	ClearAfter string `yaml:"clear_after"` // duration string like "30s"
}

// ==============================================================================
// Digital Accounts Configuration
// ==============================================================================

// AccountsConfig holds digital account tracking settings
type AccountsConfig struct {
	Enabled      bool                       `yaml:"enabled"`
	TrialAlerts  []TrialAlert               `yaml:"trial_alerts"`
	Detection    AccountDetectionConfig     `yaml:"detection"`
	Categories   map[string][]string        `yaml:"categories"`
}

// TrialAlert defines when to alert before trial expiration
type TrialAlert struct {
	DaysBefore int    `yaml:"days_before"`
	Urgency    string `yaml:"urgency"` // "low", "high", "critical"
}

// AccountDetectionConfig controls account detection behavior
type AccountDetectionConfig struct {
	MinConfidence float64                `yaml:"min_confidence"`
	Keywords      map[string][]string    `yaml:"keywords"`
}

// ==============================================================================
// Notifications Configuration
// ==============================================================================

// NotificationsConfig controls notification behavior
type NotificationsConfig struct {
	Desktop     DesktopNotifConfig `yaml:"desktop"`
	Mobile      MobileNotifConfig  `yaml:"mobile"`
	QuietHours  QuietHoursConfig   `yaml:"quiet_hours"`
	WeekendMode string             `yaml:"weekend_mode"` // "normal", "quiet", "disabled"
}

// DesktopNotifConfig controls desktop notifications
type DesktopNotifConfig struct {
	Enabled  bool `yaml:"enabled"`
	Duration int  `yaml:"duration"` // seconds
	Sound    bool `yaml:"sound"`
}

// MobileNotifConfig controls mobile notifications (via ntfy.sh)
type MobileNotifConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Topic    string `yaml:"topic"`
	Server   string `yaml:"server"`
	Priority int    `yaml:"priority"` // 1-5
}

// QuietHoursConfig defines quiet hours settings
type QuietHoursConfig struct {
	Start       string `yaml:"start"`        // "HH:MM" format
	End         string `yaml:"end"`          // "HH:MM" format
	AllowUrgent bool   `yaml:"allow_urgent"` // Allow priority emails during quiet hours
}

// ==============================================================================
// Helper Methods
// ==============================================================================

// GetCleanupInterval returns the cleanup interval as a time.Duration
func (m *MonitoringConfig) GetCleanupInterval() (time.Duration, error) {
	if m.Database.CleanupInterval == "0" || m.Database.CleanupInterval == "" {
		return 0, nil
	}
	return time.ParseDuration(m.Database.CleanupInterval)
}

// GetOTPExpiryDuration returns the OTP expiry as a time.Duration
func (o *OTPConfig) GetOTPExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(o.ExpiryDuration)
}

// GetClearAfterDuration returns the clipboard clear duration as time.Duration
func (c *ClipboardConfig) GetClearAfterDuration() (time.Duration, error) {
	return time.ParseDuration(c.ClearAfter)
}

// GetCacheTTL returns the cache TTL as a time.Duration
func (c *CacheConfig) GetCacheTTL() (time.Duration, error) {
	return time.ParseDuration(c.TTL)
}
