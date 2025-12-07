/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"github.com/datateamsix/email-sentinel/internal/ui"
	"github.com/spf13/cobra"
)

// menuCmd represents the menu command
var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Launch interactive menu interface",
	Long: `Launch the interactive menu interface for Email Sentinel.

The menu provides a user-friendly way to:
  - Start/stop monitoring
  - Manage email filters
  - Configure notifications
  - View status and alert history
  - Adjust settings

Example:
  email-sentinel menu`,
	Run: runMenu,
}

func init() {
	rootCmd.AddCommand(menuCmd)
}

func runMenu(cmd *cobra.Command, args []string) {
	ui.RunInteractiveMenu()
}
