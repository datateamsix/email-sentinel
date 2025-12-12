/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package accounts

import (
	"regexp"
	"time"
)

// DetectionResult represents the result of account detection
type DetectionResult struct {
	ServiceName    string     // Name of the service (e.g., "Netflix", "Adobe")
	EmailAddress   string     // Email address used for the account
	AccountType    string     // "trial", "paid", "free"
	PriceMonthly   float64    // Monthly price (normalized)
	TrialEndDate   *time.Time // When trial expires (if applicable)
	CancelURL      string     // Cancellation URL (if detected)
	Category       string     // Service category (streaming, software, cloud, productivity)
	Confidence     float64    // Detection confidence score (0.0 to 1.0)
	GmailMessageID string     // Gmail message ID for reference
}

// DetectionPattern represents a pattern for matching account-related emails
type DetectionPattern struct {
	Name         string           // Pattern identifier (e.g., "trial_start", "subscription_renewal")
	Type         string           // Account type: "trial", "paid", "free", "cancellation"
	Keywords     []string         // Keywords that trigger this pattern
	ServiceRegex *regexp.Regexp   // Regex to extract service name
	PriceRegex   *regexp.Regexp   // Regex to extract price
	DateRegex    *regexp.Regexp   // Regex to extract trial end date
	Confidence   float64          // Base confidence score (0.0 to 1.0)
	Category     string           // Service category hint
}

// DetectionContext contains the email context for account detection
type DetectionContext struct {
	Subject  string // Email subject
	Body     string // Email body (may be empty if not available)
	Snippet  string // Email snippet/preview
	Sender   string // Sender email address
	ToEmail  string // Recipient email address (the user's email)
	ReceivedDate time.Time // When email was received
	MessageID string // Gmail message ID
}

// AccountConfig represents the configuration for account detection
type AccountConfig struct {
	Enabled            bool          // Enable/disable account detection
	MinConfidence      float64       // Minimum confidence threshold (0.0 to 1.0)
	TrialAlerts        []TrialAlert  // Trial expiration alerts configuration
	Categories         map[string][]string // Service categories
	DetectionKeywords  map[string][]string // Keywords for detection by type
}

// TrialAlert represents a trial expiration alert configuration
type TrialAlert struct {
	DaysBefore int    // Days before expiration to alert
	Urgency    string // Alert urgency level: "low", "high", "critical"
}

// PriceInfo represents extracted price information
type PriceInfo struct {
	Amount   float64 // Price amount
	Currency string  // Currency code (e.g., "USD")
	Period   string  // Billing period: "monthly", "annual", "one-time"
}

// ServiceCategory represents service categorization
type ServiceCategory struct {
	Name     string   // Category name (streaming, software, cloud, productivity)
	Services []string // Service names in this category
}
