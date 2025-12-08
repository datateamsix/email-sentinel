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

	"github.com/datateamsix/email-sentinel/internal/storage"
)

// alertsClearCmd represents the alerts clear command
var alertsClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all email alerts",
	Long: `Delete all email alerts from the database.

This helps keep the database clean by removing all stored alerts.
You'll be prompted for confirmation unless --force is used.

Examples:
  # Clear all alerts with confirmation
  email-sentinel alerts clear

  # Clear all alerts without confirmation
  email-sentinel alerts clear --force`,
	Run: runAlertsClear,
}

var forceAlertsClear bool

func init() {
	alertsCmd.AddCommand(alertsClearCmd)
	alertsClearCmd.Flags().BoolVarP(&forceAlertsClear, "force", "f", false, "Skip confirmation prompt")
}

func runAlertsClear(cmd *cobra.Command, args []string) {
	// Open database
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("❌ Error opening database: %v\n", err)
		fmt.Println("   Tip: Database may not exist. Start monitoring with 'email-sentinel start' first.")
		os.Exit(1)
	}
	defer storage.CloseDB(db)

	// Count current alerts
	count, err := storage.CountTodayAlerts(db)
	if err != nil {
		// Try to count all alerts instead
		alerts, err := storage.GetRecentAlerts(db, 10000)
		if err != nil {
			fmt.Printf("❌ Error counting alerts: %v\n", err)
			os.Exit(1)
		}
		count = len(alerts)
	}

	if count == 0 {
		fmt.Println("✨ No alerts to clear")
		return
	}

	// Prompt for confirmation unless --force is used
	if !forceAlertsClear {
		fmt.Printf("Found %d alert(s)\n", count)
		fmt.Print("Delete all alerts? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ Error reading input: %v\n", err)
			os.Exit(1)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return
		}
	}

	// Delete all alerts
	deleted, err := storage.DeleteAllAlerts(db)
	if err != nil {
		fmt.Printf("❌ Error deleting alerts: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🗑️  Cleared %d alert(s)\n", deleted)
}
