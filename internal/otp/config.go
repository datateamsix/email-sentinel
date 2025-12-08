/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package otp

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// OTPRulesYAML represents the YAML structure for OTP rules
type OTPRulesYAML struct {
	Enabled             bool              `yaml:"enabled"`
	ExpiryDuration      string            `yaml:"expiry_duration"`
	ConfidenceThreshold float64           `yaml:"confidence_threshold"`
	AutoCopy            bool              `yaml:"auto_copy_to_clipboard"`
	AutoClearDuration   string            `yaml:"clipboard_auto_clear"`
	CustomPatterns      []CustomPattern   `yaml:"custom_patterns"`
	TrustedSenders      []string          `yaml:"trusted_otp_senders"`
}

// LoadOTPRules loads OTP rules from a YAML file
func LoadOTPRules(path string) (*OTPRules, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read OTP rules file: %w", err)
	}

	var yamlRules OTPRulesYAML
	if err := yaml.Unmarshal(data, &yamlRules); err != nil {
		return nil, fmt.Errorf("failed to parse OTP rules YAML: %w", err)
	}

	// Parse durations
	expiryDuration, err := time.ParseDuration(yamlRules.ExpiryDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid expiry_duration: %w", err)
	}

	autoClearDuration, err := time.ParseDuration(yamlRules.AutoClearDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid clipboard_auto_clear: %w", err)
	}

	rules := &OTPRules{
		Enabled:              yamlRules.Enabled,
		ExpiryDuration:       expiryDuration,
		ConfidenceThreshold:  yamlRules.ConfidenceThreshold,
		AutoCopy:             yamlRules.AutoCopy,
		AutoClearDuration:    autoClearDuration,
		EnableSecureClipboard: yamlRules.AutoCopy, // Enable if auto-copy is on
		CustomPatterns:       yamlRules.CustomPatterns,
		TrustedSenders:       yamlRules.TrustedSenders,
		MaxProcessingTime:    500 * time.Millisecond,
	}

	return rules, nil
}

// SaveOTPRules saves OTP rules to a YAML file
func SaveOTPRules(path string, rules *OTPRules) error {
	yamlRules := OTPRulesYAML{
		Enabled:             rules.Enabled,
		ExpiryDuration:      rules.ExpiryDuration.String(),
		ConfidenceThreshold: rules.ConfidenceThreshold,
		AutoCopy:            rules.AutoCopy,
		AutoClearDuration:   rules.AutoClearDuration.String(),
		CustomPatterns:      rules.CustomPatterns,
		TrustedSenders:      rules.TrustedSenders,
	}

	data, err := yaml.Marshal(&yamlRules)
	if err != nil {
		return fmt.Errorf("failed to marshal OTP rules: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write OTP rules file: %w", err)
	}

	return nil
}

// DefaultOTPRules returns sensible default OTP rules
func DefaultOTPRules() *OTPRules {
	return &OTPRules{
		Enabled:              true,
		ExpiryDuration:       5 * time.Minute,
		ConfidenceThreshold:  0.7,
		AutoCopy:             false,
		AutoClearDuration:    2 * time.Minute,
		EnableSecureClipboard: false,
		CustomPatterns:       []CustomPattern{},
		TrustedSenders: []string{
			"accounts.google.com",
			"noreply@google.com",
			"amazon.com",
			"noreply@github.com",
			"github.com",
			"microsoft.com",
			"account.microsoft.com",
			"paypal.com",
			"venmo.com",
			"apple.com",
			"appleid.apple.com",
			"noreply@",
			"no-reply@",
			"@auth0.com",
			"@okta.com",
			"@twilio.com",
		},
		BlockedPatterns:   []string{},
		MaxProcessingTime: 500 * time.Millisecond,
	}
}

// MergeWithDefaults merges user rules with defaults for missing values
func MergeWithDefaults(userRules *OTPRules) *OTPRules {
	defaults := DefaultOTPRules()

	if userRules.ExpiryDuration == 0 {
		userRules.ExpiryDuration = defaults.ExpiryDuration
	}

	if userRules.ConfidenceThreshold == 0 {
		userRules.ConfidenceThreshold = defaults.ConfidenceThreshold
	}

	if userRules.AutoClearDuration == 0 {
		userRules.AutoClearDuration = defaults.AutoClearDuration
	}

	if userRules.MaxProcessingTime == 0 {
		userRules.MaxProcessingTime = defaults.MaxProcessingTime
	}

	if len(userRules.TrustedSenders) == 0 {
		userRules.TrustedSenders = defaults.TrustedSenders
	}

	return userRules
}

// GenerateExampleYAML generates an example OTP rules YAML configuration
func GenerateExampleYAML() string {
	return `# OTP/2FA Code Detection Rules
otp_rules:
  # Enable/disable OTP detection
  enabled: true

  # How long OTP codes remain active before expiring
  # Format: duration string (5m, 10m, 30m, 1h)
  expiry_duration: "5m"

  # Minimum confidence score to treat as valid OTP (0.0 to 1.0)
  # Higher = fewer false positives, but might miss some codes
  confidence_threshold: 0.7

  # Auto-copy most recent OTP to clipboard when detected
  auto_copy_to_clipboard: false

  # Auto-clear clipboard after this duration (security feature)
  clipboard_auto_clear: "2m"

  # Custom OTP patterns (advanced users)
  # Use Go regex syntax
  custom_patterns:
    # Example: Match codes like "Code: ABC123"
    - name: "custom_alphanumeric"
      regex: "Code:\\s*([A-Z0-9]{6})"
      confidence: 0.8

  # Sender domains known to send OTP codes
  # Emails from these domains get higher confidence scores
  trusted_otp_senders:
    - "accounts.google.com"
    - "amazon.com"
    - "github.com"
    - "microsoft.com"
    - "paypal.com"
    - "noreply@"
`
}
