package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/gmail"
	"github.com/datateamsix/email-sentinel/internal/storage"
)

// Dashboard displays system status
type Dashboard struct {
	RefreshInterval time.Duration
	reader          *bufio.Reader
}

// DashboardData holds all status information
type DashboardData struct {
	// Service
	IsRunning   bool
	PID         int
	Uptime      time.Duration
	LastCheck   time.Time
	NextCheck   time.Time
	LastRun     time.Time
	HasStateInfo bool

	// Gmail
	Email       string
	AuthValid   bool
	TokenExpiry time.Time
	TokenExists bool

	// Filters
	FilterCount int
	Filters     []FilterSummary

	// Notifications
	DesktopEnabled bool
	MobileEnabled  bool
	NtfyTopic      string

	// Stats (Last 24h)
	EmailsChecked     int64
	FiltersMatched    int64
	NotificationsSent int64
	PollingInterval   int
}

// FilterSummary represents a brief filter overview
type FilterSummary struct {
	Name    string
	Summary string // Brief description
}

// NewDashboard creates a dashboard
func NewDashboard() *Dashboard {
	return &Dashboard{
		RefreshInterval: 5 * time.Second,
		reader:          bufio.NewReader(os.Stdin),
	}
}

// Display shows the dashboard once
func (d *Dashboard) Display() error {
	data, err := GatherDashboardData()
	if err != nil {
		return fmt.Errorf("failed to gather dashboard data: %w", err)
	}

	d.render(data)
	return nil
}

// Run shows dashboard with interactive controls
func (d *Dashboard) Run() error {
	for {
		ClearScreen()

		// Gather and display data
		data, err := GatherDashboardData()
		if err != nil {
			PrintError(fmt.Sprintf("Error gathering data: %v", err))
			d.waitForInput()
			return err
		}

		d.render(data)

		// Show controls
		fmt.Println()
		ColorGreen.Print("Select option: ")
		choice := d.getUserInput()

		switch strings.ToLower(choice) {
		case "r", "refresh":
			// Loop will refresh automatically
			continue
		case "s", "start", "stop":
			PrintInfo("Start/Stop functionality requires the 'email-sentinel start' command")
			d.waitForInput()
		case "b", "back":
			return nil
		default:
			PrintError("Invalid choice")
			time.Sleep(1 * time.Second)
		}
	}
}

