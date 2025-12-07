/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "email-sentinel",
	Short: "Monitor Gmail and get instant notifications for filtered emails",
	Long: `Email Sentinel - Real-time Gmail monitoring with custom filters

Monitor your Gmail inbox and receive instant desktop and mobile notifications
when emails match your custom filters. Perfect for tracking important emails
without constantly checking your inbox.

Features:
  • Gmail API integration with OAuth 2.0 authentication
  • Flexible filters by sender, subject, or both (AND/OR logic)
  • Desktop notifications (Windows, macOS, Linux)
  • Mobile push notifications via ntfy.sh (free, no account needed)
  • Low resource usage with configurable polling intervals

New in this version:
  • Windows toast notifications - Rich, clickable notifications in Action Center
  • Smart priority rules - YAML-based rules engine for urgent email classification
  • Alert history - SQLite database stores all alerts with daily auto-cleanup
  • System tray app - Background mode with tray icon showing recent alerts
  • Direct Gmail links - Click any alert to open email in browser

Quick Start:
  1. email-sentinel init                    # Authenticate with Gmail
  2. email-sentinel filter add              # Create your first filter
  3. email-sentinel test desktop            # Test notifications
  4. email-sentinel start --tray            # Start monitoring with tray icon
  5. email-sentinel alerts                  # View alert history

More Info: https://github.com/yourusername/email-sentinel`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.email-sentinel.yaml)")
}


