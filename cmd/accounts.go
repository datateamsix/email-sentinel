/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/datateamsix/email-sentinel/internal/ui"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage digital accounts (subscriptions, trials, free accounts)",
	Long: `Manage your digital accounts including subscriptions, trials, and free accounts.

Email Sentinel automatically detects account-related emails and tracks:
- Which email address you used for each service
- Trial expiration dates
- Monthly subscription costs
- Total spending across all services

Available Commands:
  list     List all accounts or filter by type
  search   Search for a specific service
  remove   Remove an account by ID
  refresh  Re-scan Gmail to detect accounts

Examples:
  email-sentinel accounts list
  email-sentinel accounts list --trials
  email-sentinel accounts list --paid
  email-sentinel accounts search netflix`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
}

// formatAccount formats an account for display
func formatAccount(acc storage.Account, index int) string {
	var sb strings.Builder

	// Status icon
	statusIcon := "✅"
	if acc.Status == "cancelled" {
		statusIcon = "❌"
	}

	// Type icon
	typeIcon := ""
	switch acc.AccountType {
	case "trial":
		typeIcon = "🆓"
	case "paid":
		typeIcon = "💳"
	case "free":
		typeIcon = "🎁"
	}

	// Header with service name
	sb.WriteString(fmt.Sprintf("[%d] %s %s %s",
		index,
		statusIcon,
		typeIcon,
		ui.ColorBold.Sprint(acc.ServiceName),
	))

	// Show trial expiration if applicable
	if acc.TrialEndDate != nil {
		daysUntil := time.Until(*acc.TrialEndDate).Hours() / 24
		if daysUntil > 0 {
			if daysUntil <= 1 {
				sb.WriteString(fmt.Sprintf("  %s", ui.ColorRed.Sprintf("🔥 Expires in %d day(s)", int(daysUntil)+1)))
			} else if daysUntil <= 3 {
				sb.WriteString(fmt.Sprintf("  %s", ui.ColorYellow.Sprintf("⚠️  Expires in %d days", int(daysUntil)+1)))
			} else {
				sb.WriteString(fmt.Sprintf("  (Expires in %d days)", int(daysUntil)+1))
			}
		} else {
			sb.WriteString(fmt.Sprintf("  %s", ui.ColorRed.Sprint("❌ Expired")))
		}
	}

	// Show price if available
	if acc.PriceMonthly > 0 {
		sb.WriteString(fmt.Sprintf("  $%.2f/mo", acc.PriceMonthly))
	}

	sb.WriteString("\n")

	// Details
	sb.WriteString(fmt.Sprintf("    Email: %s\n", ui.ColorCyan.Sprint(acc.EmailAddress)))
	sb.WriteString(fmt.Sprintf("    Type: %s", acc.AccountType))

	if acc.Category != "" && acc.Category != "other" {
		sb.WriteString(fmt.Sprintf(" | Category: %s", acc.Category))
	}

	sb.WriteString("\n")

	if acc.CancelURL != "" {
		sb.WriteString(fmt.Sprintf("    Cancel: %s\n", ui.ColorGray.Sprint(acc.CancelURL)))
	}

	sb.WriteString(fmt.Sprintf("    Detected: %s\n", formatTimestamp(acc.DetectedAt)))

	return sb.String()
}

// formatAccountSummary formats a summary of accounts
func formatAccountSummary(accounts []Account, totalSpend float64) string {
	var sb strings.Builder

	trialCount := 0
	paidCount := 0
	freeCount := 0
	expiringCount := 0

	emailsUsed := make(map[string]bool)

	for _, acc := range accounts {
		if acc.Status != "active" {
			continue
		}

		emailsUsed[acc.EmailAddress] = true

		switch acc.AccountType {
		case "trial":
			trialCount++
			if acc.TrialEndDate != nil {
				daysUntil := time.Until(*acc.TrialEndDate).Hours() / 24
				if daysUntil > 0 && daysUntil <= 7 {
					expiringCount++
				}
			}
		case "paid":
			paidCount++
		case "free":
			freeCount++
		}
	}

	sb.WriteString(ui.ColorBold.Sprintf("\n📊 Account Summary\n"))
	sb.WriteString(fmt.Sprintf("   Total accounts: %d (%d active)\n", len(accounts), trialCount+paidCount+freeCount))
	sb.WriteString(fmt.Sprintf("   Trials: %d", trialCount))
	if expiringCount > 0 {
		sb.WriteString(ui.ColorYellow.Sprintf(" (%d expiring soon)", expiringCount))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("   Paid: %d\n", paidCount))
	sb.WriteString(fmt.Sprintf("   Free: %d\n", freeCount))
	sb.WriteString(fmt.Sprintf("   Emails used: %d\n", len(emailsUsed)))

	if totalSpend > 0 {
		sb.WriteString(fmt.Sprintf("\n💰 Total: $%.2f/month ($%.2f/year)\n", totalSpend, totalSpend*12))
	}

	return sb.String()
}

// Account type for display (mirrors storage.Account)
type Account struct {
	ID             int64
	ServiceName    string
	EmailAddress   string
	AccountType    string
	Status         string
	PriceMonthly   float64
	TrialEndDate   *time.Time
	DetectedAt     time.Time
	Category       string
	CancelURL      string
}
