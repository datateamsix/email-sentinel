/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/config"
	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/gmail"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show email-sentinel configuration status",
	Long: `Display the current status of email-sentinel configuration.

Shows:
- Authentication status
- Number of configured filters
- Configuration settings
- Notification settings

Example:
  email-sentinel status`,
	Run: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	fmt.Println("📊 Email Sentinel Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")

	// Check authentication
	if gmail.TokenExists() {
		fmt.Println("✅ Authentication: Configured")
		tokenPath, _ := config.TokenPath()
		fmt.Printf("   Token: %s\n", tokenPath)
	} else {
		fmt.Println("❌ Authentication: Not configured")
		fmt.Println("   Run: email-sentinel init")
	}
	fmt.Println("")

	// Check filters
	cfg, err := filter.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️  Configuration: Error loading (%v)\n", err)
		return
	}

	fmt.Printf("📋 Filters: %d configured\n", len(cfg.Filters))
	if len(cfg.Filters) > 0 {
		for i, f := range cfg.Filters {
			fmt.Printf("   [%d] %s\n", i+1, f.Name)
		}
	} else {
		fmt.Println("   Run: email-sentinel filter add")
	}
	fmt.Println("")

	// Show settings
	fmt.Printf("⚙️  Polling Interval: %d seconds\n", cfg.PollingInterval)
	fmt.Println("")

	// Notification settings
	fmt.Println("🔔 Notifications:")
	if cfg.Notifications.Desktop {
		fmt.Println("   Desktop: Enabled")
	} else {
		fmt.Println("   Desktop: Disabled")
	}

	if cfg.Notifications.Mobile.Enabled {
		fmt.Printf("   Mobile: Enabled (topic: %s)\n", cfg.Notifications.Mobile.NtfyTopic)
	} else {
		fmt.Println("   Mobile: Disabled")
	}
	fmt.Println("")

	// Config file location
	configPath, _ := config.ConfigPath()
	fmt.Printf("📁 Config File: %s\n", configPath)
}
