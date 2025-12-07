package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/datateamsix/email-sentinel/internal/config"
	"github.com/datateamsix/email-sentinel/internal/gmail"
	"gopkg.in/yaml.v3"
)

// MessageMetadata contains the email data used for rule evaluation
type MessageMetadata struct {
	Sender  string
	Subject string
	Snippet string
	Body    string
}

// PriorityRules defines the conditions for marking emails as urgent (priority 1)
type PriorityRules struct {
	UrgentKeywords []string `yaml:"urgent_keywords"`
	VIPSenders     []string `yaml:"vip_senders"`
	VIPDomains     []string `yaml:"vip_domains"`
}

// NotificationSettings controls when and how notifications are sent
type NotificationSettings struct {
	QuietHoursStart string `yaml:"quiet_hours_start"` // e.g., "22:00"
	QuietHoursEnd   string `yaml:"quiet_hours_end"`   // e.g., "08:00"
	WeekendMode     string `yaml:"weekend_mode"`      // "normal", "quiet", "disabled"
}

// Rules represents the complete rules configuration
type Rules struct {
	PriorityRules        PriorityRules        `yaml:"priority_rules"`
	NotificationSettings NotificationSettings `yaml:"notification_settings"`
}

// DefaultRules returns a Rules struct with sensible defaults
func DefaultRules() *Rules {
	return &Rules{
		PriorityRules: PriorityRules{
			UrgentKeywords: []string{
				"urgent",
				"asap",
				"action required",
				"issue",
				"deadline",
				"invoice",
				"payment",
				"eod",
				"today",
			},
			VIPSenders: []string{
				// Users can add their important contacts
			},
			VIPDomains: []string{
				// Users can add important domains
			},
		},
		NotificationSettings: NotificationSettings{
			QuietHoursStart: "",        // Empty = disabled
			QuietHoursEnd:   "",        // Empty = disabled
			WeekendMode:     "normal",  // normal, quiet, disabled
		},
	}
}

// RulesPath returns the path where the rules.yaml file should be stored
func RulesPath() (string, error) {
	configDir, err := config.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "rules.yaml"), nil
}

// LoadRules loads the rules configuration from YAML file
// If the file doesn't exist, it creates one with defaults
func LoadRules(path string) (*Rules, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, create default
		rules := DefaultRules()
		if err := SaveRules(path, rules); err != nil {
			return nil, fmt.Errorf("failed to create default rules file: %w", err)
		}
		return rules, nil
	}

	// Read existing file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var rules Rules
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse rules YAML: %w", err)
	}

	return &rules, nil
}

// SaveRules writes the rules configuration to a YAML file
func SaveRules(path string, rules *Rules) error {
	data, err := yaml.Marshal(rules)
	if err != nil {
		return fmt.Errorf("failed to marshal rules: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write rules file: %w", err)
	}

	return nil
}

// EvaluatePriorityRules determines if a message should be marked as priority (1) or normal (0)
// Returns 1 if:
//   - Subject or snippet contains urgent keywords
//   - Sender matches VIP senders list
//   - Sender's domain matches VIP domains list
// Otherwise returns 0
func EvaluatePriorityRules(rules *Rules, msg MessageMetadata) int {
	if rules == nil {
		return 0 // No rules, default to normal priority
	}

	// Check urgent keywords in subject and snippet
	searchText := strings.ToLower(msg.Subject + " " + msg.Snippet + " " + msg.Body)
	for _, keyword := range rules.PriorityRules.UrgentKeywords {
		if strings.Contains(searchText, strings.ToLower(keyword)) {
			return 1 // Urgent keyword found
		}
	}

	// Extract sender email address
	senderEmail := gmail.GetFromAddress(msg.Sender)
	senderEmailLower := strings.ToLower(senderEmail)

	// Check VIP senders (exact match)
	for _, vipSender := range rules.PriorityRules.VIPSenders {
		if strings.ToLower(vipSender) == senderEmailLower {
			return 1 // VIP sender
		}
	}

	// Check VIP domains
	senderDomain := gmail.GetFromDomain(msg.Sender)
	senderDomainLower := strings.ToLower(senderDomain)

	for _, vipDomain := range rules.PriorityRules.VIPDomains {
		if strings.ToLower(vipDomain) == senderDomainLower {
			return 1 // VIP domain
		}
	}

	return 0 // Normal priority
}

// IsQuietTime checks if the current time falls within quiet hours
// Returns true if notifications should be suppressed
func (r *Rules) IsQuietTime() bool {
	// TODO: Implement time-based quiet hours checking
	// This would parse QuietHoursStart/End and compare with current time
	// For now, returns false (quiet hours disabled)
	return false
}

// ShouldNotifyOnWeekend checks if notifications should be sent on weekends
// based on the weekend_mode setting
func (r *Rules) ShouldNotifyOnWeekend() bool {
	// TODO: Implement weekend checking
	// "normal" = notify as usual
	// "quiet" = only notify for priority 1
	// "disabled" = no notifications
	return true
}
