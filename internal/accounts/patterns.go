/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package accounts

import (
	"regexp"
)

// GetDefaultPatterns returns the default set of account detection patterns
// These are generic patterns that work across most services
func GetDefaultPatterns() []DetectionPattern {
	return []DetectionPattern{
		// Trial start patterns
		{
			Name:         "trial_start_generic",
			Type:         "trial",
			Keywords:     []string{"free trial", "trial period", "trial started", "trial membership", "start your trial", "trial begins"},
			ServiceRegex: regexp.MustCompile(`(?i)(?:welcome to|thanks for joining|you.?re now a member of|trial for)\s+([A-Z][A-Za-z0-9\s]+?)(?:\s+(?:premium|plus|pro|free trial))?(?:\.|!|,)`),
			PriceRegex:   regexp.MustCompile(`(?i)\$(\d+(?:\.\d{2})?)\s*(?:per|/)\s*(?:month|mo)`),
			DateRegex:    regexp.MustCompile(`(?i)(?:trial\s+)?(?:ends?|expires?)\s+(?:on\s+)?(\d{1,2}[-/]\d{1,2}[-/]\d{2,4}|\w+\s+\d{1,2},?\s+\d{4})`),
			Confidence:   0.85,
		},
		{
			Name:         "trial_ending_soon",
			Type:         "trial",
			Keywords:     []string{"trial expires", "trial ends", "trial ending soon", "last chance", "trial will expire"},
			ServiceRegex: regexp.MustCompile(`(?i)(?:your|the)\s+([A-Z][A-Za-z0-9\s]+?)\s+(?:trial|free trial|membership)`),
			PriceRegex:   regexp.MustCompile(`(?i)\$(\d+(?:\.\d{2})?)\s*(?:per|/)\s*(?:month|mo)`),
			DateRegex:    regexp.MustCompile(`(?i)(?:on|in)\s+(\d{1,2})\s+(?:day|hour)`),
			Confidence:   0.90,
		},

		// Subscription/payment patterns
		{
			Name:         "subscription_payment",
			Type:         "paid",
			Keywords:     []string{"subscription renewed", "payment successful", "subscription confirmed", "payment processed", "billing successful"},
			ServiceRegex: regexp.MustCompile(`(?i)(?:for|your)\s+([A-Z][A-Za-z0-9\s]+?)\s+(?:subscription|membership|plan)`),
			PriceRegex:   regexp.MustCompile(`(?i)(?:total|amount|charged|paid):\s*\$(\d+(?:\.\d{2})?)`),
			Confidence:   0.90,
		},
		{
			Name:         "recurring_payment",
			Type:         "paid",
			Keywords:     []string{"monthly subscription", "annual subscription", "recurring payment", "auto-renew", "automatic renewal"},
			ServiceRegex: regexp.MustCompile(`(?i)([A-Z][A-Za-z0-9\s]+?)\s+(?:subscription|membership|plan)`),
			PriceRegex:   regexp.MustCompile(`(?i)\$(\d+(?:\.\d{2})?)`),
			Confidence:   0.85,
		},

		// Account creation patterns
		{
			Name:         "account_created",
			Type:         "free",
			Keywords:     []string{"welcome to", "account created", "verify your email", "confirm your account", "registration successful", "account activated"},
			ServiceRegex: regexp.MustCompile(`(?i)(?:welcome to|you.?ve joined|thanks for signing up for)\s+([A-Z][A-Za-z0-9\s]+?)(?:\.|!|,)`),
			Confidence:   0.75,
		},

		// Cancellation patterns
		{
			Name:         "subscription_cancelled",
			Type:         "cancellation",
			Keywords:     []string{"subscription cancelled", "subscription canceled", "membership ended", "auto-renew disabled", "will not be charged"},
			ServiceRegex: regexp.MustCompile(`(?i)(?:your|the)\s+([A-Z][A-Za-z0-9\s]+?)\s+(?:subscription|membership)`),
			Confidence:   0.85,
		},
	}
}

// Common price extraction patterns
var (
	// PricePatterns contains various price extraction patterns
	PricePatterns = []*regexp.Regexp{
		regexp.MustCompile(`\$(\d+(?:\.\d{2})?)`),                                    // $9.99
		regexp.MustCompile(`USD\s*(\d+(?:\.\d{2})?)`),                                 // USD 9.99
		regexp.MustCompile(`(?i)(?:price|cost|total|amount):\s*\$?(\d+(?:\.\d{2})?)`), // Price: $9.99
		regexp.MustCompile(`(?i)(\d+(?:\.\d{2})?)\s*(?:per|/)\s*(?:month|mo)`),       // 9.99 per month
	}

	// DatePatterns contains various date extraction patterns
	DatePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`),                           // 2025-01-15
		regexp.MustCompile(`(\d{1,2}/\d{1,2}/\d{4})`),                       // 1/15/2025
		regexp.MustCompile(`(\w+\s+\d{1,2},?\s+\d{4})`),                     // January 15, 2025
		regexp.MustCompile(`(?i)in\s+(\d+)\s+day`),                          // in 7 days
		regexp.MustCompile(`(?i)(?:on|expires)\s+(\w+\s+\d{1,2})`),         // on January 15
	}

	// EmailPattern extracts email addresses
	EmailPattern = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)

	// CancelURLPattern extracts cancellation URLs
	CancelURLPattern = regexp.MustCompile(`(?i)(?:cancel|unsubscribe|manage.*subscription).*?(https?://[^\s<>"]+)`)
)

// ServiceCategoryKeywords maps keywords to service categories
var ServiceCategoryKeywords = map[string][]string{
	"streaming": {
		"netflix", "hulu", "disney", "spotify", "apple music", "youtube premium",
		"prime video", "hbo", "paramount", "peacock", "tidal", "pandora",
	},
	"software": {
		"adobe", "microsoft", "github", "notion", "grammarly", "canva",
		"figma", "sketch", "photoshop", "office", "creative cloud",
	},
	"cloud": {
		"aws", "amazon web services", "google cloud", "dropbox", "icloud",
		"onedrive", "azure", "cloudflare", "backblaze",
	},
	"productivity": {
		"chatgpt", "slack", "zoom", "asana", "trello", "monday",
		"clickup", "notion", "evernote", "todoist",
	},
}

// DetermineCategory attempts to determine the service category based on service name
func DetermineCategory(serviceName string) string {
	serviceNameLower := toLower(serviceName)

	for category, keywords := range ServiceCategoryKeywords {
		for _, keyword := range keywords {
			if contains(serviceNameLower, keyword) {
				return category
			}
		}
	}

	return "other"
}

// Helper functions
func toLower(s string) string {
	result := ""
	for _, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			result += string(ch + 32)
		} else {
			result += string(ch)
		}
	}
	return result
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