// render displays the dashboard
func (d *Dashboard) render(data *DashboardData) {
	width := 63

	// Header
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", width-2) + "╗"))
	d.printCenteredRow("EMAIL SENTINEL STATUS", width)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", width-2) + "╣"))
	d.printEmptyRow(width)

	// Service Status
	d.printSectionTitle("Service Status", width)
	d.printDivider(width)

	if data.IsRunning {
		statusLine := fmt.Sprintf("  Watcher:     %s Running", ColorGreen.Sprint("●"))
		if data.PID > 0 {
			statusLine += fmt.Sprintf(" (PID: %d)", data.PID)
		}
		d.printRow(statusLine, width)

		if data.Uptime > 0 {
			d.printRow(fmt.Sprintf("  Uptime:      %s", formatDuration(data.Uptime)), width)
		}

		if !data.LastCheck.IsZero() {
			d.printRow(fmt.Sprintf("  Last Check:  %s", formatRelativeTime(data.LastCheck)), width)
		}

		if !data.NextCheck.IsZero() {
			d.printRow(fmt.Sprintf("  Next Check:  in %s", formatRelativeTime(data.NextCheck)), width)
		}
	} else {
		d.printRow(fmt.Sprintf("  Watcher:     %s Stopped", ColorGray.Sprint("○")), width)
		if !data.LastRun.IsZero() {
			d.printRow(fmt.Sprintf("  Last Run:    %s", formatRelativeTime(data.LastRun)), width)
		}
	}
	d.printEmptyRow(width)

	// Gmail Connection
	d.printSectionTitle("Gmail Connection", width)
	d.printDivider(width)

	if data.TokenExists {
		if data.Email != "" {
			d.printRow(fmt.Sprintf("  Account:     %s", data.Email), width)
		}

		if data.AuthValid {
			d.printRow(fmt.Sprintf("  Auth Status: %s Valid", ColorGreen.Sprint("✓")), width)
		} else {
			d.printRow(fmt.Sprintf("  Auth Status: %s Invalid/Expired", ColorRed.Sprint("✗")), width)
		}

		if !data.TokenExpiry.IsZero() && data.TokenExpiry.After(time.Now()) {
			timeUntilExpiry := time.Until(data.TokenExpiry)
			d.printRow(fmt.Sprintf("  Token Expiry: in %s", formatDuration(timeUntilExpiry)), width)
		}
	} else {
		d.printRow(fmt.Sprintf("  Auth Status: %s Not configured", ColorRed.Sprint("✗")), width)
		d.printRow("  Run: email-sentinel init", width)
	}
	d.printEmptyRow(width)

	// Filters
	d.printSectionTitle("Filters", width)
	d.printDivider(width)

	if data.FilterCount > 0 {
		d.printRow(fmt.Sprintf("  Active Filters: %d", data.FilterCount), width)
		d.printRow("  ┌─────────────────────────────────────────────────────┐", width)

		// Show up to 5 filters
		displayCount := data.FilterCount
		if displayCount > 5 {
			displayCount = 5
		}

		for i := 0; i < displayCount; i++ {
			filterLine := fmt.Sprintf("  │ %d. %-20s %s", i+1, data.Filters[i].Name, data.Filters[i].Summary)
			// Truncate if too long
			if len(stripANSI(filterLine)) > width-8 {
				filterLine = filterLine[:width-11] + "..."
			}
			d.printRow(filterLine, width)
		}

		if data.FilterCount > 5 {
			d.printRow(fmt.Sprintf("  │ ... and %d more", data.FilterCount-5), width)
		}

		d.printRow("  └─────────────────────────────────────────────────────┘", width)
	} else {
		d.printRow("  Active Filters: 0", width)
		d.printRow("  Run: email-sentinel filter add", width)
	}
	d.printEmptyRow(width)

	// Notifications
	d.printSectionTitle("Notifications", width)
	d.printDivider(width)

	if data.DesktopEnabled {
		d.printRow(fmt.Sprintf("  Desktop:     %s Enabled", ColorGreen.Sprint("✓")), width)
	} else {
		d.printRow(fmt.Sprintf("  Desktop:     %s Disabled", ColorGray.Sprint("✗")), width)
	}

	if data.MobileEnabled && data.NtfyTopic != "" {
		d.printRow(fmt.Sprintf("  Mobile:      %s Enabled (topic: %s)", ColorGreen.Sprint("✓"), data.NtfyTopic), width)
	} else {
		d.printRow(fmt.Sprintf("  Mobile:      %s Disabled", ColorGray.Sprint("✗")), width)
	}
	d.printEmptyRow(width)

	// Statistics
	d.printSectionTitle("Statistics (Last 24h)", width)
	d.printDivider(width)

	// Note: These are estimates based on polling interval and alerts
	if data.EmailsChecked > 0 {
		d.printRow(fmt.Sprintf("  Emails Checked:   %s", formatNumber(data.EmailsChecked)), width)
	} else {
		d.printRow("  Emails Checked:   N/A (not running)", width)
	}

	d.printRow(fmt.Sprintf("  Filters Matched:  %d", data.FiltersMatched), width)
	d.printRow(fmt.Sprintf("  Notifications:    %d", data.NotificationsSent), width)
	d.printEmptyRow(width)

	// Footer
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", width-2) + "╣"))
	d.printRow("  [r] Refresh   [s] Start/Stop   [b] Back", width)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", width-2) + "╝"))
}

