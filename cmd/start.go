/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	googlemail "google.golang.org/api/gmail/v1"

	"github.com/datateamsix/email-sentinel/internal/accounts"
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
var searchScope string // Gmail search scope (inbox, all, all-except-trash, spam-only)

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

Gmail Scope:
Each filter can specify which Gmail categories to search (inbox, primary,
social, promotions, etc.). The --search flag overrides all per-filter scopes.

Examples:
  # Run in foreground with logs (uses per-filter scopes)
  email-sentinel start

  # Run with system tray
  email-sentinel start --tray

  # Override all filters to search only social category
  email-sentinel start --search social

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
	startCmd.Flags().StringVar(&searchScope, "search", "", "Override filter scopes with global search: inbox, all, primary, social, promotions, updates, forums, all-except-trash")
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
		fmt.Println("⚠️  No filters configured yet.")
		fmt.Println("\n📝 You can add filters in several ways:")
		fmt.Println("   • Command line: email-sentinel filter add")
		fmt.Println("   • Interactive menu: email-sentinel menu")
		if trayMode {
			fmt.Println("   • System tray: Right-click the tray icon > Manage Filters")
		}
		fmt.Println("\n▶️  Starting monitoring service (no matches will trigger until you add filters)...")
		fmt.Println()
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

	// Build Gmail search query from scope flag (if provided)
	var gmailSearchQuery string
	if searchScope != "" {
		gmailSearchQuery = buildGmailSearchQuery(searchScope)
		fmt.Printf("   Global search override: %s (query: '%s')\n", searchScope, gmailSearchQuery)
	} else {
		fmt.Println("   Using per-filter Gmail scopes")
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
	if err := checkEmailsWithRecovery(client, cfg, seenMessages, db, priorityRules, aiService, gmailSearchQuery); err != nil {
		failureCount++
		lastFailureTime = time.Now()
	}

	for {
		select {
		case <-ticker.C:
			// Check for expired filters and clean them up
			removed, err := filter.CleanupExpiredFilters()
			if err != nil {
				fmt.Printf("⚠️  Error checking for expired filters: %v\n", err)
			} else if len(removed) > 0 {
				for _, name := range removed {
					fmt.Printf("🗑️  Filter '%s' expired and was automatically removed\n", name)
					// Send notification about expired filter
					notify.SendDesktopNotification(
						"Filter Expired",
						fmt.Sprintf("Filter '%s' has expired and been removed", name),
					)
				}
				// Reload config since filters were removed
				cfg, err = filter.LoadConfig()
				if err != nil {
					fmt.Printf("⚠️  Error reloading config after cleanup: %v\n", err)
				}
			}

			// Check for expiring trials and send alerts
			checkExpiringTrials(db)

			// Circuit breaker: implement exponential backoff on repeated failures
			if failureCount > 0 && time.Since(lastFailureTime) < backoffDuration {
				fmt.Printf("[%s] Backing off due to %d consecutive failures... waiting %v\n",
					time.Now().Format("15:04:05"), failureCount, backoffDuration)
				continue
			}

			// Attempt email check with recovery
			if err := checkEmailsWithRecovery(client, cfg, seenMessages, db, priorityRules, aiService, gmailSearchQuery); err != nil {
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
// buildGmailSearchQuery converts a search scope string to a Gmail search query
func buildGmailSearchQuery(scope string) string {
	switch strings.ToLower(scope) {
	case "all":
		return "" // Empty query = search everything
	case "all-except-trash":
		return "-in:trash"
	case "spam-only":
		return "in:spam"
	case "promotions":
		return "category:promotions"
	case "social":
		return "category:social"
	case "updates":
		return "category:updates"
	case "forums":
		return "category:forums"
	case "inbox":
		return "in:inbox"
	default:
		// Default to inbox if unknown scope
		fmt.Printf("⚠️  Unknown search scope '%s', defaulting to 'inbox'\n", scope)
		return "in:inbox"
	}
}

func checkEmailsWithRecovery(client *gmail.Client, cfg *filter.Config, seenMessages *state.SeenMessages, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service, searchQuery string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in checkEmails: %v", r)
			fmt.Printf("\n❌ PANIC RECOVERED in email checking: %v\n", r)
		}
	}()

	return checkEmails(client, cfg, seenMessages, db, priorityRules, aiService, searchQuery)
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

func checkEmails(client *gmail.Client, cfg *filter.Config, seenMessages *state.SeenMessages, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service, searchQuery string) error {
	// Get all unique scopes from filters for optimized fetching
	uniqueScopes, err := filter.GetAllUniqueScopes()
	if err != nil {
		fmt.Printf("⚠️  Error getting filter scopes: %v\n", err)
		return err
	}

	// If global search query is provided (via --search flag), use it
	// Otherwise, fetch messages for each unique scope
	var allMessages []*googlemail.Message
	var fetchErr error

	if searchQuery != "" {
		// Global scope override from command line flag
		allMessages, fetchErr = client.GetRecentMessagesWithQuery(10, searchQuery)
	} else {
		// Fetch messages for each unique filter scope
		messageMap := make(map[string]*googlemail.Message)
		for _, scope := range uniqueScopes {
			query := filter.BuildGmailSearchQuery(scope)
			messages, err := client.GetRecentMessagesWithQuery(10, query)
			if err != nil {
				fmt.Printf("⚠️  Error fetching messages for scope '%s': %v\n", scope, err)
				fetchErr = err
				continue
			}

			// Deduplicate messages by ID
			for _, msg := range messages {
				messageMap[msg.Id] = msg
			}
		}

		// Convert map to slice
		allMessages = make([]*googlemail.Message, 0, len(messageMap))
		for _, msg := range messageMap {
			allMessages = append(allMessages, msg)
		}
	}

	if fetchErr != nil {
		return fetchErr
	}

	matchCount := 0

	for _, msg := range allMessages {
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
			time.Now().Format("15:04:05"), len(allMessages))
	}

	return nil
}

// processMessage processes a single email message and handles all matched filters
func processMessage(msg *googlemail.Message, cfg *filter.Config, db *sql.DB, priorityRules *rules.Rules, aiService *ai.Service) bool {
	// Parse message
	email := gmail.ParseMessage(msg)

	// Detect digital accounts (subscriptions, trials, etc.) - runs on ALL emails
	detectAndSaveAccount(email, db)

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
	saveAndNotifyAlert(db, alert, cfg)

	// Generate AI summary asynchronously if enabled
	if aiService != nil {
		generateAISummaryAsync(aiService, *alert)
	}
}

// sendNotificationsForMatch sends mobile notifications for a matched filter
// Desktop notifications are handled by saveAndNotifyAlert() to avoid duplicates
func sendNotificationsForMatch(match filter.MatchResult, email *gmail.EmailMessage, cfg *filter.Config) {
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
func saveAndNotifyAlert(db *sql.DB, alert *storage.Alert, cfg *filter.Config) {
	// Save alert with retry logic to prevent data loss
	if err := storage.InsertAlertWithRetry(db, alert); err != nil {
		// Critical: Even retry and fallback failed
		fmt.Printf("   ❌ CRITICAL: Failed to save alert (retry + fallback failed): %v\n", err)
	}

	// Send desktop notification (Windows toast or Unix notification) if enabled
	// This provides a rich, platform-specific notification with AI summaries
	if cfg.Notifications.Desktop {
		if err := notify.SendAlertNotification(*alert); err != nil {
			fmt.Printf("   ⚠️  Desktop notification failed: %v\n", err)
		}
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

// detectAndSaveAccount detects and saves digital account information from emails
func detectAndSaveAccount(email *gmail.EmailMessage, db *sql.DB) {
	// Load app config to get account settings
	appCfg, err := appconfig.Load()
	if err != nil || !appCfg.Accounts.Enabled {
		// Silently skip if config not available or accounts disabled
		return
	}

	// Load account configuration
	accountCfg := accounts.LoadConfigFromAppConfig(appCfg)
	if !accountCfg.Enabled {
		return
	}

	// Create detector
	detector := accounts.NewDetector(accountCfg.MinConfidence, accountCfg.Categories)

	// Create detection context
	ctx := accounts.DetectionContext{
		Subject:      email.Subject,
		Body:         "",          // Body not available in snippet API
		Snippet:      email.Snippet,
		Sender:       email.From,
		ToEmail:      "",          // We'll try to extract this
		ReceivedDate: time.Now(),  // Use current time as we don't have exact received date
		MessageID:    email.ID,    // Use Gmail message ID
	}

	// Attempt to extract recipient email from snippet
	if ctx.ToEmail == "" {
		// Try to get from Gmail headers if available
		ctx.ToEmail = extractRecipientFromEmail(email)
	}

	// Detect account
	result, err := detector.DetectAccount(ctx)
	if err != nil {
		// Silent failure - don't spam logs for detection errors
		return
	}

	if result == nil {
		// No account detected - this is normal, not an error
		return
	}

	// Convert to storage model
	now := time.Now()
	account := &storage.Account{
		ServiceName:    result.ServiceName,
		EmailAddress:   result.EmailAddress,
		AccountType:    result.AccountType,
		Status:         "active",
		PriceMonthly:   result.PriceMonthly,
		TrialEndDate:   result.TrialEndDate,
		GmailMessageID: result.GmailMessageID,
		DetectedAt:     now,
		UpdatedAt:      now,
		Confidence:     result.Confidence,
		CancelURL:      result.CancelURL,
		Category:       result.Category,
	}

	// Save to database
	if err := storage.InsertAccount(db, account); err != nil {
		// Only log if it's not a duplicate
		if !strings.Contains(err.Error(), "UNIQUE") {
			fmt.Printf("   ⚠️  Failed to save account: %v\n", err)
		}
		return
	}

	// Log successful detection
	typeIcon := "💳"
	if account.AccountType == "trial" {
		typeIcon = "🆓"
	} else if account.AccountType == "free" {
		typeIcon = "🎁"
	}

	fmt.Printf("   %s ACCOUNT DETECTED: %s (%s) | Email: %s\n",
		typeIcon,
		account.ServiceName,
		account.AccountType,
		account.EmailAddress,
	)

	if account.TrialEndDate != nil {
		daysUntil := time.Until(*account.TrialEndDate).Hours() / 24
		if daysUntil > 0 {
			fmt.Printf("      Trial expires in %d days\n", int(daysUntil)+1)
		}
	}

	if account.PriceMonthly > 0 {
		fmt.Printf("      Price: $%.2f/month\n", account.PriceMonthly)
	}
}

// extractRecipientFromEmail attempts to extract the recipient email address
func extractRecipientFromEmail(email *gmail.EmailMessage) string {
	// Try to extract from snippet (look for "sent to:", "delivered to:", etc.)
	text := email.Subject + " " + email.Snippet

	// Common patterns
	patterns := []string{
		`(?i)sent to:?\s*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
		`(?i)delivered to:?\s*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
		`(?i)for:?\s*([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
	}

	for _, patternStr := range patterns {
		if matches := regexp.MustCompile(patternStr).FindStringSubmatch(text); len(matches) > 1 {
			return matches[1]
		}
	}

	// If not found, return empty - will default to user's primary email
	return ""
}

// checkExpiringTrials checks for expiring trials and sends alerts
func checkExpiringTrials(db *sql.DB) {
	// Load app config to get trial alert settings
	appCfg, err := appconfig.Load()
	if err != nil || !appCfg.Accounts.Enabled {
		return
	}

	// Get all active trials
	trials, err := storage.GetActiveTrials(db)
	if err != nil {
		// Silent failure - don't spam logs
		return
	}

	// Check each trial against alert thresholds
	for _, trial := range trials {
		if trial.TrialEndDate == nil {
			continue
		}

		daysUntil := time.Until(*trial.TrialEndDate).Hours() / 24

		// Skip if already expired
		if daysUntil < 0 {
			continue
		}

		// Check against each alert threshold
		for _, alert := range appCfg.Accounts.TrialAlerts {
			if daysUntil <= float64(alert.DaysBefore) && daysUntil > float64(alert.DaysBefore-1) {
				// Should send alert
				sendTrialExpirationAlert(trial, alert.Urgency, int(daysUntil)+1)
			}
		}
	}
}

// sendTrialExpirationAlert sends a notification for an expiring trial
func sendTrialExpirationAlert(trial storage.Account, urgency string, daysUntil int) {
	var title, message, icon string

	switch urgency {
	case "critical":
		icon = "🔥"
		title = "FINAL WARNING: Trial Expires Soon!"
	case "high":
		icon = "⚠️"
		title = "Trial Expires Soon"
	default:
		icon = "📅"
		title = "Trial Expiring"
	}

	if daysUntil == 1 {
		message = fmt.Sprintf("%s %s trial expires tomorrow", icon, trial.ServiceName)
	} else {
		message = fmt.Sprintf("%s %s trial expires in %d days", icon, trial.ServiceName, daysUntil)
	}

	if trial.PriceMonthly > 0 {
		message += fmt.Sprintf(" ($%.2f/month)", trial.PriceMonthly)
	}

	// Log to console
	fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)

	// Send desktop notification
	if err := notify.SendDesktopNotification(title, message); err != nil {
		// Silent failure for notifications
		return
	}
}
