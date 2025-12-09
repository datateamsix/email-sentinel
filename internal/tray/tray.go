package tray

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"fyne.io/systray"
)

// TrayApp represents the system tray application
type TrayApp struct {
	db              *sql.DB
	recentAlerts    []*systray.MenuItem
	alertUpdateChan chan storage.Alert
	quitChan        chan struct{}
	mu              sync.Mutex
	hasUrgent       bool
	refreshTimer    *time.Timer
	refreshMu       sync.Mutex
	cleanupInterval time.Duration
}

// Config holds configuration for the tray app
type Config struct {
	DB              *sql.DB
	CleanupInterval time.Duration // How often to cleanup old alerts (0 = disabled)
}

var (
	globalApp      *TrayApp
	mRecentAlerts  *systray.MenuItem
	mOpenHistory   *systray.MenuItem
	mClearAlerts   *systray.MenuItem
	mQuit          *systray.MenuItem
)

// Run starts the system tray application
// This function blocks until the tray is quit
func Run(cfg Config) {
	globalApp = &TrayApp{
		db:              cfg.DB,
		alertUpdateChan: make(chan storage.Alert, 100),
		quitChan:        make(chan struct{}),
		recentAlerts:    make([]*systray.MenuItem, 0),
		cleanupInterval: cfg.CleanupInterval,
	}

	systray.Run(onReady, onExit)
}

// onReady is called when the system tray is ready
func onReady() {
	// Set initial icon and title (only if valid icon data exists)
	if icon := GetNormalIcon(); icon != nil && len(icon) > 0 {
		systray.SetIcon(icon)
	}
	systray.SetTitle("Email Sentinel")
	systray.SetTooltip("Email Sentinel - Monitoring Gmail")

	// Create menu items
	mRecentAlerts = systray.AddMenuItem("Recent Alerts", "View recent email alerts")
	systray.AddSeparator()
	mOpenHistory = systray.AddMenuItem("Open History", "View all alerts in terminal")
	mClearAlerts = systray.AddMenuItem("Clear Alerts", "Delete all alerts from history")
	systray.AddSeparator()
	mQuit = systray.AddMenuItem("Quit", "Quit Email Sentinel")

	// Load initial alerts
	go globalApp.loadRecentAlerts()

	// Start event handlers
	go globalApp.handleMenuEvents()
	go globalApp.handleAlertUpdates()

	log.Println("📱 System tray initialized")
}

// onExit is called when the system tray is exiting
func onExit() {
	log.Println("🛑 System tray shutting down")
	close(globalApp.quitChan)
}

// scheduleRefresh schedules a debounced refresh of the alerts menu
// Multiple calls within 500ms will be batched into a single refresh
func (app *TrayApp) scheduleRefresh() {
	app.refreshMu.Lock()
	defer app.refreshMu.Unlock()

	// Reset the timer if it exists, otherwise create a new one
	if app.refreshTimer != nil {
		app.refreshTimer.Stop()
	}

	app.refreshTimer = time.AfterFunc(500*time.Millisecond, func() {
		app.loadRecentAlerts()
	})
}

