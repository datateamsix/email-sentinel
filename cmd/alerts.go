/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
)

// alertsCmd represents the alerts command
var alertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "View today's email alerts",
	Long: `Display all email alerts that were triggered today.

This shows notifications that were sent during the current monitoring session,
allowing you to review missed alerts. The database is automatically wiped
at 12:00 AM daily.

Examples:
  # View today's alerts
  email-sentinel alerts

  # View last 5 alerts
  email-sentinel alerts --recent 5`,
	Run: runAlerts,
}

var recentLimit int

func init() {
	rootCmd.AddCommand(alertsCmd)
	alertsCmd.Flags().IntVarP(&recentLimit, "recent", "r", 0, "Show only N most recent alerts (0 = all today)")
}

func runAlerts(cmd *cobra.Command, args []string) {
	// Initialize database
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("❌ Error opening alert database: %v\n", err)
		os.Exit(1)
	}
	defer storage.CloseDB(db)

	var alerts []storage.Alert

	if recentLimit > 0 {
		// Get N most recent alerts
		alerts, err = storage.GetRecentAlerts(db, recentLimit)
		if err != nil {
			fmt.Printf("❌ Error fetching recent alerts: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Get all alerts from today
		alerts, err = storage.GetTodayAlerts(db)
		if err != nil {
			fmt.Printf("❌ Error fetching today's alerts: %v\n", err)
			os.Exit(1)
		}
	}

	if len(alerts) == 0 {
		if recentLimit > 0 {
			fmt.Println("📭 No alerts found")
		} else {
			fmt.Println("📭 No alerts today")
		}
		return
	}

	// Display header
	if recentLimit > 0 {
		fmt.Printf("📬 Last %d Alert(s)\n\n", len(alerts))
	} else {
		count, _ := storage.CountTodayAlerts(db)
		fmt.Printf("📬 Today's Alerts (%d total)\n\n", count)
	}

	// Display each alert
	for i, alert := range alerts {
		// Add priority indicator
		priorityIcon := "📩" // Normal priority
		if alert.Priority == 1 {
			priorityIcon = "🔥" // High priority
		}

		fmt.Printf("[%d] %s %s\n", i+1, priorityIcon, alert.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Filter: %s\n", alert.FilterName)
		if alert.Priority == 1 {
			fmt.Printf("    Priority: HIGH\n")
		}
		fmt.Printf("    From:   %s\n", alert.Sender)
		fmt.Printf("    Subject: %s\n", alert.Subject)

		if alert.Snippet != "" {
			// Truncate snippet if too long
			snippet := alert.Snippet
			if len(snippet) > 100 {
				snippet = snippet[:97] + "..."
			}
			fmt.Printf("    Preview: %s\n", snippet)
		}

		fmt.Printf("    Link:   %s\n", alert.GmailLink)
		fmt.Println()
	}
}
