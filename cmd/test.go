/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/notify"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test notifications and configuration",
	Long: `Test that desktop and mobile notifications are working correctly.

This command will send test notifications to verify your system is configured
properly before you start monitoring.

Subcommands:
  desktop     Test desktop notification
  mobile      Test mobile notification (requires ntfy_topic configured)
  toast       Test Windows toast notification (Windows only)
  filter      Test if an email would match a filter

Examples:
  email-sentinel test desktop
  email-sentinel test mobile
  email-sentinel test toast
  email-sentinel test toast --priority  (test high-priority notification)
  email-sentinel test filter "Job Alerts" "from:linkedin.com" "subject:interview"`,
}

var testDesktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Send a test desktop notification",
	Long: `Send a test desktop notification to verify your system permissions.

This will attempt to display a native OS notification. If you don't see it:
- Windows: Check Settings → System → Notifications
- macOS: Check System Preferences → Notifications → Terminal
- Linux: Ensure notification daemon is running (notify-send)`,
	Run: runTestDesktop,
}

var testMobileCmd = &cobra.Command{
	Use:   "mobile",
	Short: "Send a test mobile notification",
	Long: `Send a test mobile push notification via ntfy.sh.

Requires:
- Mobile notifications enabled: email-sentinel config set mobile true
- Topic configured: email-sentinel config set ntfy_topic "your-topic"
- ntfy app installed on your phone subscribed to the topic`,
	Run: runTestMobile,
}

var testToastCmd = &cobra.Command{
	Use:   "toast",
	Short: "Send a test Windows toast notification",
	Long: `Send a test Windows toast notification to verify Action Center notifications.

This sends a rich, clickable notification that appears in Windows Action Center.
The notification includes:
  • Clickable link to open in browser
  • Email preview with sender and snippet
  • Priority styling for urgent alerts

Windows only. On other platforms, use 'test desktop' instead.

Examples:
  email-sentinel test toast               # Test normal priority
  email-sentinel test toast --priority    # Test high priority`,
	Run: runTestToast,
}

var testFilterCmd = &cobra.Command{
	Use:   "filter <filter-name> <from> <subject>",
	Short: "Test if an email would match a filter",
	Long: `Test if a hypothetical email would match a specific filter.

This is useful for validating your filter patterns before real emails arrive.

Example:
  email-sentinel test filter "Job Alerts" "recruiter@linkedin.com" "New job opportunity"`,
	Args: cobra.ExactArgs(3),
	Run:  runTestFilter,
}

var testPriority bool

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testDesktopCmd)
	testCmd.AddCommand(testMobileCmd)
	testCmd.AddCommand(testToastCmd)
	testCmd.AddCommand(testFilterCmd)

	// Add priority flag to toast test
	testToastCmd.Flags().BoolVarP(&testPriority, "priority", "p", false, "Test high-priority notification")
}

func runTestDesktop(cmd *cobra.Command, args []string) {
	fmt.Println("🔔 Sending test desktop notification...")
	fmt.Println("")

	err := notify.SendDesktopNotification(
		"Email Sentinel Test",
		"If you can see this, desktop notifications are working! ✅",
	)

	if err != nil {
		fmt.Printf("❌ Desktop notification failed: %v\n", err)
		fmt.Println("")
		fmt.Println("Troubleshooting:")
		fmt.Println("  Windows: Settings → System → Notifications → Enable notifications")
		fmt.Println("  macOS: System Preferences → Notifications → Terminal → Allow")
		fmt.Println("  Linux: Install libnotify-bin or equivalent notification daemon")
		os.Exit(1)
	}

	fmt.Println("✅ Test notification sent!")
	fmt.Println("")
	fmt.Println("Did you see a notification pop up?")
	fmt.Println("If not, check your system's notification settings:")
	fmt.Println("")
	fmt.Println("Windows:")
	fmt.Println("  • Settings → System → Notifications → Notifications")
	fmt.Println("  • Turn off 'Focus Assist' (Action Center)")
	fmt.Println("")
	fmt.Println("macOS:")
	fmt.Println("  • System Preferences → Notifications")
	fmt.Println("  • Find 'Terminal' or your terminal app")
	fmt.Println("  • Enable 'Allow Notifications'")
	fmt.Println("  • Turn off 'Do Not Disturb'")
	fmt.Println("")
	fmt.Println("Linux:")
	fmt.Println("  • Check notification daemon is running:")
	fmt.Println("    systemctl --user status dunst (or notification-daemon)")
	fmt.Println("  • Test with: notify-send 'Test' 'Message'")
}

