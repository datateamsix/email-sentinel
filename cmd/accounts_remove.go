/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/datateamsix/email-sentinel/internal/ui"
)

// accountsRemoveCmd represents the accounts remove command
var accountsRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove an account by ID",
	Long: `Remove an account from the database by its ID.

The ID is shown in brackets when you list accounts.

Example:
  email-sentinel accounts list
  email-sentinel accounts remove 3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Parse ID
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("%s Invalid account ID: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}

		// Initialize database
		db, err := storage.InitDB()
		if err != nil {
			fmt.Printf("%s Failed to initialize database: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}
		defer storage.CloseDB(db)

		// Delete the account
		err = storage.DeleteAccount(db, id)
		if err != nil {
			fmt.Printf("%s Failed to remove account: %v\n", ui.ColorRed.Sprint("✗"), err)
			return
		}

		fmt.Printf("%s Account #%d removed successfully\n", ui.ColorGreen.Sprint("✓"), id)
	},
}

func init() {
	accountsCmd.AddCommand(accountsRemoveCmd)
}
