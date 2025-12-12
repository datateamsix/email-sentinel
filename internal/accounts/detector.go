/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package accounts

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Detector handles account detection from emails
type Detector struct {
	patterns      []DetectionPattern
	minConfidence float64
	categories    map[string][]string
}

// NewDetector creates a new account detector
func NewDetector(minConfidence float64, categories map[string][]string) *Detector {
	return &Detector{
		patterns:      GetDefaultPatterns(),
		minConfidence: minConfidence,
		categories:    categories,
	}
}

// DetectAccount analyzes an email and attempts to detect account information
func (d *Detector) DetectAccount(ctx DetectionContext) (*DetectionResult, error) {
	// Combine all text for analysis
	fullText := ctx.Subject + " " + ctx.Snippet + " " + ctx.Body

	// Try each pattern
	for _, pattern := range d.patterns {
		if d.matchesPattern(fullText, pattern) {
			result := d.extractAccountInfo(ctx, pattern, fullText)
			if result != nil && result.Confidence >= d.minConfidence {
				return result, nil
			}
		}
	}

	return nil, nil // No account detected
}

// matchesPattern checks if the text contains keywords from the pattern
func (d *Detector) matchesPattern(text string, pattern DetectionPattern) bool {
	textLower := toLower(text)

	// Check if any keyword matches
	for _, keyword := range pattern.Keywords {
		if contains(textLower, toLower(keyword)) {
			return true
		}
	}

	return false
}

// extractAccountInfo extracts account information using the pattern
func (d *Detector) extractAccountInfo(ctx DetectionContext, pattern DetectionPattern, fullText string) *DetectionResult {
	result := &DetectionResult{
		AccountType:    pattern.Type,
		Confidence:     pattern.Confidence,
		GmailMessageID: ctx.MessageID,
		EmailAddress:   ctx.ToEmail,
	}

	// Extract service name
	if pattern.ServiceRegex != nil {
		if matches := pattern.ServiceRegex.FindStringSubmatch(fullText); len(matches) > 1 {
			result.ServiceName = strings.TrimSpace(matches[1])
		}
	}

	// If no service name found from regex, try to extract from sender
	if result.ServiceName == "" {
		result.ServiceName = extractServiceFromSender(ctx.Sender)
	}

	// Extract price if pattern has price regex
	if pattern.PriceRegex != nil {
		if price := d.extractPrice(fullText, pattern.PriceRegex); price > 0 {
			result.PriceMonthly = price
		}
	}

	// Extract trial end date if applicable
	if pattern.Type == "trial" && pattern.DateRegex != nil {
		if trialEnd := d.extractDate(fullText, pattern.DateRegex, ctx.ReceivedDate); trialEnd != nil {
			result.TrialEndDate = trialEnd
		}
	}

	// Extract cancel URL
	result.CancelURL = extractCancelURL(fullText)

	// Determine category
	if result.ServiceName != "" {
		result.Category = DetermineCategory(result.ServiceName)
	}

	// Boost confidence if we have good data
	if result.ServiceName != "" {
		result.Confidence += 0.05
	}
	if result.PriceMonthly > 0 {
		result.Confidence += 0.05
	}
	if result.TrialEndDate != nil {
		result.Confidence += 0.05
	}

	// Cap confidence at 1.0
	if result.Confidence > 1.0 {
		result.Confidence = 1.0
	}

	// Only return if we have at least a service name
	if result.ServiceName == "" {
		return nil
	}

	return result
}

// extractPrice attempts to extract a price from text using the pattern
func (d *Detector) extractPrice(text string, priceRegex *regexp.Regexp) float64 {
	// Try pattern-specific regex first
	if matches := priceRegex.FindStringSubmatch(text); len(matches) > 1 {
		if price, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return price
		}
	}

	// Try common price patterns
	for _, pattern := range PricePatterns {
		if matches := pattern.FindStringSubmatch(text); len(matches) > 1 {
			if price, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return price
			}
		}
	}

	return 0
}

// extractDate attempts to extract a date from text
func (d *Detector) extractDate(text string, dateRegex *regexp.Regexp, baseDate time.Time) *time.Time {
	// Try pattern-specific regex first
	if matches := dateRegex.FindStringSubmatch(text); len(matches) > 1 {
		dateStr := matches[1]

		// Try parsing as "in N days"
		if strings.Contains(toLower(text), "in") && strings.Contains(toLower(dateStr), "day") {
			// Extract number
			var days int
			for i := 0; i < len(dateStr); i++ {
				if dateStr[i] >= '0' && dateStr[i] <= '9' {
					numStr := ""
					for i < len(dateStr) && dateStr[i] >= '0' && dateStr[i] <= '9' {
						numStr += string(dateStr[i])
						i++
					}
					if n, err := strconv.Atoi(numStr); err == nil {
						days = n
						break
					}
				}
			}
			if days > 0 {
				date := baseDate.Add(time.Duration(days) * 24 * time.Hour)
				return &date
			}
		}

		// Try common date formats
		if date := parseDate(dateStr); date != nil {
			return date
		}
	}

	// Try common date patterns
	for _, pattern := range DatePatterns {
		if matches := pattern.FindStringSubmatch(text); len(matches) > 1 {
			if date := parseDate(matches[1]); date != nil {
				return date
			}
		}
	}

	return nil
}

// parseDate attempts to parse a date string in various formats
func parseDate(dateStr string) *time.Time {
	formats := []string{
		"2006-01-02",
		"1/2/2006",
		"01/02/2006",
		"January 2, 2006",
		"Jan 2, 2006",
		"January 2 2006",
		"Jan 2 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t
		}
	}

	return nil
}

// extractServiceFromSender extracts service name from sender email
func extractServiceFromSender(sender string) string {
	// Extract domain
	parts := strings.Split(sender, "@")
	if len(parts) != 2 {
		return ""
	}

	domain := parts[1]

	// Remove common suffixes
	domain = strings.TrimSuffix(domain, ".com")
	domain = strings.TrimSuffix(domain, ".io")
	domain = strings.TrimSuffix(domain, ".net")
	domain = strings.TrimSuffix(domain, ".org")

	// Split by dots and take the main part
	domainParts := strings.Split(domain, ".")
	if len(domainParts) > 0 {
		// Capitalize first letter
		serviceName := domainParts[0]
		if len(serviceName) > 0 {
			return strings.ToUpper(serviceName[0:1]) + serviceName[1:]
		}
	}

	return ""
}

// extractCancelURL attempts to find a cancellation URL in the text
func extractCancelURL(text string) string {
	if matches := CancelURLPattern.FindStringSubmatch(text); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Helper: Extract recipient email if it's in "sent to:" format
func extractRecipientEmail(text string) string {
	// Look for "sent to:", "delivered to:", etc.
	patterns := []string{
		`sent to:?\s*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
		`delivered to:?\s*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
		`for:?\s*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
	}

	textLower := toLower(text)
	for _, patternStr := range patterns {
		pattern := regexp.MustCompile(`(?i)` + patternStr)
		if matches := pattern.FindStringSubmatch(textLower); len(matches) > 1 {
			return matches[1]
		}
	}

	// Fallback: find any email in the text
	if matches := EmailPattern.FindStringSubmatch(text); len(matches) > 0 {
		return matches[0]
	}

	return ""
}