func runTestMobile(cmd *cobra.Command, args []string) {
	fmt.Println("📱 Sending test mobile notification...")
	fmt.Println("")

	cfg, err := filter.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	if !cfg.Notifications.Mobile.Enabled {
		fmt.Println("❌ Mobile notifications are disabled")
		fmt.Println("\nEnable with: email-sentinel config set mobile true")
		os.Exit(1)
	}

	if cfg.Notifications.Mobile.NtfyTopic == "" {
		fmt.Println("❌ No ntfy topic configured")
		fmt.Println("\nSet topic: email-sentinel config set ntfy_topic \"your-topic\"")
		os.Exit(1)
	}

	fmt.Printf("Sending to topic: %s\n", cfg.Notifications.Mobile.NtfyTopic)
	fmt.Println("")

	err = notify.SendMobileNotification(
		cfg.Notifications.Mobile.NtfyTopic,
		"Email Sentinel Test",
		"If you can see this on your phone, mobile notifications are working! ✅",
	)

	if err != nil {
		fmt.Printf("❌ Mobile notification failed: %v\n", err)
		fmt.Println("")
		fmt.Println("Troubleshooting:")
		fmt.Println("  1. Verify ntfy app is installed on your phone")
		fmt.Println("  2. Check you're subscribed to topic:", cfg.Notifications.Mobile.NtfyTopic)
		fmt.Println("  3. Test manually: https://ntfy.sh/" + cfg.Notifications.Mobile.NtfyTopic)
		os.Exit(1)
	}

	fmt.Println("✅ Test notification sent!")
	fmt.Println("")
	fmt.Println("Check your phone for a notification from ntfy.sh")
	fmt.Println("")
	fmt.Println("If you didn't receive it:")
	fmt.Println("  • Open ntfy app and verify subscription to:", cfg.Notifications.Mobile.NtfyTopic)
	fmt.Println("  • Check phone's notification settings allow ntfy notifications")
	fmt.Println("  • Try a different topic name (must be unique)")
}

func runTestToast(cmd *cobra.Command, args []string) {
	fmt.Println("🪟 Sending test Windows toast notification...")
	fmt.Println("")

	var err error
	if testPriority {
		fmt.Println("Testing HIGH PRIORITY notification...")
		err = notify.SendPriorityTestNotification()
	} else {
		fmt.Println("Testing normal priority notification...")
		err = notify.SendTestNotification()
	}

	if err != nil {
		fmt.Printf("❌ Toast notification failed: %v\n", err)
		fmt.Println("")
		fmt.Println("Troubleshooting:")
		fmt.Println("  • This feature is Windows-only")
		fmt.Println("  • Check Settings → System → Notifications")
		fmt.Println("  • Ensure notifications are enabled for Email Sentinel")
		fmt.Println("  • Turn off Focus Assist in Action Center")
		os.Exit(1)
	}

	fmt.Println("✅ Toast notification sent!")
	fmt.Println("")
	fmt.Println("Check your Windows Action Center for the notification.")
	fmt.Println("Features to verify:")
	fmt.Println("  ✓ Notification shows email subject as title")
	fmt.Println("  ✓ Shows sender and email preview")
	fmt.Println("  ✓ Click 'Open Email' button to test link")
	if testPriority {
		fmt.Println("  ✓ Shows 🔥 icon and HIGH PRIORITY label")
		fmt.Println("  ✓ Uses reminder audio (more urgent)")
	} else {
		fmt.Println("  ✓ Shows 📧 icon and filter name")
	}
	fmt.Println("")
	fmt.Println("Tip: Try with --priority flag to test urgent notifications")
}

func runTestFilter(cmd *cobra.Command, args []string) {
	filterName := args[0]
	fromEmail := args[1]
	subjectLine := args[2]

	filters, err := filter.ListFilters()
	if err != nil {
		fmt.Printf("❌ Error loading filters: %v\n", err)
		os.Exit(1)
	}

	// Find the filter
	var targetFilter *filter.Filter
	for _, f := range filters {
		if f.Name == filterName {
			targetFilter = &f
			break
		}
	}

	if targetFilter == nil {
		fmt.Printf("❌ Filter '%s' not found\n", filterName)
		fmt.Println("\nAvailable filters:")
		for _, f := range filters {
			fmt.Printf("  - %s\n", f.Name)
		}
		os.Exit(1)
	}

	fmt.Printf("🧪 Testing Filter: %s\n", filterName)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	fmt.Printf("Email From:    %s\n", fromEmail)
	fmt.Printf("Email Subject: %s\n", subjectLine)
	fmt.Println("")

	matches := filter.MatchesFilter(*targetFilter, fromEmail, subjectLine)

	if matches {
		fmt.Println("✅ MATCH - This email would trigger a notification!")
	} else {
		fmt.Println("❌ NO MATCH - This email would be ignored")
		fmt.Println("")
		fmt.Println("Filter details:")
		if len(targetFilter.From) > 0 {
			fmt.Printf("  From patterns: %v\n", targetFilter.From)
		}
		if len(targetFilter.Subject) > 0 {
			fmt.Printf("  Subject patterns: %v\n", targetFilter.Subject)
		}
		fmt.Printf("  Match mode: %s\n", targetFilter.Match)
	}
}