// GatherDashboardData collects all status information
func GatherDashboardData() (*DashboardData, error) {
	data := &DashboardData{
		HasStateInfo: false,
	}

	// Check Gmail authentication
	data.TokenExists = gmail.TokenExists()
	data.AuthValid = data.TokenExists

	// Try to load token to check expiry
	if data.TokenExists {
		token, err := gmail.LoadToken()
		if err == nil && token != nil {
			data.TokenExpiry = token.Expiry
			data.AuthValid = token.Valid()

			// Try to extract email from token (if available)
			// This would require calling the Gmail API, so we'll skip for now
			data.Email = "user@gmail.com" // Placeholder
		}
	}

	// Load filters
	cfg, err := filter.LoadConfig()
	if err != nil {
		// If config doesn't exist, use defaults
		cfg = filter.DefaultConfig()
	}

	data.FilterCount = len(cfg.Filters)
	data.Filters = make([]FilterSummary, 0, data.FilterCount)

	for _, f := range cfg.Filters {
		summary := ""
		if len(f.From) > 0 {
			fromList := strings.Join(f.From, ", ")
			if len(fromList) > 30 {
				fromList = fromList[:27] + "..."
			}
			summary = fmt.Sprintf("from: %s", fromList)
		}
		if len(f.Subject) > 0 {
			subjList := strings.Join(f.Subject, ", ")
			if len(subjList) > 25 {
				subjList = subjList[:22] + "..."
			}
			if summary != "" {
				summary += " | "
			}
			summary += fmt.Sprintf("subj: %s", subjList)
		}

		data.Filters = append(data.Filters, FilterSummary{
			Name:    f.Name,
			Summary: summary,
		})
	}

	// Notification settings
	data.DesktopEnabled = cfg.Notifications.Desktop
	data.MobileEnabled = cfg.Notifications.Mobile.Enabled
	data.NtfyTopic = cfg.Notifications.Mobile.NtfyTopic
	data.PollingInterval = cfg.PollingInterval

	// Statistics from database
	db, err := storage.InitDB()
	if err == nil && db != nil {
		defer storage.CloseDB(db)

		// Count today's alerts
		count, err := storage.CountTodayAlerts(db)
		if err == nil {
			data.FiltersMatched = int64(count)
			data.NotificationsSent = int64(count) // Each alert = 1+ notifications
		}

		// Estimate emails checked (rough calculation)
		// If we have alerts from today, estimate based on polling interval
		if count > 0 {
			// Assume 24 hours of monitoring, with polling interval
			checksPerHour := 3600 / int64(cfg.PollingInterval)
			totalChecks := checksPerHour * 24
			emailsPerCheck := int64(10) // Default messages fetched per check
			data.EmailsChecked = totalChecks * emailsPerCheck
		}
	}

	// Service status - since we don't track PID/uptime, mark as not running
	// This would require implementing a daemon/PID file tracking system
	data.IsRunning = false

	return data, nil
}

// Helper functions

// printSectionTitle prints a section header
func (d *Dashboard) printSectionTitle(title string, width int) {
	d.printRow(fmt.Sprintf("  %s", ColorBold.Sprint(title)), width)
}

// printDivider prints a divider line
func (d *Dashboard) printDivider(width int) {
	d.printRow(fmt.Sprintf("  %s", strings.Repeat("─", width-6)), width)
}

// printRow prints a row with borders
func (d *Dashboard) printRow(content string, width int) {
	visibleLen := len(stripANSI(content))
	padding := width - visibleLen - 4
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("%s %s%s %s\n",
		ColorCyan.Sprint("║"),
		content,
		strings.Repeat(" ", padding),
		ColorCyan.Sprint("║"),
	)
}

// printCenteredRow prints centered text with borders
func (d *Dashboard) printCenteredRow(text string, width int) {
	textLen := len(text)
	totalPadding := width - textLen - 4
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	fmt.Printf("%s%s%s%s%s%s\n",
		ColorCyan.Sprint("║"),
		strings.Repeat(" ", leftPadding),
		ColorBold.Sprint(text),
		strings.Repeat(" ", rightPadding),
		ColorCyan.Sprint("║"),
	)
}

// printEmptyRow prints an empty row with borders
func (d *Dashboard) printEmptyRow(width int) {
	fmt.Printf("%s%s%s\n",
		ColorCyan.Sprint("║"),
		strings.Repeat(" ", width-2),
		ColorCyan.Sprint("║"),
	)
}

// getUserInput gets user input
func (d *Dashboard) getUserInput() string {
	input, _ := d.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// waitForInput waits for user to press enter
func (d *Dashboard) waitForInput() {
	fmt.Println()
	ColorDim.Print("Press Enter to continue...")
	d.reader.ReadString('\n')
}

// formatDuration formats a duration in human-readable format
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

// formatRelativeTime formats a time relative to now
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	var diff time.Duration

	if t.After(now) {
		// Future time
		diff = t.Sub(now)
		return formatDuration(diff)
	}

	// Past time
	diff = now.Sub(t)

	if diff < time.Minute {
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := int(diff.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}

// formatNumber formats a number with thousands separators
func formatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	s := fmt.Sprintf("%d", n)
	result := ""
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}

// ShowDashboard is a convenience function to display the dashboard once
func ShowDashboard() error {
	dashboard := NewDashboard()
	return dashboard.Display()
}

// RunInteractiveDashboard starts the interactive dashboard
func RunInteractiveDashboard() error {
	dashboard := NewDashboard()
	return dashboard.Run()
}
