/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// dbCmd represents the database management command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long: `Database management commands for Email Sentinel.

Subcommands:
  backup     Create a database backup

Examples:
  # Create a manual backup
  email-sentinel db backup`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
