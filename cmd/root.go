/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/ui"
)

var versionFlag bool

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

Interactive Mode:
  Run 'email-sentinel' without arguments to open the interactive menu.

More Info: https://github.com/datateamsix/email-sentinel`,
	Run: func(cmd *cobra.Command, args []string) {
		// Handle version flag
		if versionFlag {
			fmt.Printf("Email Sentinel v%s\n", ui.AppVersion)
			return
		}

		// No subcommand provided - launch interactive mode
		runInteractive()
	},
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
	// Add version flag
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.email-sentinel.yaml)")
}

// runInteractive launches the interactive menu system
func runInteractive() {
	// Clear screen and show banner
	ui.ClearScreen()
	ui.PrintBanner(ui.AppVersion)

	// Check if first-time setup needed
	if ui.ShouldRunWizard() {
		fmt.Println()
		ui.PrintInfo("Welcome! It looks like this is your first time running Email Sentinel.")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(ui.ColorGreen.Sprint("Would you like to run the setup wizard? [Y/n]: "))

		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "" || response == "y" || response == "yes" {
			wizard := ui.NewWizard()
			if err := wizard.Run(); err != nil {
				if err.Error() != "wizard cancelled by user" {
					ui.PrintError(fmt.Sprintf("Setup wizard error: %v", err))
					fmt.Println()
					ui.PrintInfo("You can run the wizard again later from the main menu.")
					fmt.Println()
					fmt.Print("Press Enter to continue...")
					reader.ReadString('\n')
				} else {
					// User quit wizard
					ui.PrintInfo("Setup wizard cancelled. You can run it again from the main menu.")
					fmt.Println()
					fmt.Print("Press Enter to continue...")
					reader.ReadString('\n')
				}
			}
		} else {
			ui.PrintInfo("Skipping setup wizard. You can run it later from the main menu.")
			fmt.Println()
			fmt.Print("Press Enter to continue...")
			reader.ReadString('\n')
		}
	}

	// Launch main menu
	ui.RunInteractiveMenu()
}
