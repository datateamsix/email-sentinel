/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package otp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Detector handles OTP code detection
type Detector struct {
	patterns []OTPPattern
	rules    *OTPRules
}

// NewDetector creates a new OTP detector with the given rules
func NewDetector(rules *OTPRules) (*Detector, error) {
	if err := validateOTPRules(rules); err != nil {
		return nil, fmt.Errorf("invalid OTP rules: %w", err)
	}

	detector := &Detector{
		patterns: GetBuiltInPatterns(),
		rules:    rules,
	}

	// Add custom patterns
	for _, customPattern := range rules.CustomPatterns {
		regex, err := regexp.Compile(customPattern.Regex)
		if err != nil {
			return nil, fmt.Errorf("invalid custom pattern %s: %w", customPattern.Name, err)
		}

		detector.patterns = append(detector.patterns, OTPPattern{
			Name:         customPattern.Name,
			Regex:        regex,
			Confidence:   customPattern.Confidence,
			CaptureGroup: 1,
			Validator:    nil,
		})
	}

	return detector, nil
}

// Detect finds OTP codes in the given context
func (d *Detector) Detect(ctx DetectionContext) *OTPResult {
	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(context.Background(), d.rules.MaxProcessingTime)
	defer cancel()

	// Search in priority order: subject, snippet, body
	sources := []struct {
		text   string
		source string
	}{
		{ctx.Subject, "subject"},
		{ctx.Snippet, "snippet"},
		{ctx.Body, "body"},
	}

	var bestMatch *OTPResult

	for _, src := range sources {
		if src.text == "" {
			continue
		}

		// Check for timeout
		select {
		case <-timeoutCtx.Done():
			return bestMatch // Return best match so far
		default:
		}

		result := d.detectInText(src.text, src.source, ctx.Sender, ctx.Subject)
		if result != nil {
			if bestMatch == nil || result.Confidence > bestMatch.Confidence {
				bestMatch = result
			}
		}
	}

	// Apply confidence threshold
	if bestMatch != nil && bestMatch.Confidence < d.rules.ConfidenceThreshold {
		return nil
	}

	return bestMatch
}

// detectInText searches for OTP codes in a specific text
func (d *Detector) detectInText(text string, source string, sender string, subject string) *OTPResult {
	var bestMatch *OTPResult

	for _, pattern := range d.patterns {
		matches := pattern.Regex.FindStringSubmatch(text)
		if len(matches) <= pattern.CaptureGroup {
			continue
		}

		code := matches[pattern.CaptureGroup]

		// Validate code if validator exists
		if pattern.Validator != nil && !pattern.Validator(code) {
			continue
		}

		// Normalize code
		code = NormalizeCode(code)

		// Check for false positives
		if IsLikelyFalsePositive(code, text) {
			continue
		}

		// Calculate confidence
		confidence := d.calculateConfidence(
			code,
			pattern.Confidence,
			text,
			sender,
			subject,
		)

		result := &OTPResult{
			Code:       code,
			Confidence: confidence,
			Source:     source,
			Pattern:    pattern.Name,
			ExpiresAt:  time.Now().Add(d.rules.ExpiryDuration),
		}

		if bestMatch == nil || result.Confidence > bestMatch.Confidence {
			bestMatch = result
		}
	}

	return bestMatch
}

// calculateConfidence computes the confidence score with adjustments
func (d *Detector) calculateConfidence(code string, baseConfidence float64, text string, sender string, subject string) float64 {
	confidence := baseConfidence

	// Boost for trusted senders
	for _, trustedSender := range d.rules.TrustedSenders {
		if strings.Contains(strings.ToLower(sender), strings.ToLower(trustedSender)) {
			confidence += 0.1
			break
		}
	}

	// Boost for OTP context in subject
	if HasOTPContext(subject) {
		confidence += 0.1
	}

	// Boost if code appears multiple times
	codeCount := strings.Count(strings.ToUpper(text), strings.ToUpper(code))
	if codeCount > 1 {
		confidence += 0.05
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// RegisterPattern adds a custom pattern to the detector
func (d *Detector) RegisterPattern(pattern OTPPattern) {
	d.patterns = append(d.patterns, pattern)
}

// DetectOTP is a convenience function for quick OTP detection
func DetectOTP(subject, body, snippet, sender string, rules *OTPRules) *OTPResult {
	detector, err := NewDetector(rules)
	if err != nil {
		return nil
	}

	ctx := DetectionContext{
		Subject: subject,
		Body:    body,
		Snippet: snippet,
		Sender:  sender,
	}

	return detector.Detect(ctx)
}

// ValidateOTP performs basic validation on an OTP code
func ValidateOTP(code string) bool {
	if len(code) < 4 || len(code) > 8 {
		return false
	}

	// Must be alphanumeric
	for _, c := range code {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}

	return true
}

// validateOTPRules validates OTP rules configuration
func validateOTPRules(rules *OTPRules) error {
	if rules == nil {
		return fmt.Errorf("rules cannot be nil")
	}

	if rules.ConfidenceThreshold < 0 || rules.ConfidenceThreshold > 1 {
		return fmt.Errorf("confidence threshold must be between 0 and 1, got %f", rules.ConfidenceThreshold)
	}

	if rules.ExpiryDuration < 0 {
		return fmt.Errorf("expiry duration cannot be negative")
	}

	if rules.MaxProcessingTime <= 0 {
		rules.MaxProcessingTime = 500 * time.Millisecond // Default
	}

	return nil
}
