package filter

import (
	"fmt"
	"strings"

	"github.com/datateamsix/email-sentinel/internal/config"
)

// LoadConfig loads the config or returns default
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	if !config.ConfigExists() {
		return cfg, nil
	}

	if err := config.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// SaveConfig saves the config to disk
func SaveConfig(cfg *Config) error {
	return config.Save(cfg)
}

// AddFilter adds a new filter to the config
func AddFilter(f Filter) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	// Check for duplicate name
	for _, existing := range cfg.Filters {
		if strings.EqualFold(existing.Name, f.Name) {
			return fmt.Errorf("filter '%s' already exists", f.Name)
		}
	}

	cfg.Filters = append(cfg.Filters, f)
	return SaveConfig(cfg)
}

// UpdateFilter updates a filter at a specific index
func UpdateFilter(index int, updated Filter) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	if index < 0 || index >= len(cfg.Filters) {
		return fmt.Errorf("filter index out of range")
	}

	// Check for duplicate name (excluding current filter)
	for i, existing := range cfg.Filters {
		if i != index && strings.EqualFold(existing.Name, updated.Name) {
			return fmt.Errorf("filter '%s' already exists", updated.Name)
		}
	}

	cfg.Filters[index] = updated
	return SaveConfig(cfg)
}

// RemoveFilter removes a filter by name
func RemoveFilter(name string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	found := false
	newFilters := []Filter{}
	for _, f := range cfg.Filters {
		if strings.EqualFold(f.Name, name) {
			found = true
			continue
		}
		newFilters = append(newFilters, f)
	}

	if !found {
		return fmt.Errorf("filter '%s' not found", name)
	}

	cfg.Filters = newFilters
	return SaveConfig(cfg)
}

// ListFilters returns all filters
func ListFilters() ([]Filter, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Filters, nil
}

// MatchesFilter checks if an email matches a given filter
func MatchesFilter(f Filter, fromAddress string, subject string) bool {
	fromAddress = strings.ToLower(fromAddress)
	subject = strings.ToLower(subject)

	fromMatched := false
	subjectMatched := false

	// Check from patterns
	if len(f.From) == 0 {
		fromMatched = true // No from filter means auto-match
	} else {
		for _, pattern := range f.From {
			if strings.Contains(fromAddress, strings.ToLower(pattern)) {
				fromMatched = true
				break
			}
		}
	}

	// Check subject patterns
	if len(f.Subject) == 0 {
		subjectMatched = true // No subject filter means auto-match
	} else {
		for _, pattern := range f.Subject {
			if strings.Contains(subject, strings.ToLower(pattern)) {
				subjectMatched = true
				break
			}
		}
	}

	// Apply match mode
	if f.Match == "all" {
		// AND logic - both must match (considering empty patterns)
		if len(f.From) > 0 && len(f.Subject) > 0 {
			return fromMatched && subjectMatched
		}
		// If only one type of filter exists, just check that one
		if len(f.From) > 0 {
			return fromMatched
		}
		if len(f.Subject) > 0 {
			return subjectMatched
		}
		return false
	}

	// "any" (OR) logic - either can match
	if len(f.From) > 0 && fromMatched {
		return true
	}
	if len(f.Subject) > 0 && subjectMatched {
		return true
	}

	return false
}

// CheckAllFilters checks an email against all filters and returns matching filter names
func CheckAllFilters(fromAddress string, subject string) ([]string, error) {
	filters, err := ListFilters()
	if err != nil {
		return nil, err
	}

	var matchedFilters []string
	for _, f := range filters {
		if MatchesFilter(f, fromAddress, subject) {
			matchedFilters = append(matchedFilters, f.Name)
		}
	}

	return matchedFilters, nil
}
