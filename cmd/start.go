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
	googlemail "google.golang.org/api/gmail/v1"

	"github.com/datateamsix/email-sentinel/internal/ai"
	"github.com/datateamsix/email-sentinel/internal/appconfig"
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
var cleanupInterval int // in minutes
var aiSummaryEnabled bool

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
	startCmd.Flags().IntVar(&cleanupInterval, "cleanup-interval", 60, "Auto-cleanup interval in minutes (0=disabled, default=60)")
	startCmd.Flags().BoolVar(&aiSummaryEnabled, "ai-summary", false, "Enable AI-powered email summaries")
}

func runStart(cmd *cobra.Command, args []string) {
	// Check if token exists
	if !gmail.TokenExists() {
		fmt.Println("❌ Not initialized. Run 'email-sentinel init' first.")
		os.Exit(1)
	}

	// Load unified configuration
	appCfg, err := appconfig.Load()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Load filter configuration (separate from app-config for now)
	cfg, err := filter.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading filter config: %v\n", err)
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

	// Run automatic backup on startup to ensure we have a recent backup
	storage.AutoBackupOnStartup(db)

	// Start daily cleanup scheduler (runs at 12:00 AM)
	stopCleanup := make(chan struct{})
	defer close(stopCleanup)
	go storage.StartDailyCleanup(db, stopCleanup)

	// Create priority rules from unified config
	priorityRules := &rules.Rules{
		PriorityRules: rules.PriorityRules{
			UrgentKeywords: appCfg.Priority.UrgentKeywords,
			VIPSenders:     appCfg.Priority.VIPSenders,
			VIPDomains:     appCfg.Priority.VIPDomains,
		},
		NotificationSettings: rules.NotificationSettings{
			QuietHoursStart: appCfg.Notifications.QuietHours.Start,
			QuietHoursEnd:   appCfg.Notifications.QuietHours.End,
			WeekendMode:     appCfg.Notifications.WeekendMode,
		},
	}

	// Initialize AI service if enabled via flag or config
	var aiService *ai.Service
	if aiSummaryEnabled || appCfg.AISummary.Enabled {
		// Create AI config from unified config
		aiConfig := createAIConfigFromAppConfig(appCfg)

		aiService, err = ai.NewService(aiConfig, db)
		if err != nil {
			fmt.Printf("⚠️  AI summary disabled: %v\n", err)
			fmt.Println("   Tip: Set API key environment variable (GEMINI_API_KEY, ANTHROPIC_API_KEY, or OPENAI_API_KEY)")
		}
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
	if aiService != nil {
		fmt.Println("   AI summaries: enabled")
		fmt.Printf("   AI provider: %s\n", appCfg.AISummary.Provider)
	}

	// Start system tray if requested
	if trayMode {
		fmt.Println("   System tray: enabled")
		if cleanupInterval > 0 {
			fmt.Printf("   Auto-cleanup: every %d minutes\n", cleanupInterval)
		} else {
			fmt.Println("   Auto-cleanup: disabled")
		}
		fmt.Println("\n📱 Starting system tray... (Look for icon in taskbar)")
		fmt.Println("   Right-click tray icon for menu options")

		// Run tray in a goroutine - it blocks, so we run monitoring in main goroutine
		go func() {
			tray.Run(tray.Config{
				DB:              db,
				CleanupInterval: time.Duration(cleanupInterval) * time.Minute,
			})
		}()

		// Give tray time to initialize
		time.Sleep(2 * time.Second)
	}

	fmt.Println("\n🔍 Watching for new emails... (Press Ctrl+C to stop)")
	fmt.Println("")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start monitoring loop with circuit breaker
	ticker := time.NewTicker(time.Duration(cfg.PollingInterval) * time.Second)
	defer ticker.Stop()

	// Circuit breaker state
	var (
		failureCount    int
		lastFailureTime time.Time
		backoffDuration = time.Duration(cfg.PollingInterval) * time.Second
	)

	// Do initial check
	if err := checkEmailsWithRecovery(client, cfg, seenMessages, db, priorityRules, aiService); err != nil {
		failureCount++
		lastFailureTime = time.Now()
	}

	for {
		select {
		case <-ticker.C:
			// Circuit breaker: implement exponential backoff on repeated failures
			if failureCount > 0 && time.Since(lastFailureTime) < backoffDuration {
				fmt.Printf("[%s] Backing off due to %d consecutive failures... waiting %v\n",
					time.Now().Format("15:04:05"), failureCount, backoffDuration)
				continue
			}

			// Attempt email check with recovery
			if err := checkEmailsWithRecovery(client, cfg, seenMessages, db, priorityRules, aiService); err != nil {
				failureCount++
				lastFailureTime = time.Now()

				// Exponential backoff: 45s, 90s, 180s, 360s (max 6 minutes)
				backoffDuration = time.Duration(cfg.PollingInterval*(1<<uint(min(failureCount-1, 3)))) * time.Second

				if failureCount >= 5 {
					fmt.Printf("\n❌ CRITICAL: %d consecutive Gmail API failures\n", failureCount)
					fmt.Printf("   Last error: %v\n", err)
					fmt.Printf("   Backing off for %v before next attempt\n", backoffDuration)
					fmt.Printf("   Check your network connection and Gmail API quota\n\n")
				}
			} else {
				// Success - reset circuit breaker
				if failureCount > 0 {
					fmt.Printf("[%s] ✅ Gmail API recovered after %d failures\n",
						time.Now().Format("15:04:05"), failureCount)
					failureCount = 0
					backoffDuration = time.Duration(cfg.PollingInterval) * time.Second
				}
			}

		case <-sigChan:
			fmt.Println("\n\n⏹️  Stopping Email Sentinel...")
			if trayMode {
				tray.Quit()
			}
			return
		}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// checkEmailsWithRecovery wraps checkEmails with panic recovery
func checkEmailsWithRecovery(client *gmail.Client, cfg *filter.Config, seenMessages *state.SeenMessages, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in checkEmails: %v", r)
			fmt.Printf("\n❌ PANIC RECOVERED in email checking: %v\n", r)
		}
	}()

	return checkEmails(client, cfg, seenMessages, db, priorityRules, aiService)
}

// createAIConfigFromAppConfig converts the unified AppConfig to the AI config format
func createAIConfigFromAppConfig(appCfg *appconfig.AppConfig) *ai.Config {
	return &ai.Config{
		AISummary: ai.AISummaryConfig{
			Enabled:  appCfg.AISummary.Enabled,
			Provider: appCfg.AISummary.Provider,
			API: ai.APIConfig{
				Claude: ai.ClaudeConfig{
					APIKey:      os.Getenv("ANTHROPIC_API_KEY"),
					Model:       appCfg.AISummary.Providers.Claude.Model,
					MaxTokens:   appCfg.AISummary.Providers.Claude.MaxTokens,
					Temperature: appCfg.AISummary.Providers.Claude.Temperature,
				},
				OpenAI: ai.OpenAIConfig{
					APIKey:      os.Getenv("OPENAI_API_KEY"),
					Model:       appCfg.AISummary.Providers.OpenAI.Model,
					MaxTokens:   appCfg.AISummary.Providers.OpenAI.MaxTokens,
					Temperature: appCfg.AISummary.Providers.OpenAI.Temperature,
				},
				Gemini: ai.GeminiConfig{
					APIKey:      os.Getenv("GEMINI_API_KEY"),
					Model:       appCfg.AISummary.Providers.Gemini.Model,
					MaxTokens:   appCfg.AISummary.Providers.Gemini.MaxTokens,
					Temperature: appCfg.AISummary.Providers.Gemini.Temperature,
				},
			},
			Behavior: ai.BehaviorConfig{
				EnableCache: appCfg.AISummary.Cache.Enabled,
				// Set defaults for fields not in new config
				MaxSummaryLength:       500,
				PriorityOnly:           false,
				TimeoutSeconds:         30,
				RetryAttempts:          3,
				IncludeInNotifications: true,
				ShowAIIcon:             true,
			},
			RateLimit: ai.RateLimitConfig{
				MaxPerHour: appCfg.AISummary.Providers.Gemini.RateLimit.RequestsPerMinute * 60,
				MaxPerDay:  appCfg.AISummary.Providers.Gemini.RateLimit.RequestsPerDay,
			},
			Prompt: ai.PromptConfig{
				System:       appCfg.AISummary.Prompt.System,
				UserTemplate: "Summarize this email:\n\nFrom: {{.From}}\nSubject: {{.Subject}}\n\n{{.Body}}",
			},
		},
	}
}

func checkEmails(client *gmail.Client, cfg *filter.Config, seenMessages *state.SeenMessages, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service) error {
	// Fetch recent messages
	messages, err := client.GetRecentMessages(10)
	if err != nil {
		fmt.Printf("⚠️  Error fetching messages: %v\n", err)
		return err
	}

	matchCount := 0

	for _, msg := range messages {
		// Skip if already seen
		if seenMessages.IsSeen(msg.Id) {
			continue
		}

		// Mark as seen immediately
		seenMessages.MarkSeen(msg.Id)

		// Process this message
		matched := processMessage(msg, cfg, db, priorityRules, aiService)
		if matched {
			matchCount++
		}
	}

	if matchCount == 0 {
		fmt.Printf("[%s] Checked %d messages, no new matches\n",
			time.Now().Format("15:04:05"), len(messages))
	}

	return nil
}

// processMessage processes a single email message and handles all matched filters
func processMessage(msg *googlemail.Message, cfg *filter.Config, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service) bool {
	// Parse message
	email := gmail.ParseMessage(msg)

	// Check against all filters (with metadata including labels)
	matchedFilters, err := filter.CheckAllFiltersWithMetadata(email.From, email.Subject)
	if err != nil {
		fmt.Printf("⚠️  Error checking filters: %v\n", err)
		return false
	}

	// If no matches, return early
	if len(matchedFilters) == 0 {
		return false
	}

	// Process each matched filter
	for _, match := range matchedFilters {
		processFilterMatch(msg, email, match, cfg, db, priorityRules, aiService)
	}

	return true
}

// processFilterMatch handles a single filter match including notifications and storage
func processFilterMatch(msg *googlemail.Message, email *gmail.EmailMessage, match filter.MatchResult, cfg *filter.Config, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service) {
	// Log the match
	labelStr := ""
	if len(match.Labels) > 0 {
		labelStr = fmt.Sprintf(" 🏷️ %s", strings.Join(match.Labels, ", "))
	}
	fmt.Printf("📧 MATCH [%s]%s From: %s | Subject: %s\n",
		match.Name, labelStr, email.From, email.Subject)

	// Send notifications (desktop and mobile)
	sendNotificationsForMatch(match, email, cfg)

	// Evaluate priority using rules engine
	priority := evaluateMessagePriority(email, priorityRules)

	// Create and save alert
	alert := createAlert(msg, email, match, priority)
	saveAndNotifyAlert(db, alert)

	// Generate AI summary asynchronously if enabled
	if aiService != nil {
		generateAISummaryAsync(aiService, *alert)
	}
}

// sendNotificationsForMatch sends desktop and mobile notifications for a matched filter
func sendNotificationsForMatch(match filter.MatchResult, email *gmail.EmailMessage, cfg *filter.Config) {
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
}

// evaluateMessagePriority determines the priority level of a message
func evaluateMessagePriority(email *gmail.EmailMessage, priorityRules *rules.Rules) int {
	msgMeta := rules.MessageMetadata{
		Sender:  email.From,
		Subject: email.Subject,
		Snippet: email.Snippet,
		Body:    "", // Body not available in snippet API call
	}
	return rules.EvaluatePriorityRules(priorityRules, msgMeta)
}

// createAlert creates an Alert struct from message data
func createAlert(msg *googlemail.Message, email *gmail.EmailMessage, match filter.MatchResult, priority int) *storage.Alert {
	return &storage.Alert{
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
}

// saveAndNotifyAlert saves an alert to the database and sends system notifications
func saveAndNotifyAlert(db *sql.DB, alert *storage.Alert) {
	// Save alert with retry logic to prevent data loss
	if err := storage.InsertAlertWithRetry(db, alert); err != nil {
		// Critical: Even retry and fallback failed
		fmt.Printf("   ❌ CRITICAL: Failed to save alert (retry + fallback failed): %v\n", err)
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

// generateAISummaryAsync generates an AI summary in a separate goroutine with panic recovery
func generateAISummaryAsync(aiService *ai.Service, alert storage.Alert) {
	go func(alertCopy storage.Alert) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("   ❌ PANIC in AI summary goroutine: %v\n", r)
				fmt.Printf("      Alert: %s from %s\n", alertCopy.Subject, alertCopy.Sender)
			}
		}()

		summary, err := aiService.GenerateSummary(
			alertCopy.MessageID,
			alertCopy.Sender,
			alertCopy.Subject,
			"", // body not available in snippet API
			alertCopy.Snippet,
			alertCopy.Priority,
		)
		if err != nil {
			fmt.Printf("   ⚠️  AI summary failed: %v\n", err)
			return
		}
		if summary != nil {
			fmt.Printf("   🤖 AI: %s\n", summary.Summary)
		}
	}(alert)
}
