/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/datateamsix/email-sentinel/internal/appconfig"
	"github.com/spf13/cobra"
)

// configMigrateCmd represents the config migrate command
var configMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate from old config files to unified app-config.yaml",
	Long: `Migrates configuration from the old separate files (ai-config.yaml,
rules.yaml, otp_rules.yaml) to the new unified app-config.yaml format.

This command will:
1. Look for old config files in the config directory
2. Merge them into a new unified app-config.yaml
3. Preserve all your existing settings
4. Keep the old files as backup (does not delete them)

Example:
  email-sentinel config migrate`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Starting configuration migration...")

		// Try to load config (which will trigger migration if needed)
		cfg, err := appconfig.Load()
		if err != nil {
			fmt.Printf("❌ Migration failed: %v\n", err)
			os.Exit(1)
		}

		// Show what was loaded
		fmt.Println("\n✅ Configuration loaded successfully!")
		fmt.Println("\n📋 Configuration Summary:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Monitoring
		fmt.Printf("\n📊 Monitoring:\n")
		fmt.Printf("   Polling Interval: %d seconds\n", cfg.Monitoring.PollingInterval)
		fmt.Printf("   WAL Mode: %v\n", cfg.Monitoring.Database.WALMode)
		fmt.Printf("   Cleanup Interval: %s\n", cfg.Monitoring.Database.CleanupInterval)

		// AI Summary
		fmt.Printf("\n🤖 AI Summary:\n")
		fmt.Printf("   Enabled: %v\n", cfg.AISummary.Enabled)
		fmt.Printf("   Provider: %s\n", cfg.AISummary.Provider)
		fmt.Printf("   Cache Enabled: %v\n", cfg.AISummary.Cache.Enabled)

		// Priority Rules
		fmt.Printf("\n⚡ Priority Rules:\n")
		fmt.Printf("   Urgent Keywords: %d configured\n", len(cfg.Priority.UrgentKeywords))
		fmt.Printf("   VIP Senders: %d configured\n", len(cfg.Priority.VIPSenders))
		fmt.Printf("   VIP Domains: %d configured\n", len(cfg.Priority.VIPDomains))

		// OTP
		fmt.Printf("\n🔐 OTP Detection:\n")
		fmt.Printf("   Enabled: %v\n", cfg.OTP.Enabled)
		fmt.Printf("   Expiry Duration: %s\n", cfg.OTP.ExpiryDuration)
		fmt.Printf("   Trusted Senders: %d configured\n", len(cfg.OTP.TrustedSenders))
		fmt.Printf("   Auto-copy to Clipboard: %v\n", cfg.OTP.Clipboard.AutoCopy)

		// Notifications
		fmt.Printf("\n🔔 Notifications:\n")
		fmt.Printf("   Desktop: %v\n", cfg.Notifications.Desktop.Enabled)
		fmt.Printf("   Mobile: %v\n", cfg.Notifications.Mobile.Enabled)
		fmt.Printf("   Weekend Mode: %s\n", cfg.Notifications.WeekendMode)
		if cfg.Notifications.QuietHours.Start != "" {
			fmt.Printf("   Quiet Hours: %s - %s\n", cfg.Notifications.QuietHours.Start, cfg.Notifications.QuietHours.End)
		}

		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		// Show config path
		configPath, _ := appconfig.ConfigPath()
		fmt.Printf("\n📁 Config file: %s\n", configPath)
		fmt.Println("\n💡 Tip: You can now edit app-config.yaml to customize your settings")
		fmt.Println("   The old config files are kept as backup and not deleted")
	},
}

func init() {
	configCmd.AddCommand(configMigrateCmd)
}