// loadRecentAlerts loads the 10 most recent alerts from the database
func (app *TrayApp) loadRecentAlerts() {
	app.mu.Lock()
	defer app.mu.Unlock()

	// Clear existing submenu items
	for _, item := range app.recentAlerts {
		item.Hide()
	}
	app.recentAlerts = make([]*systray.MenuItem, 0)

	// Fetch recent alerts from database
	alerts, err := storage.GetRecentAlerts(app.db, 10)
	if err != nil {
		log.Printf("Error loading recent alerts: %v", err)
		noAlerts := mRecentAlerts.AddSubMenuItem("No alerts yet", "")
		noAlerts.Disable()
		return
	}

	if len(alerts) == 0 {
		noAlerts := mRecentAlerts.AddSubMenuItem("No alerts yet", "")
		noAlerts.Disable()
		return
	}

	// Check if any are urgent
	hasUrgent := false
	for _, alert := range alerts {
		if alert.Priority == 1 {
			hasUrgent = true
			break
		}
	}
	app.hasUrgent = hasUrgent

	// Update icon based on urgent status (only if valid icon data exists)
	if hasUrgent {
		if icon := GetUrgentIcon(); icon != nil && len(icon) > 0 {
			systray.SetIcon(icon)
		}
		systray.SetTooltip("Email Sentinel - ⚠️ Urgent alerts!")
	} else {
		if icon := GetNormalIcon(); icon != nil && len(icon) > 0 {
			systray.SetIcon(icon)
		}
		systray.SetTooltip("Email Sentinel - Monitoring Gmail")
	}

	// Add each alert as a submenu item
	for _, alert := range alerts {
		app.addAlertMenuItem(alert)
	}
}

// addAlertMenuItem adds a single alert to the recent alerts submenu
func (app *TrayApp) addAlertMenuItem(alert storage.Alert) {
	// Determine icon based on priority and filter labels
	icon := "📧"

	// Check if this is an OTP-related alert
	isOTP := false
	for _, label := range alert.FilterLabels {
		if label == "otp" || label == "OTP" {
			isOTP = true
			break
		}
	}

	// Check if this alert has an AI summary
	hasAISummary := alert.AISummary != nil

	// Set icon: OTP takes precedence, then AI summary, then priority, then default
	if isOTP {
		icon = "🔐" // Lock icon for OTP messages
	} else if hasAISummary {
		icon = "🤖" // AI icon for summarized emails
	} else if alert.Priority == 1 {
		icon = "🔥" // Fire icon for high priority
	}

	// Truncate subject if too long
	subject := alert.Subject
	if len(subject) > 50 {
		subject = subject[:47] + "..."
	}

	// Format time
	timeStr := alert.Timestamp.Format("15:04")
	if !isToday(alert.Timestamp) {
		timeStr = alert.Timestamp.Format("Jan 2")
	}

	title := fmt.Sprintf("%s [%s] %s", icon, timeStr, subject)

	// Enhanced tooltip with filter info and AI summary
	tooltip := fmt.Sprintf("From: %s\nFilter: %s\nClick to open in Gmail", alert.Sender, alert.FilterName)
	if isOTP {
		tooltip = fmt.Sprintf("🔐 OTP Message\nFrom: %s\nFilter: %s\nClick to open in Gmail", alert.Sender, alert.FilterName)
	}

	// Add AI summary to tooltip if available
	if hasAISummary && alert.AISummary != nil {
		tooltip += fmt.Sprintf("\n\n🤖 AI Summary:\n%s", alert.AISummary.Summary)

		if len(alert.AISummary.Questions) > 0 {
			tooltip += fmt.Sprintf("\n\n❓ Questions (%d):", len(alert.AISummary.Questions))
			for i, q := range alert.AISummary.Questions {
				if i < 3 { // Show max 3 questions in tooltip
					tooltip += fmt.Sprintf("\n  • %s", q)
				}
			}
			if len(alert.AISummary.Questions) > 3 {
				tooltip += fmt.Sprintf("\n  ... and %d more", len(alert.AISummary.Questions)-3)
			}
		}

		if len(alert.AISummary.ActionItems) > 0 {
			tooltip += fmt.Sprintf("\n\n✅ Action Items (%d):", len(alert.AISummary.ActionItems))
			for i, item := range alert.AISummary.ActionItems {
				if i < 3 { // Show max 3 action items in tooltip
					tooltip += fmt.Sprintf("\n  • %s", item)
				}
			}
			if len(alert.AISummary.ActionItems) > 3 {
				tooltip += fmt.Sprintf("\n  ... and %d more", len(alert.AISummary.ActionItems)-3)
			}
		}
	}

	menuItem := mRecentAlerts.AddSubMenuItem(title, tooltip)
	app.recentAlerts = append(app.recentAlerts, menuItem)

	// Handle clicks on this alert (open Gmail link)
	go func(link string, item *systray.MenuItem) {
		for {
			select {
			case <-item.ClickedCh:
				openBrowser(link)
			case <-app.quitChan:
				return
			}
		}
	}(alert.GmailLink, menuItem)
}

