/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package accounts

import (
	"github.com/datateamsix/email-sentinel/internal/appconfig"
)

// LoadConfigFromAppConfig converts AppConfig.Accounts to AccountConfig
func LoadConfigFromAppConfig(appCfg *appconfig.AppConfig) *AccountConfig {
	if appCfg == nil {
		return DefaultAccountConfig()
	}

	cfg := &AccountConfig{
		Enabled:            appCfg.Accounts.Enabled,
		MinConfidence:      appCfg.Accounts.Detection.MinConfidence,
		Categories:         appCfg.Accounts.Categories,
		DetectionKeywords:  appCfg.Accounts.Detection.Keywords,
		TrialAlerts:        make([]TrialAlert, 0),
	}

	// Convert trial alerts
	for _, alert := range appCfg.Accounts.TrialAlerts {
		cfg.TrialAlerts = append(cfg.TrialAlerts, TrialAlert{
			DaysBefore: alert.DaysBefore,
			Urgency:    alert.Urgency,
		})
	}

	// Apply defaults if not set
	if cfg.MinConfidence == 0 {
		cfg.MinConfidence = 0.7
	}

	if len(cfg.TrialAlerts) == 0 {
		cfg.TrialAlerts = []TrialAlert{
			{DaysBefore: 3, Urgency: "high"},
			{DaysBefore: 1, Urgency: "critical"},
		}
	}

	return cfg
}

// DefaultAccountConfig returns default account configuration
func DefaultAccountConfig() *AccountConfig {
	return &AccountConfig{
		Enabled:       true,
		MinConfidence: 0.7,
		TrialAlerts: []TrialAlert{
			{DaysBefore: 3, Urgency: "high"},
			{DaysBefore: 1, Urgency: "critical"},
		},
		Categories: map[string][]string{
			"streaming": {
				"Netflix", "Hulu", "Disney+", "Spotify", "Apple Music", "YouTube Premium",
			},
			"software": {
				"Adobe", "Microsoft 365", "GitHub", "Notion", "Grammarly", "Canva",
			},
			"cloud": {
				"AWS", "Google Cloud", "Dropbox", "iCloud", "OneDrive",
			},
			"productivity": {
				"ChatGPT", "Slack", "Zoom", "Asana", "Trello",
			},
		},
		DetectionKeywords: map[string][]string{
			"trial": {
				"free trial", "trial period", "trial started", "trial membership",
				"start your trial", "trial expires", "trial ends",
			},
			"subscription": {
				"subscription renewed", "payment successful", "subscription confirmed",
				"monthly subscription", "annual subscription", "recurring payment", "auto-renew",
			},
			"account_created": {
				"welcome to", "account created", "verify your email", "confirm your account",
				"registration successful", "account activated",
			},
			"cancellation": {
				"subscription cancelled", "subscription canceled", "membership ended",
				"auto-renew disabled", "will not be charged",
			},
		},
	}
}
