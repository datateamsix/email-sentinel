/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package otp

import (
	"regexp"
	"strconv"
	"strings"
)

// GetBuiltInPatterns returns the default OTP detection patterns
func GetBuiltInPatterns() []OTPPattern {
	return []OTPPattern{
		// High confidence patterns with context keywords
		{
			Name:         "code_keyword_6digit",
			Regex:        regexp.MustCompile(`(?i)(?:code|otp|token|pin)[\s:]*(\d{6})`),
			Confidence:   0.85,
			CaptureGroup: 1,
			Validator:    validateNumeric,
		},
		{
			Name:         "your_code_is",
			Regex:        regexp.MustCompile(`(?i)your\s+(?:code|otp|token|verification code|pin)\s+(?:is|:)[\s:]*([A-Z0-9]{4,8})`),
			Confidence:   0.90,
			CaptureGroup: 1,
			Validator:    validateAlphanumeric,
		},
		{
			Name:         "verification_code",
			Regex:        regexp.MustCompile(`(?i)verification\s+code[\s:]*(\d{6})`),
			Confidence:   0.85,
			CaptureGroup: 1,
			Validator:    validateNumeric,
		},
		{
			Name:         "use_code",
			Regex:        regexp.MustCompile(`(?i)use\s+(?:code|otp)[\s:]*([A-Z0-9]{4,8})`),
			Confidence:   0.80,
			CaptureGroup: 1,
			Validator:    validateAlphanumeric,
		},
		{
			Name:         "code_in_quotes",
			Regex:        regexp.MustCompile(`(?i)["']([0-9]{6})["']`),
			Confidence:   0.75,
			CaptureGroup: 1,
			Validator:    validateNumeric,
		},
		// Medium confidence patterns
		{
			Name:         "8_digit_numeric",
			Regex:        regexp.MustCompile(`\b(\d{8})\b`),
			Confidence:   0.50,
			CaptureGroup: 1,
			Validator:    validateNumeric,
		},
		{
			Name:         "6_digit_numeric",
			Regex:        regexp.MustCompile(`\b(\d{6})\b`),
			Confidence:   0.60,
			CaptureGroup: 1,
			Validator:    validateNumeric,
		},
		{
			Name:         "6_char_alphanumeric",
			Regex:        regexp.MustCompile(`\b([A-Z0-9]{6})\b`),
			Confidence:   0.55,
			CaptureGroup: 1,
			Validator:    validateAlphanumeric,
		},
		// Lower confidence patterns
		{
			Name:         "4_digit_pin",
			Regex:        regexp.MustCompile(`(?i)(?:pin|code)[\s:]*(\d{4})\b`),
			Confidence:   0.70,
			CaptureGroup: 1,
			Validator:    validateNumeric,
		},
		{
			Name:         "hyphenated_code",
			Regex:        regexp.MustCompile(`\b(\d{3}-\d{3})\b`),
			Confidence:   0.65,
			CaptureGroup: 1,
			Validator:    nil, // No validator, keep hyphens
		},
	}
}

// validateNumeric checks if the code is all numeric
func validateNumeric(code string) bool {
	_, err := strconv.Atoi(code)
	return err == nil
}

// validateAlphanumeric checks if the code is alphanumeric
func validateAlphanumeric(code string) bool {
	if len(code) == 0 {
		return false
	}

	for _, c := range code {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}

// IsLikelyFalsePositive checks if a code is likely a false positive
func IsLikelyFalsePositive(code string, context string) bool {
	context = strings.ToLower(context)

	// Check for invoice/order patterns
	falsePositiveKeywords := []string{
		"invoice", "order", "transaction", "receipt",
		"reference", "confirmation", "tracking",
		"phone", "fax", "ext", "extension",
	}

	for _, keyword := range falsePositiveKeywords {
		if strings.Contains(context, keyword) {
			return true
		}
	}

	// Check for sequential digits (123456, 654321)
	if isSequential(code) {
		return true
	}

	// Check for repeating digits (111111, 000000)
	if isRepeating(code) {
		return true
	}

	return false
}

// isSequential checks if the code consists of sequential digits
func isSequential(code string) bool {
	if len(code) < 4 {
		return false
	}

	// Check ascending
	isAscending := true
	for i := 1; i < len(code); i++ {
		if code[i] != code[i-1]+1 {
			isAscending = false
			break
		}
	}

	// Check descending
	isDescending := true
	for i := 1; i < len(code); i++ {
		if code[i] != code[i-1]-1 {
			isDescending = false
			break
		}
	}

	return isAscending || isDescending
}

// isRepeating checks if the code consists of repeating digits
func isRepeating(code string) bool {
	if len(code) < 4 {
		return false
	}

	first := code[0]
	for i := 1; i < len(code); i++ {
		if code[i] != first {
			return false
		}
	}

	return true
}

// HasOTPContext checks if text contains OTP-related keywords
func HasOTPContext(text string) bool {
	text = strings.ToLower(text)

	otpKeywords := []string{
		"otp", "code", "verification", "authenticate", "authentication",
		"2fa", "two-factor", "two factor", "security code", "access code",
		"confirmation code", "verify", "login code", "signin code",
	}

	for _, keyword := range otpKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}

// NormalizeCode removes common separators from codes
func NormalizeCode(code string) string {
	// Remove hyphens, spaces, underscores
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "_", "")
	code = strings.TrimSpace(code)

	return strings.ToUpper(code)
}
