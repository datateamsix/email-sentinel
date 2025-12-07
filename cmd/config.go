/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/config"
	"github.com/datateamsix/email-sentinel/internal/filter"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and modify configuration",
	Long: `View and modify email-sentinel configuration settings.

Subcommands:
  show      Display current configuration
  set       Modify configuration values

Examples:
  # Show current config
  email-sentinel config show

  # Set polling interval
  email-sentinel config set polling 30

  # Enable mobile notifications
  email-sentinel config set mobile true

  # Set ntfy topic
  email-sentinel config set ntfy_topic "my-topic"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default to show if no subcommand
		if len(args) == 0 {
			runConfigShow(cmd, args)
		} else {
			cmd.Help()
		}
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Run:   runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  polling          Polling interval in seconds (default: 45)
  desktop          Enable/disable desktop notifications (true/false)
  mobile           Enable/disable mobile notifications (true/false)
  ntfy_topic       Set ntfy.sh topic for mobile notifications

Examples:
  email-sentinel config set polling 60
  email-sentinel config set desktop false
  email-sentinel config set mobile true
  email-sentinel config set ntfy_topic "my-secret-topic"`,
	Args: cobra.ExactArgs(2),
	Run:  runConfigSet,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) {
	cfg, err := filter.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	configPath, _ := config.ConfigPath()

	fmt.Println("\n⚙️  Email Sentinel Configuration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nConfig File: %s\n", configPath)
	fmt.Println("")
	fmt.Printf("Polling Interval:     %d seconds\n", cfg.PollingInterval)
	fmt.Printf("Desktop Notifications: %v\n", cfg.Notifications.Desktop)
	fmt.Printf("Mobile Notifications:  %v\n", cfg.Notifications.Mobile.Enabled)
	if cfg.Notifications.Mobile.Enabled {
		fmt.Printf("Ntfy Topic:           %s\n", cfg.Notifications.Mobile.NtfyTopic)
	}
	fmt.Printf("\nFilters:              %d configured\n", len(cfg.Filters))
	fmt.Println("")
}

func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	cfg, err := filter.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch key {
	case "polling":
		interval, err := strconv.Atoi(value)
		if err != nil || interval < 10 {
			fmt.Println("❌ Polling interval must be a number >= 10")
			os.Exit(1)
		}
		cfg.PollingInterval = interval
		fmt.Printf("✅ Set polling interval to %d seconds\n", interval)

	case "desktop":
		if value == "true" || value == "1" || value == "yes" {
			cfg.Notifications.Desktop = true
			fmt.Println("✅ Desktop notifications enabled")
		} else if value == "false" || value == "0" || value == "no" {
			cfg.Notifications.Desktop = false
			fmt.Println("✅ Desktop notifications disabled")
		} else {
			fmt.Println("❌ Value must be true or false")
			os.Exit(1)
		}

	case "mobile":
		if value == "true" || value == "1" || value == "yes" {
			cfg.Notifications.Mobile.Enabled = true
			fmt.Println("✅ Mobile notifications enabled")
			if cfg.Notifications.Mobile.NtfyTopic == "" {
				fmt.Println("\n⚠️  Don't forget to set ntfy_topic:")
				fmt.Println("   email-sentinel config set ntfy_topic \"your-topic\"")
			}
		} else if value == "false" || value == "0" || value == "no" {
			cfg.Notifications.Mobile.Enabled = false
			fmt.Println("✅ Mobile notifications disabled")
		} else {
			fmt.Println("❌ Value must be true or false")
			os.Exit(1)
		}

	case "ntfy_topic":
		cfg.Notifications.Mobile.NtfyTopic = value
		fmt.Printf("✅ Set ntfy topic to: %s\n", value)

	default:
		fmt.Printf("❌ Unknown config key: %s\n", key)
		fmt.Println("\nAvailable keys: polling, desktop, mobile, ntfy_topic")
		os.Exit(1)
	}

	// Save config
	if err := filter.SaveConfig(cfg); err != nil {
		fmt.Printf("❌ Error saving config: %v\n", err)
		os.Exit(1)
	}
}