// handleMenuEvents handles clicks on main menu items
func (app *TrayApp) handleMenuEvents() {
	for {
		select {
		case <-mOpenHistory.ClickedCh:
			app.openHistory()

		case <-mClearAlerts.ClickedCh:
			app.clearAlerts()

		case <-mQuit.ClickedCh:
			log.Println("Quit requested from tray menu")
			systray.Quit()
			return

		case <-app.quitChan:
			return
		}
	}
}

// handleAlertUpdates processes new alerts sent via UpdateTrayOnNewAlert
func (app *TrayApp) handleAlertUpdates() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Setup cleanup ticker if enabled (interval > 0)
	var cleanupTicker *time.Ticker
	var cleanupChan <-chan time.Time

	if app.cleanupInterval > 0 {
		cleanupTicker = time.NewTicker(app.cleanupInterval)
		defer cleanupTicker.Stop()
		cleanupChan = cleanupTicker.C
		log.Printf("🗑️  Auto-cleanup enabled: checking every %v", app.cleanupInterval)
	} else {
		log.Println("🗑️  Auto-cleanup disabled")
	}

	for {
		select {
		case alert := <-app.alertUpdateChan:
			log.Printf("📱 Tray: New alert received - %s", alert.Subject)

			// Temporarily switch to urgent icon if it's a priority alert
			if alert.Priority == 1 {
				app.mu.Lock()
				app.hasUrgent = true
				app.mu.Unlock()

				if icon := GetUrgentIcon(); icon != nil && len(icon) > 0 {
					systray.SetIcon(icon)
					systray.SetTooltip("Email Sentinel - ⚠️ New urgent alert!")
				}

				// Schedule a refresh after icon flash (debounced)
				go func() {
					time.Sleep(5 * time.Second)
					app.scheduleRefresh()
				}()
			}

			// Schedule refresh for the new alert (debounced)
			app.scheduleRefresh()

		case <-ticker.C:
			// Periodically refresh the alerts (debounced)
			app.scheduleRefresh()

		case <-cleanupChan:
			// Delete alerts older than 24 hours (only if cleanup is enabled)
			deleted, err := storage.DeleteAlerts24HoursOld(app.db)
			if err != nil {
				log.Printf("⚠️  Error cleaning up 24-hour-old alerts: %v", err)
			} else if deleted > 0 {
				log.Printf("🗑️  Cleaned up %d alert(s) older than 24 hours", deleted)
				app.scheduleRefresh()
			}

		case <-app.quitChan:
			return
		}
	}
}

// UpdateTrayOnNewAlert is called when a new alert is created
// This updates the tray icon and menu with the new alert
func UpdateTrayOnNewAlert(alert storage.Alert) {
	if globalApp != nil {
		select {
		case globalApp.alertUpdateChan <- alert:
			// Alert queued for processing
		default:
			// Channel full, skip this update
			log.Println("⚠️  Tray alert channel full, skipping update")
		}
	}
}

