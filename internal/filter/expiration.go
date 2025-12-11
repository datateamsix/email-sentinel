package filter

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseExpiration parses expiration string and returns a time.Time pointer
// Supports: "1d", "7d", "30d", "60d", "90d", "never", "YYYY-MM-DD"
// Returns nil for "never" or invalid input
func ParseExpiration(input string) (*time.Time, error) {
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" || input == "never" || input == "0" {
		return nil, nil
	}

	// Try parsing as duration (e.g., "7d", "30d")
	if strings.HasSuffix(input, "d") {
		daysStr := strings.TrimSuffix(input, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return nil, fmt.Errorf("invalid duration format: %s (expected format: 7d, 30d, etc.)", input)
		}
		if days <= 0 {
			return nil, fmt.Errorf("expiration days must be positive: %d", days)
		}
		expiresAt := time.Now().Add(time.Duration(days) * 24 * time.Hour)
		return &expiresAt, nil
	}

	// Try parsing as date (YYYY-MM-DD)
	expiresAt, err := time.Parse("2006-01-02", input)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s (expected: YYYY-MM-DD or duration like 7d)", input)
	}

	// Ensure date is in the future
	if expiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("expiration date must be in the future: %s", input)
	}

	return &expiresAt, nil
}

// FormatExpiration returns a human-readable expiration status
func FormatExpiration(expiresAt *time.Time) string {
	if expiresAt == nil {
		return "never expires"
	}

	now := time.Now()
	if expiresAt.Before(now) {
		// Expired
		duration := now.Sub(*expiresAt)
		if duration < 24*time.Hour {
			return fmt.Sprintf("expired %d hours ago", int(duration.Hours()))
		}
		return fmt.Sprintf("expired %d days ago", int(duration.Hours()/24))
	}

	// Not expired yet
	duration := expiresAt.Sub(now)
	if duration < 24*time.Hour {
		return fmt.Sprintf("expires in %d hours", int(duration.Hours()))
	}
	days := int(duration.Hours() / 24)
	if days == 1 {
		return "expires tomorrow"
	}
	return fmt.Sprintf("expires in %d days", days)
}

// IsExpired checks if a filter has expired (accounting for 24-hour grace period)
func IsExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	// Add 24-hour grace period
	graceDeadline := expiresAt.Add(24 * time.Hour)
	return time.Now().After(graceDeadline)
}

// IsInGracePeriod checks if a filter is expired but within grace period
func IsInGracePeriod(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	now := time.Now()
	return now.After(*expiresAt) && now.Before(expiresAt.Add(24*time.Hour))
}

// CleanupExpiredFilters removes filters that have expired beyond the grace period
// Returns a list of filter names that were removed
func CleanupExpiredFilters() ([]string, error) {
	filters, err := ListFilters()
	if err != nil {
		return nil, err
	}

	var removed []string
	var remaining []Filter

	for _, f := range filters {
		if IsExpired(f.ExpiresAt) {
			// Filter has expired beyond grace period, remove it
			removed = append(removed, f.Name)
		} else {
			// Keep this filter
			remaining = append(remaining, f)
		}
	}

	// If any filters were removed, save the updated config
	if len(removed) > 0 {
		cfg, err := LoadConfig()
		if err != nil {
			return nil, err
		}
		cfg.Filters = remaining
		if err := SaveConfig(cfg); err != nil {
			return nil, err
		}
	}

	return removed, nil
}
