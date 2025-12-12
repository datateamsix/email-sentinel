/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/datateamsix/email-sentinel/internal/ui"
)

var (
	listTrialsOnly bool
	listPaidOnly   bool
	listFreeOnly   bool
)

// accountsListCmd represents the accounts list command
var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts or filter by type",
	Long: `List all detected digital accounts or filter by type.

Examples:
  email-sentinel accounts list              # Show all accounts
  email-sentinel accounts list --trials     # Show only trials
  email-sentinel accounts list --paid       # Show only paid subscriptions`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database
		db, err := storage.InitDB()
		if err != nil {
			fmt.Printf("%s Failed to initialize database: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}
		defer storage.CloseDB(db)

		var accounts []storage.Account

		// Determine which accounts to show
		if listTrialsOnly {
			accounts, err = storage.GetAccountsByType(db, "trial")
		} else if listPaidOnly {
			accounts, err = storage.GetAccountsByType(db, "paid")
		} else if listFreeOnly {
			accounts, err = storage.GetAccountsByType(db, "free")
		} else {
			accounts, err = storage.GetAllAccounts(db)
		}

		if err != nil {
			fmt.Printf("%s Failed to get accounts: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}

		if len(accounts) == 0 {
			fmt.Println(ui.ColorYellow.Sprint("No accounts found."))
			fmt.Println("\nEmail Sentinel will automatically detect accounts as you receive emails.")
			fmt.Println("Try running: email-sentinel start")
			return
		}

		// Get total monthly spend
		totalSpend, err := storage.GetTotalMonthlySpend(db)
		if err != nil {
			totalSpend = 0
		}

		// Display header
		title := "All Accounts"
		if listTrialsOnly {
			title = "Trial Accounts"
		} else if listPaidOnly {
			title = "Paid Subscriptions"
		} else if listFreeOnly {
			title = "Free Accounts"
		}

		fmt.Printf("\n%s\n", ui.ColorBold.Sprintf("📋 %s (%d total)", title, len(accounts)))
		fmt.Println(ui.ColorGray.Sprint("─────────────────────────────────────────────────────────────────"))

		// Convert to display type
		displayAccounts := make([]Account, len(accounts))
		for i, acc := range accounts {
			displayAccounts[i] = Account{
				ID:             acc.ID,
				ServiceName:    acc.ServiceName,
				EmailAddress:   acc.EmailAddress,
				AccountType:    acc.AccountType,
				Status:         acc.Status,
				PriceMonthly:   acc.PriceMonthly,
				TrialEndDate:   acc.TrialEndDate,
				DetectedAt:     acc.DetectedAt,
				Category:       acc.Category,
				CancelURL:      acc.CancelURL,
			}
		}

		// Display each account
		for i, acc := range displayAccounts {
			fmt.Println(formatAccount(storage.Account{
				ID:             acc.ID,
				ServiceName:    acc.ServiceName,
				EmailAddress:   acc.EmailAddress,
				AccountType:    acc.AccountType,
				Status:         acc.Status,
				PriceMonthly:   acc.PriceMonthly,
				TrialEndDate:   acc.TrialEndDate,
				DetectedAt:     acc.DetectedAt,
				Category:       acc.Category,
				CancelURL:      acc.CancelURL,
			}, i+1))
		}

		// Display summary
		fmt.Println(formatAccountSummary(displayAccounts, totalSpend))
	},
}

func init() {
	accountsCmd.AddCommand(accountsListCmd)

	accountsListCmd.Flags().BoolVar(&listTrialsOnly, "trials", false, "Show only trial accounts")
	accountsListCmd.Flags().BoolVar(&listPaidOnly, "paid", false, "Show only paid subscriptions")
	accountsListCmd.Flags().BoolVar(&listFreeOnly, "free", false, "Show only free accounts")
}
