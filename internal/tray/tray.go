package tray

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/getlantern/systray"
)

// TrayApp represents the system tray application
type TrayApp struct {
	db              *sql.DB
	recentAlerts    []*systray.MenuItem
	alertUpdateChan chan storage.Alert
	quitChan        chan struct{}
	mu              sync.Mutex
	hasUrgent       bool
}

// Config holds configuration for the tray app
type Config struct {
	DB *sql.DB
}

var (
	globalApp     *TrayApp
	mRecentAlerts *systray.MenuItem
	mOpenHistory  *systray.MenuItem
	mQuit         *systray.MenuItem
)

// Run starts the system tray application
// This function blocks until the tray is quit
func Run(cfg Config) {
	globalApp = &TrayApp{
		db:              cfg.DB,
		alertUpdateChan: make(chan storage.Alert, 100),
		quitChan:        make(chan struct{}),
		recentAlerts:    make([]*systray.MenuItem, 0),
	}

	systray.Run(onReady, onExit)
}

// onReady is called when the system tray is ready
func onReady() {
	// Set initial icon and title
	systray.SetIcon(GetNormalIcon())
	systray.SetTitle("Email Sentinel")
	systray.SetTooltip("Email Sentinel - Monitoring Gmail")

	// Create menu items
	mRecentAlerts = systray.AddMenuItem("Recent Alerts", "View recent email alerts")
	systray.AddSeparator()
	mOpenHistory = systray.AddMenuItem("Open History", "View all alerts in terminal")
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

	// Update icon based on urgent status
	if hasUrgent {
		systray.SetIcon(GetUrgentIcon())
		systray.SetTooltip("Email Sentinel - ⚠️ Urgent alerts!")
	} else {
		systray.SetIcon(GetNormalIcon())
		systray.SetTooltip("Email Sentinel - Monitoring Gmail")
	}

	// Add each alert as a submenu item
	for _, alert := range alerts {
		app.addAlertMenuItem(alert)
	}
}

// addAlertMenuItem adds a single alert to the recent alerts submenu
func (app *TrayApp) addAlertMenuItem(alert storage.Alert) {
	// Format the menu item title
	icon := "📧"
	if alert.Priority == 1 {
		icon = "🔥"
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
	tooltip := fmt.Sprintf("From: %s\nClick to open in Gmail", alert.Sender)

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

	for {
		select {
		case alert := <-app.alertUpdateChan:
			log.Printf("📱 Tray: New alert received - %s", alert.Subject)

			// Temporarily switch to urgent icon if it's a priority alert
			if alert.Priority == 1 {
				app.hasUrgent = true
				systray.SetIcon(GetUrgentIcon())
				systray.SetTooltip("Email Sentinel - ⚠️ New urgent alert!")

				// Flash the icon by switching back after a few seconds
				go func() {
					time.Sleep(5 * time.Second)
					// Check if there are still urgent alerts
					app.loadRecentAlerts()
				}()
			}

			// Reload the alerts menu
			app.loadRecentAlerts()

		case <-ticker.C:
			// Periodically refresh the alerts
			app.loadRecentAlerts()

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

// openHistory opens the alerts history in a terminal window
func (app *TrayApp) openHistory() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Open new cmd window and run the alerts command
		cmd = exec.Command("cmd", "/c", "start", "cmd", "/k", "email-sentinel.exe alerts")
	case "darwin":
		// macOS - open new Terminal window
		cmd = exec.Command("open", "-a", "Terminal", "email-sentinel alerts")
	default:
		// Linux - try common terminal emulators
		terminals := []string{"gnome-terminal", "konsole", "xterm"}
		for _, term := range terminals {
			if _, err := exec.LookPath(term); err == nil {
				cmd = exec.Command(term, "-e", "email-sentinel alerts")
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

// openBrowser opens the given URL in the default browser
func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
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
