/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package otp

import (
	"regexp"
	"time"
)

// OTPResult represents the result of OTP detection
type OTPResult struct {
	Code       string    // The extracted OTP code
	Confidence float64   // Confidence score (0.0 to 1.0)
	Source     string    // Where the code was found: "subject", "body", or "snippet"
	Pattern    string    // Name of the pattern that matched
	ExpiresAt  time.Time // When this OTP code expires
}

// OTPPattern represents a pattern for matching OTP codes
type OTPPattern struct {
	Name         string           // Pattern identifier
	Regex        *regexp.Regexp   // Compiled regex pattern
	Confidence   float64          // Base confidence score (0.0 to 1.0)
	CaptureGroup int              // Which regex group contains the code
	Validator    func(string) bool // Optional validator function
}

// OTPRules represents the configuration for OTP detection
type OTPRules struct {
	Enabled              bool            // Enable/disable OTP detection
	ExpiryDuration       time.Duration   // How long codes remain valid
	ConfidenceThreshold  float64         // Minimum confidence to accept (0.0 to 1.0)
	AutoCopy             bool            // Auto-copy to clipboard
	AutoClearDuration    time.Duration   // Auto-clear clipboard duration
	EnableSecureClipboard bool           // Enable secure clipboard features
	CustomPatterns       []CustomPattern // User-defined patterns
	TrustedSenders       []string        // Email domains/addresses that boost confidence
	BlockedPatterns      []string        // Patterns to never match (e.g., invoice numbers)
	MaxProcessingTime    time.Duration   // Maximum time for detection
}

// CustomPattern represents a user-defined OTP pattern
type CustomPattern struct {
	Name       string  // Pattern name
	Regex      string  // Regex pattern (will be compiled)
	Confidence float64 // Base confidence score
}

// DetectionContext contains the context for OTP detection
type DetectionContext struct {
	Subject string // Email subject
	Body    string // Email body (may be empty if not available)
	Snippet string // Email snippet/preview
	Sender  string // Sender email address
}

// ValidationResult represents the result of OTP validation
type ValidationResult struct {
	IsValid bool   // Whether the code is valid
	Reason  string // Reason for validation failure (if any)
}