// openHistory opens the alerts history in a terminal window with management commands
func (app *TrayApp) openHistory() {
	var cmd *exec.Cmd

	// Create a script that shows alerts and management commands
	script := `email-sentinel.exe alerts && echo. && echo ============================================ && echo Alert Management Commands: && echo. && echo   email-sentinel alerts          - View all alerts && echo   email-sentinel alerts clear    - Clear all alerts && echo. && echo Filter Management Commands: && echo. && echo   email-sentinel filter list     - List all filters && echo   email-sentinel filter add      - Add new filter && echo   email-sentinel filter remove   - Remove a filter && echo. && echo OTP Commands: && echo. && echo   email-sentinel otp list        - View OTP codes && echo   email-sentinel otp clear       - Clear old OTPs && echo. && echo ============================================ && echo. && pause`

	switch runtime.GOOS {
	case "windows":
		// Open new cmd window and run the script
		cmd = exec.Command("cmd", "/c", "start", "cmd", "/k", script)
	case "darwin":
		// macOS - create a temporary script file
		script = `email-sentinel alerts && echo "" && echo "============================================" && echo "Alert Management Commands:" && echo "" && echo "  email-sentinel alerts          - View all alerts" && echo "  email-sentinel alerts clear    - Clear all alerts" && echo "" && echo "Filter Management Commands:" && echo "" && echo "  email-sentinel filter list     - List all filters" && echo "  email-sentinel filter add      - Add new filter" && echo "  email-sentinel filter remove   - Remove a filter" && echo "" && echo "OTP Commands:" && echo "" && echo "  email-sentinel otp list        - View OTP codes" && echo "  email-sentinel otp clear       - Clear old OTPs" && echo "" && echo "============================================" && echo "" && read -p "Press any key to continue..."`
		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to do script "`+script+`"`)
	default:
		// Linux - create a bash script
		script = `email-sentinel alerts && echo "" && echo "============================================" && echo "Alert Management Commands:" && echo "" && echo "  email-sentinel alerts          - View all alerts" && echo "  email-sentinel alerts clear    - Clear all alerts" && echo "" && echo "Filter Management Commands:" && echo "" && echo "  email-sentinel filter list     - List all filters" && echo "  email-sentinel filter add      - Add new filter" && echo "  email-sentinel filter remove   - Remove a filter" && echo "" && echo "OTP Commands:" && echo "" && echo "  email-sentinel otp list        - View OTP codes" && echo "  email-sentinel otp clear       - Clear old OTPs" && echo "" && echo "============================================" && echo "" && read -p "Press any key to continue..."`
		terminals := []string{"gnome-terminal", "konsole", "xterm"}
		for _, term := range terminals {
			if _, err := exec.LookPath(term); err == nil {
				cmd = exec.Command(term, "-e", "bash", "-c", script)
				break
			}
		}
	}

	if cmd != nil {
		if err := cmd.Start(); err != nil {
			log.Printf("Error opening history: %v", err)
		}
	}
}

// clearAlerts deletes all alerts from the database and refreshes the tray
func (app *TrayApp) clearAlerts() {
	deleted, err := storage.DeleteAllAlerts(app.db)
	if err != nil {
		log.Printf("❌ Error clearing alerts: %v", err)
		return
	}

	if deleted > 0 {
		log.Printf("🗑️  Cleared %d alert(s) from tray", deleted)
		app.scheduleRefresh()
	} else {
		log.Println("✨ No alerts to clear")
	}
}

// isValidGmailURL validates that a URL is a legitimate Gmail link
// This prevents command injection attacks via malicious email subjects
func isValidGmailURL(urlStr string) bool {
	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Must be HTTPS
	if parsedURL.Scheme != "https" {
		return false
	}

	// Must be mail.google.com domain
	if !strings.HasSuffix(parsedURL.Host, "mail.google.com") {
		return false
	}

	// Path should start with /mail/
	if !strings.HasPrefix(parsedURL.Path, "/mail/") {
		return false
	}

	return true
}

// openBrowser opens the given URL in the default browser
// URL is validated before execution to prevent command injection
func openBrowser(urlStr string) {
	// Validate URL to prevent command injection attacks
	if !isValidGmailURL(urlStr) {
		log.Printf("⚠️  Security: Blocked invalid Gmail URL: %s", urlStr)
		return
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", urlStr)
	case "darwin":
		cmd = exec.Command("open", urlStr)
	default:
		cmd = exec.Command("xdg-open", urlStr)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

// isToday checks if a time is today
func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// Quit sends a quit signal to the system tray
func Quit() {
	systray.Quit()
}
