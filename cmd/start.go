/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/config"
	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/gmail"
	"github.com/datateamsix/email-sentinel/internal/notify"
	"github.com/datateamsix/email-sentinel/internal/rules"
	"github.com/datateamsix/email-sentinel/internal/state"
	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/datateamsix/email-sentinel/internal/tray"
)

var daemonMode bool
var trayMode bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start monitoring Gmail for matching emails",
	Long: `Start monitoring your Gmail inbox for emails that match configured filters.

When a match is found, notifications are sent via:
- Desktop notifications (native OS)
- Mobile push notifications (via ntfy.sh, if configured)

The monitoring runs continuously, checking Gmail at regular intervals
defined in your configuration (default: 45 seconds).

Examples:
  # Run in foreground with logs
  email-sentinel start

  # Run as background daemon
  email-sentinel start --daemon`,
	Run: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&daemonMode, "daemon", "d", false, "Run as background daemon")
	startCmd.Flags().BoolVarP(&trayMode, "tray", "t", false, "Run with system tray icon")
}

func runStart(cmd *cobra.Command, args []string) {
	// Check if token exists
	if !gmail.TokenExists() {
		fmt.Println("❌ Not initialized. Run 'email-sentinel init' first.")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := filter.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(cfg.Filters) == 0 {
		fmt.Println("❌ No filters configured.")
		fmt.Println("\nAdd filters with: email-sentinel filter add")
		os.Exit(1)
	}

	// Load credentials
	credPath := findCredentials()
	if credPath == "" {
		fmt.Println("❌ credentials.json not found")
		fmt.Println("\nPlace credentials.json in:")
		fmt.Println("  - Current directory")
		configDir, _ := config.ConfigDir()
		fmt.Printf("  - Config directory: %s\n", configDir)
		os.Exit(1)
	}

	oauthConfig, err := gmail.LoadCredentials(credPath)
	if err != nil {
		fmt.Printf("❌ Error loading credentials: %v\n", err)
		os.Exit(1)
	}

	// Load token
	token, err := gmail.LoadToken()
	if err != nil {
		fmt.Printf("❌ Error loading token: %v\n", err)
		fmt.Println("\nRe-run: email-sentinel init")
		os.Exit(1)
	}

	// Create Gmail client
	client, err := gmail.NewClient(token, oauthConfig)
	if err != nil {
		fmt.Printf("❌ Error creating Gmail client: %v\n", err)
		os.Exit(1)
	}

	// Initialize seen messages tracker
	seenMessages, err := state.NewSeenMessages()
	if err != nil {
		fmt.Printf("❌ Error initializing state: %v\n", err)
		os.Exit(1)
	}

	// Initialize alert storage database
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("❌ Error initializing alert storage: %v\n", err)
		os.Exit(1)
	}
	defer storage.CloseDB(db)

	// Start daily cleanup scheduler (runs at 12:00 AM)
	stopCleanup := make(chan struct{})
	defer close(stopCleanup)
	go storage.StartDailyCleanup(db, stopCleanup)

	// Load priority rules
	rulesPath, err := rules.RulesPath()
	if err != nil {
		fmt.Printf("❌ Error getting rules path: %v\n", err)
		os.Exit(1)
	}

	priorityRules, err := rules.LoadRules(rulesPath)
	if err != nil {
		fmt.Printf("❌ Error loading priority rules: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Email Sentinel Started")
	fmt.Printf("   Monitoring %d filter(s)\n", len(cfg.Filters))
	fmt.Printf("   Polling interval: %d seconds\n", cfg.PollingInterval)
	if cfg.Notifications.Desktop {
		fmt.Println("   Desktop notifications: enabled")
	}
	if cfg.Notifications.Mobile.Enabled {
		fmt.Println("   Mobile notifications: enabled")
	}

	// Start system tray if requested
	if trayMode {
		fmt.Println("   System tray: enabled")
		fmt.Println("\n📱 Starting system tray... (Look for icon in taskbar)")
		fmt.Println("   Right-click tray icon for menu options")

		// Run tray in a goroutine - it blocks, so we run monitoring in main goroutine
		go func() {
			tray.Run(tray.Config{DB: db})
		}()

		// Give tray time to initialize
		time.Sleep(2 * time.Second)
	}

	fmt.Println("\n🔍 Watching for new emails... (Press Ctrl+C to stop)")
	fmt.Println("")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start monitoring loop
	ticker := time.NewTicker(time.Duration(cfg.PollingInterval) * time.Second)
	defer ticker.Stop()

	// Do initial check
	checkEmails(client, cfg, seenMessages, db, priorityRules)

	for {
		select {
		case <-ticker.C:
			checkEmails(client, cfg, seenMessages, db, priorityRules)
		case <-sigChan:
			fmt.Println("\n\n⏹️  Stopping Email Sentinel...")
			if trayMode {
				tray.Quit()
			}
			return
		}
	}
}

func checkEmails(client *gmail.Client, cfg *filter.Config, seenMessages *state.SeenMessages, db *sql.DB, priorityRules *rules.Rules) {
	// Fetch recent messages
	messages, err := client.GetRecentMessages(10)
	if err != nil {
		fmt.Printf("⚠️  Error fetching messages: %v\n", err)
		return
	}

	matchCount := 0

	for _, msg := range messages {
		// Skip if already seen
		if seenMessages.IsSeen(msg.Id) {
			continue
		}

		// Mark as seen immediately
		seenMessages.MarkSeen(msg.Id)

		// Parse message
		email := gmail.ParseMessage(msg)

		// Check against all filters (with metadata including labels)
		matchedFilters, err := filter.CheckAllFiltersWithMetadata(email.From, email.Subject)
		if err != nil {
			fmt.Printf("⚠️  Error checking filters: %v\n", err)
			continue
		}

		// If matched, send notifications
		if len(matchedFilters) > 0 {
			matchCount++

			for _, match := range matchedFilters {
				labelStr := ""
				if len(match.Labels) > 0 {
					labelStr = fmt.Sprintf(" 🏷️ %s", strings.Join(match.Labels, ", "))
				}
				fmt.Printf("📧 MATCH [%s]%s From: %s | Subject: %s\n",
					match.Name, labelStr, email.From, email.Subject)

				// Send desktop notification with labels
				if cfg.Notifications.Desktop {
					if err := notify.SendEmailAlertWithLabels(match.Name, match.Labels, email.From, email.Subject); err != nil {
						fmt.Printf("   ⚠️  Desktop notification failed: %v\n", err)
					}
				}

				// Send mobile notification with labels
				if cfg.Notifications.Mobile.Enabled && cfg.Notifications.Mobile.NtfyTopic != "" {
					if err := notify.SendMobileEmailAlertWithLabels(
						cfg.Notifications.Mobile.NtfyTopic,
						match.Name,
						match.Labels,
						email.From,
						email.Subject,
					); err != nil {
						fmt.Printf("   ⚠️  Mobile notification failed: %v\n", err)
					}
				}

				// Evaluate priority using rules engine
				msgMeta := rules.MessageMetadata{
					Sender:  email.From,
					Subject: email.Subject,
					Snippet: email.Snippet,
					Body:    "", // Body not available in snippet API call
				}
				priority := rules.EvaluatePriorityRules(priorityRules, msgMeta)

				// Save alert to database
				alert := &storage.Alert{
					Timestamp:    time.Now(),
					Sender:       email.From,
					Subject:      email.Subject,
					Snippet:      email.Snippet,
					Labels:       strings.Join(msg.LabelIds, ","),
					MessageID:    msg.Id,
					GmailLink:    gmail.BuildGmailLink(msg.Id),
					FilterName:   match.Name,
					FilterLabels: match.Labels,
					Priority:     priority,
				}

				if err := storage.InsertAlert(db, alert); err != nil {
					// Don't fail on storage errors, just log
					fmt.Printf("   ⚠️  Failed to save alert to database: %v\n", err)
				}

				// Send Windows toast notification (in addition to existing notifications)
				// This provides a rich, clickable notification in Windows Action Center
				if err := notify.SendAlertNotification(*alert); err != nil {
					fmt.Printf("   ⚠️  Toast notification failed: %v\n", err)
				}

				// Update system tray if enabled
				if trayMode {
					tray.UpdateTrayOnNewAlert(*alert)
				}
			}
		}
	}

	if matchCount == 0 {
		fmt.Printf("[%s] Checked %d messages, no new matches\n",
			time.Now().Format("15:04:05"), len(messages))
	}
}
