/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/filter"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured email filters",
	Long: `Display all configured email filters with their settings.

Shows filter names, sender patterns, subject patterns, and match modes.

Example:
  email-sentinel filter list`,
	Run: runFilterList,
}

func init() {
	filterCmd.AddCommand(listCmd)
}

func runFilterList(cmd *cobra.Command, args []string) {
	filters, err := filter.ListFilters()
	if err != nil {
		fmt.Printf("❌ Error loading filters: %v\n", err)
		os.Exit(1)
	}

	if len(filters) == 0 {
		fmt.Println("No filters configured.")
		fmt.Println("\nAdd one with: email-sentinel filter add")
		return
	}

	fmt.Printf("\n📋 Email Filters (%d)\n", len(filters))
	fmt.Println(strings.Repeat("━", 60))

	for i, f := range filters {
		fmt.Printf("\n[%d] %s\n", i+1, f.Name)

		if len(f.From) > 0 {
			fmt.Printf("    From:    %s\n", strings.Join(f.From, ", "))
		} else {
			fmt.Println("    From:    (any)")
		}

		if len(f.Subject) > 0 {
			fmt.Printf("    Subject: %s\n", strings.Join(f.Subject, ", "))
		} else {
			fmt.Println("    Subject: (any)")
		}

		if len(f.Labels) > 0 {
			fmt.Printf("    Labels:  🏷️  %s\n", strings.Join(f.Labels, ", "))
		}

		matchDesc := "any (OR - either condition triggers)"
		if f.Match == "all" {
			matchDesc = "all (AND - all conditions must match)"
		}
		fmt.Printf("    Match:   %s\n", matchDesc)

		// Show Gmail scope
		scope := f.GmailScope
		if scope == "" {
			scope = "inbox"
		}
		fmt.Printf("    Scope:   📬 %s\n", scope)

		// Show expiration status
		expirationStatus := filter.FormatExpiration(f.ExpiresAt)
		if filter.IsInGracePeriod(f.ExpiresAt) {
			fmt.Printf("    Expires: ⚠️  %s (grace period - will be deleted in <24hrs)\n", expirationStatus)
		} else if filter.IsExpired(f.ExpiresAt) {
			fmt.Printf("    Expires: ❌ %s (will be deleted)\n", expirationStatus)
		} else if f.ExpiresAt != nil {
			fmt.Printf("    Expires: ⏰ %s\n", expirationStatus)
		} else {
			fmt.Printf("    Expires: ♾️  %s\n", expirationStatus)
		}
	}

	fmt.Println("")
}
