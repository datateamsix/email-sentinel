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

// accountsSearchCmd represents the accounts search command
var accountsSearchCmd = &cobra.Command{
	Use:   "search <service>",
	Short: "Search for a specific service",
	Long: `Search for accounts by service name (case-insensitive).

This is the killer feature - quickly find which email you used for a service!

Examples:
  email-sentinel accounts search netflix
  email-sentinel accounts search adobe
  email-sentinel accounts search spotify`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := args[0]

		// Initialize database
		db, err := storage.InitDB()
		if err != nil {
			fmt.Printf("%s Failed to initialize database: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}
		defer storage.CloseDB(db)

		// Search for accounts
		accounts, err := storage.SearchAccounts(db, searchTerm)
		if err != nil {
			fmt.Printf("%s Failed to search accounts: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}

		if len(accounts) == 0 {
			fmt.Printf("%s No accounts found matching '%s'\n", ui.ColorYellow.Sprint("ℹ"), searchTerm)
			fmt.Println("\nTip: Make sure Email Sentinel is running to detect accounts from your emails.")
			return
		}

		// Display header
		fmt.Printf("\n%s\n", ui.ColorBold.Sprintf("🔍 Search Results for '%s' (%d found)", searchTerm, len(accounts)))
		fmt.Println(ui.ColorGray.Sprint("─────────────────────────────────────────────────────────────────"))

		// Display each account
		for i, acc := range accounts {
			fmt.Println(formatAccount(acc, i+1))
		}

		// If only one result, highlight the email
		if len(accounts) == 1 {
			fmt.Printf("\n%s %s is using: %s\n",
				ui.ColorGreen.Sprint("✓"),
				ui.ColorBold.Sprint(accounts[0].ServiceName),
				ui.ColorCyan.Sprint(accounts[0].EmailAddress),
			)
		}
	},
}

func init() {
	accountsCmd.AddCommand(accountsSearchCmd)
}
