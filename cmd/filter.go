/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Manage email filters",
	Long: `Manage email filters to match incoming emails.

Filters can match by sender address, subject keywords, or both.
Use subcommands to add, list, edit, or remove filters.

Available Commands:
  add     Add a new filter
  list    List all filters
  edit    Edit an existing filter
  remove  Remove a filter

Examples:
  email-sentinel filter add --name "Jobs" --from "linkedin.com"
  email-sentinel filter list
  email-sentinel filter edit "Jobs"
  email-sentinel filter remove "Jobs"`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)
}
