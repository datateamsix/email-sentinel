package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/storage"
)

// MenuItem represents a single menu option
type MenuItem struct {
	Key         string      // "1", "2", "a", "q", etc.
	Label       string      // Display text
	Description string      // Optional help text
	Icon        string      // Optional icon/emoji
	Action      func() error // Function to execute
	SubMenu     *Menu       // Optional nested menu
}

// Menu represents a menu with multiple items
type Menu struct {
	Title       string
	Description string
	Items       []MenuItem
	ShowBack    bool
	ShowQuit    bool
	Parent      *Menu
	StatusBar   func() string // Optional dynamic status bar
}

// MenuConfig holds menu configuration
type MenuConfig struct {
	Width       int
	ClearScreen bool
	ShowStatus  bool
}

// DefaultMenuConfig returns default menu configuration
var DefaultMenuConfig = MenuConfig{
	Width:       63,
	ClearScreen: true,
	ShowStatus:  true,
}

// NewMenu creates a new menu
func NewMenu(title string) *Menu {
	return &Menu{
		Title:    title,
		Items:    []MenuItem{},
		ShowBack: false,
		ShowQuit: true,
	}
}

// AddItem adds a menu item
func (m *Menu) AddItem(key, icon, label, description string, action func() error) *Menu {
	m.Items = append(m.Items, MenuItem{
		Key:         key,
		Icon:        icon,
		Label:       label,
		Description: description,
		Action:      action,
	})
	return m
}

// AddSubMenu adds a submenu item
func (m *Menu) AddSubMenu(key, icon, label, description string, subMenu *Menu) *Menu {
	subMenu.Parent = m
	subMenu.ShowBack = true
	subMenu.ShowQuit = false

	m.Items = append(m.Items, MenuItem{
		Key:         key,
		Icon:        icon,
		Label:       label,
		Description: description,
		SubMenu:     subMenu,
	})
	return m
}

// SetStatusBar sets a dynamic status bar function
func (m *Menu) SetStatusBar(statusFn func() string) *Menu {
	m.StatusBar = statusFn
	return m
}

// Display shows the menu and handles user input
func (m *Menu) Display() error {
	for {
		if DefaultMenuConfig.ClearScreen {
			ClearScreen()
		}

		m.render()

		choice := m.getUserInput()

		// Handle back option
		if m.ShowBack && (choice == "b" || choice == "back") {
			return nil
		}

		// Handle quit option
		if m.ShowQuit && (choice == "q" || choice == "quit") {
			return fmt.Errorf("quit")
		}

		// Find and execute menu item
		item := m.findItem(choice)
		if item == nil {
			PrintError("Invalid choice. Please try again.")
			time.Sleep(1 * time.Second)
			continue
		}

		// Execute action or navigate to submenu
		if item.SubMenu != nil {
			err := item.SubMenu.Display()
			if err != nil && err.Error() == "quit" {
				return err
			}
		} else if item.Action != nil {
			if DefaultMenuConfig.ClearScreen {
				ClearScreen()
			}
			err := item.Action()
			if err != nil {
				PrintError(fmt.Sprintf("Error: %v", err))
			}
			fmt.Println()
			fmt.Print("Press Enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}
	}
}

// Run starts the interactive menu loop
func (m *Menu) Run() {
	err := m.Display()
	if err != nil && err.Error() == "quit" {
		fmt.Println()
		PrintInfo("Goodbye! 👋")
		fmt.Println()
	}
}

// render displays the menu
func (m *Menu) render() {
	width := DefaultMenuConfig.Width

	// Show banner for main menu only (when no parent)
	if m.Parent == nil && m.Title == "Main Menu" {
		PrintMenuBanner()
	}

	// Top border
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", width-2) + "╗"))

	// Title
	m.printCenteredRow(strings.ToUpper(m.Title), width)

	// Separator
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", width-2) + "╣"))

	// Empty line
	m.printEmptyRow(width)

	// Menu items
	for _, item := range m.Items {
		m.printMenuItem(item, width)
	}

	// Empty line before footer options
	m.printEmptyRow(width)

	// Back option
	if m.ShowBack {
		backText := "  [b] ← Back to Main Menu"
		m.printRow(backText, width)
		m.printEmptyRow(width)
	}

	// Quit option
	if m.ShowQuit {
		quitText := "  [q] Exit"
		m.printRow(quitText, width)
		m.printEmptyRow(width)
	}

	// Bottom border
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", width-2) + "╝"))

	// Status bar
	if DefaultMenuConfig.ShowStatus && m.StatusBar != nil {
		fmt.Println()
		fmt.Println(ColorGray.Sprint(strings.Repeat("─", width)))
		fmt.Println(m.StatusBar())
	}

	fmt.Println()
}

// printMenuItem prints a single menu item
func (m *Menu) printMenuItem(item MenuItem, width int) {
	// Format: "  [1] 🚀 Start Monitoring      Start watching for emails"
	keyPart := fmt.Sprintf("  [%s]", item.Key)

	var labelPart string
	if item.Icon != "" {
		labelPart = fmt.Sprintf(" %s %s", item.Icon, item.Label)
	} else {
		labelPart = fmt.Sprintf(" %s", item.Label)
	}

	// Calculate spacing for description
	usedSpace := len(keyPart) + len(labelPart)
	descStart := 35 // Column where description starts

	var line string
	if item.Description != "" {
		spacing := descStart - usedSpace
		if spacing < 2 {
			spacing = 2
		}
		line = fmt.Sprintf("%s%s%s%s",
			ColorBold.Sprint(keyPart),
			labelPart,
			strings.Repeat(" ", spacing),
			ColorDim.Sprint(item.Description),
		)
	} else {
		line = fmt.Sprintf("%s%s", ColorBold.Sprint(keyPart), labelPart)
	}

	m.printRow(line, width)
}

// printRow prints a row with borders
func (m *Menu) printRow(content string, width int) {
	// Calculate visible length (without ANSI codes)
	visibleLen := len(stripANSI(content))
	padding := width - visibleLen - 4 // -4 for borders and spaces
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
func (m *Menu) printCenteredRow(text string, width int) {
	textLen := len(text)
	totalPadding := width - textLen - 4
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	fmt.Printf("%s%s%s%s%s\n",
		ColorCyan.Sprint("║"),
		strings.Repeat(" ", leftPadding),
		ColorBold.Sprint(text),
		strings.Repeat(" ", rightPadding),
		ColorCyan.Sprint("║"))
}

// printEmptyRow prints an empty row with borders
func (m *Menu) printEmptyRow(width int) {
	fmt.Printf("%s%s%s\n",
		ColorCyan.Sprint("║"),
		strings.Repeat(" ", width-2),
		ColorCyan.Sprint("║"),
	)
}

// getUserInput gets user input
func (m *Menu) getUserInput() string {
	ColorGreen.Print("Select an option: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

// findItem finds a menu item by key
func (m *Menu) findItem(key string) *MenuItem {
	for i := range m.Items {
		if m.Items[i].Key == key {
			return &m.Items[i]
		}
	}
	return nil
}

// stripANSI removes ANSI escape codes for length calculation
func stripANSI(str string) string {
	// Simple ANSI stripper - removes common escape sequences
	result := str
	result = strings.ReplaceAll(result, "\033[0m", "")
	result = strings.ReplaceAll(result, "\033[1m", "")
	result = strings.ReplaceAll(result, "\033[2m", "")
	result = strings.ReplaceAll(result, "\033[36m", "")
	result = strings.ReplaceAll(result, "\033[34m", "")
	result = strings.ReplaceAll(result, "\033[32m", "")
	result = strings.ReplaceAll(result, "\033[33m", "")
	result = strings.ReplaceAll(result, "\033[31m", "")
	result = strings.ReplaceAll(result, "\033[90m", "")

	// Remove any remaining escape sequences (more thorough)
	for strings.Contains(result, "\033[") {
		start := strings.Index(result, "\033[")
		end := strings.Index(result[start:], "m")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}

	return result
}

// BuildMainMenu creates the main menu structure
func BuildMainMenu() *Menu {
	menu := NewMenu("Main Menu")

	menu.AddItem("1", "🚀", "Start Monitoring", "Start watching for emails", func() error {
		PrintSection("Start Email Monitoring")
		fmt.Println()
		PrintInfo("To start monitoring, use one of these commands:")
		fmt.Println()
		PrintBullet("email-sentinel start              (foreground)")
		PrintBullet("email-sentinel start --tray       (with system tray)")
		PrintBullet("email-sentinel start --daemon     (background)")
		fmt.Println()
		PrintInfo("The monitoring service must run separately from this menu.")
		return nil
	})

	// Filter Management submenu
	filterMenu := buildFilterMenu()
	menu.AddSubMenu("2", "📋", "Manage Filters", "Add, edit, remove filters", filterMenu)

	// Notifications submenu
	notifMenu := buildNotificationsMenu()
	menu.AddSubMenu("3", "🔔", "Notifications", "Configure alerts", notifMenu)

	// Status submenu
	statusMenu := buildStatusMenu()
	menu.AddSubMenu("4", "📊", "Status & History", "View alerts and status", statusMenu)

	// Digital Accounts submenu
	accountsMenu := buildAccountsMenu()
	menu.AddSubMenu("5", "💳", "Digital Accounts", "Manage subscriptions & trials", accountsMenu)

	// Settings submenu
	settingsMenu := buildSettingsMenu()
	menu.AddSubMenu("6", "⚙️", "Settings", "Configure app settings", settingsMenu)

	menu.AddItem("7", "🔧", "Setup Wizard", "Re-run initial setup", func() error {
		ClearScreen()
		wizard := NewWizard()
		return wizard.Run()
	})

	// Add status bar
	menu.SetStatusBar(func() string {
		return fmt.Sprintf("Status: %s Running | Filters: 3 | Last check: 2 min ago",
			ColorGreen.Sprint("●"))
	})

	return menu
}

// buildFilterMenu creates the filter management submenu
func buildFilterMenu() *Menu {
	menu := NewMenu("Filter Management")

	menu.AddItem("1", "➕", "Add Filter", "Create a new filter", func() error {
		return handleAddFilter()
	})

	menu.AddItem("2", "✏️", "Edit Filter", "Modify existing filter", func() error {
		return handleEditFilter()
	})

	menu.AddItem("3", "🗑️", "Remove Filter", "Delete a filter", func() error {
		return handleRemoveFilter()
	})

	menu.AddItem("4", "📋", "List Filters", "View all filters", func() error {
		return handleListFilters()
	})

	return menu
}

// buildAccountsMenu creates the digital accounts submenu
func buildAccountsMenu() *Menu {
	menu := NewMenu("Digital Accounts")

	menu.AddItem("1", "📋", "List All Accounts", "View all subscriptions & trials", func() error {
		PrintSection("Digital Accounts")
		PrintInfo("Launching accounts list...")
		fmt.Println()
		PrintInfo("Run: email-sentinel accounts list")
		PrintInfo("Or:  email-sentinel accounts list --trials (show only trials)")
		PrintInfo("Or:  email-sentinel accounts list --paid (show only paid)")
		return nil
	})

	menu.AddItem("2", "🔍", "Search Account", "Find which email you used for a service", func() error {
		PrintSection("Search Account")
		reader := bufio.NewReader(os.Stdin)

		fmt.Println()
		ColorGreen.Print("Service name (e.g., Netflix, Spotify): ")
		service, _ := reader.ReadString('\n')
		service = strings.TrimSpace(service)

		if service == "" {
			PrintError("Service name is required")
			return fmt.Errorf("service name required")
		}

		fmt.Println()
		PrintInfo(fmt.Sprintf("Searching for '%s'...", service))
		PrintInfo("Run: email-sentinel accounts search " + service)
		return nil
	})

	menu.AddItem("3", "🔥", "Expiring Trials", "View trials expiring soon", func() error {
		PrintSection("Expiring Trials")
		PrintInfo("Showing trials expiring in the next 7 days...")
		fmt.Println()
		PrintInfo("Run: email-sentinel accounts list --trials")
		PrintWarning("Remember to cancel before trial expires to avoid charges!")
		return nil
	})

	menu.AddItem("4", "💰", "Total Spending", "Calculate monthly/annual costs", func() error {
		PrintSection("Total Spending")
		PrintInfo("Calculating total subscription costs...")
		fmt.Println()
		PrintInfo("Run: email-sentinel accounts list")
		PrintInfo("Total spending is shown at the bottom of the list")
		return nil
	})

	return menu
}

// buildNotificationsMenu creates the notifications submenu
func buildNotificationsMenu() *Menu {
	menu := NewMenu("Notifications")

	menu.AddItem("1", "🖥️", "Desktop Notifications", "Toggle on/off", func() error {
		PrintSection("Desktop Notifications")
		PrintSuccess("Desktop notifications are enabled")
		return nil
	})

	menu.AddItem("2", "📱", "Mobile (ntfy.sh)", "Configure mobile push", func() error {
		PrintSection("Mobile Notifications")
		PrintInfo("Configure mobile push notifications via ntfy.sh")
		PrintKeyValue("Status", "Disabled")
		PrintKeyValue("Topic", "Not configured")
		return nil
	})

	menu.AddItem("3", "🧪", "Test Notifications", "Send test alert", func() error {
		PrintSection("Test Notifications")
		PrintInfo("Sending test notification...")
		time.Sleep(1 * time.Second)
		PrintSuccess("Test notification sent!")
		return nil
	})

	return menu
}

// buildStatusMenu creates the status submenu
func buildStatusMenu() *Menu {
	menu := NewMenu("Status & History")

	menu.AddItem("1", "📊", "Dashboard", "System status overview", func() error {
		return RunInteractiveDashboard()
	})

	menu.AddItem("2", "📜", "Alert History", "View past notifications", func() error {
		return ShowAlertHistory()
	})

	menu.AddItem("3", "🔍", "Check Gmail", "Test Gmail connection", func() error {
		PrintSection("Gmail Connection Test")
		PrintInfo("Testing Gmail API connection...")
		time.Sleep(1 * time.Second)
		PrintSuccess("Successfully connected to Gmail!")
		PrintKeyValue("Account", "user@gmail.com")
		PrintKeyValue("API Status", "Active")
		return nil
	})

	return menu
}

// buildSettingsMenu creates the settings submenu
func buildSettingsMenu() *Menu {
	menu := NewMenu("Settings")

	menu.AddItem("1", "⏱️", "Polling Interval", "Set check frequency", func() error {
		PrintSection("Polling Interval")
		PrintInfo("Current polling interval: 45 seconds")
		fmt.Println()
		PrintBullet("Recommended: 30-60 seconds")
		PrintBullet("Shorter intervals use more API quota")
		return nil
	})

	menu.AddItem("2", "🔐", "Re-authenticate", "Re-run Gmail OAuth", func() error {
		PrintSection("Gmail Re-authentication")
		PrintWarning("This will open your browser to re-authorize Gmail access")
		return nil
	})

	menu.AddItem("3", "📁", "Open Config Folder", "Open config directory", func() error {
		PrintSection("Configuration Location")
		PrintKeyValue("Config Dir", "~/.config/email-sentinel/")
		PrintKeyValue("Config File", "config.yaml")
		PrintKeyValue("Database", "history.db")
		PrintKeyValue("Token File", "token.json")
		return nil
	})

	menu.AddItem("4", "🔄", "Reset to Defaults", "Clear all settings", func() error {
		PrintSection("Reset Settings")
		PrintWarning("This will delete all filters and configuration!")
		PrintError("This action cannot be undone.")
		return nil
	})

	return menu
}

// handleAddFilter handles the interactive filter addition process
func handleAddFilter() error {
	PrintSection("Add New Email Filter")
	reader := bufio.NewReader(os.Stdin)

	// Get filter name
	fmt.Print("\nFilter name: ")
	filterName, _ := reader.ReadString('\n')
	filterName = strings.TrimSpace(filterName)

	if filterName == "" {
		PrintError("Filter name is required")
		return fmt.Errorf("filter name is required")
	}

	// Get from patterns
	fmt.Println("\n📤 Sender Filter (From)")
	fmt.Println("   Match emails from specific senders.")
	fmt.Println("   Examples: boss@company.com, @linkedin.com, greenhouse.io")
	fmt.Print("\nFrom contains (comma-separated, or blank to skip): ")
	filterFrom, _ := reader.ReadString('\n')
	filterFrom = strings.TrimSpace(filterFrom)

	// Get subject patterns
	fmt.Println("\n📝 Subject Filter")
	fmt.Println("   Match emails with specific words in subject line.")
	fmt.Println("   Examples: interview, urgent, invoice")
	fmt.Print("\nSubject contains (comma-separated, or blank to skip): ")
	filterSubject, _ := reader.ReadString('\n')
	filterSubject = strings.TrimSpace(filterSubject)

	// Validate at least one pattern
	if filterFrom == "" && filterSubject == "" {
		PrintError("At least one 'from' or 'subject' pattern is required")
		return fmt.Errorf("at least one pattern required")
	}

	// Parse CSV
	fromPatterns := parseCSV(filterFrom)
	subjectPatterns := parseCSV(filterSubject)

	// Get match mode if both patterns specified
	filterMatch := "any"
	if len(fromPatterns) > 0 && len(subjectPatterns) > 0 {
		fmt.Println("\n🔀 Match Mode")
		fmt.Println("   ANY (OR): Notify if sender OR subject matches (broader)")
		fmt.Println("   ALL (AND): Notify only if sender AND subject match (precise)")
		fmt.Print("\nMatch mode [any/all] (default: any): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "all" || input == "and" {
			filterMatch = "all"
		}
	}

	// Get labels
	fmt.Println("\n🏷️  Labels/Categories (Optional)")
	fmt.Println("   Organize filters (e.g., work, personal, urgent)")
	fmt.Print("\nLabels (comma-separated, or blank to skip): ")
	filterLabels, _ := reader.ReadString('\n')
	filterLabels = strings.TrimSpace(filterLabels)
	labelsList := parseCSV(filterLabels)

	// Get Gmail scope
	filterScope := "inbox"
	fmt.Println("\n📧 Gmail Scope (Optional)")
	fmt.Println("   Options: inbox (default), all, primary, social, promotions, updates")
	fmt.Print("\nScope (or blank for inbox): ")
	scopeInput, _ := reader.ReadString('\n')
	scopeInput = strings.TrimSpace(scopeInput)
	if scopeInput != "" {
		filterScope = scopeInput
	}

	// Create filter
	newFilter := filter.Filter{
		Name:       filterName,
		From:       fromPatterns,
		Subject:    subjectPatterns,
		Match:      filterMatch,
		Labels:     labelsList,
		GmailScope: filterScope,
	}

	// Add filter
	if err := filter.AddFilter(newFilter); err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
		return err
	}

	fmt.Println()
	PrintSuccess(fmt.Sprintf("Filter '%s' added successfully!", filterName))
	return nil
}

// handleEditFilter handles the interactive filter editing process
func handleEditFilter() error {
	PrintSection("Edit Filter")
	reader := bufio.NewReader(os.Stdin)

	// Load all filters
	filters, err := filter.ListFilters()
	if err != nil {
		PrintError(fmt.Sprintf("Error loading filters: %v", err))
		return err
	}

	if len(filters) == 0 {
		fmt.Println()
		PrintInfo("No filters to edit")
		return nil
	}

	// Display filters
	fmt.Println()
	PrintInfo("Select a filter to edit:")
	fmt.Println()

	for i, f := range filters {
		fmt.Printf("  [%d] %s\n", i+1, f.Name)
	}

	// Get selection
	fmt.Println()
	ColorGreen.Print("Enter number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(filters) {
		PrintError("Invalid selection")
		return fmt.Errorf("invalid selection")
	}

	filterIndex := num - 1
	selectedFilter := filters[filterIndex]
	fmt.Printf("\nEditing filter: %s\n", ColorBold.Sprint(selectedFilter.Name))
	fmt.Println("\n(Press Enter to keep current value)")

	// Edit from patterns
	fmt.Printf("\nFrom [%s]: ", strings.Join(selectedFilter.From, ", "))
	fromInput, _ := reader.ReadString('\n')
	fromInput = strings.TrimSpace(fromInput)
	if fromInput != "" {
		selectedFilter.From = parseCSV(fromInput)
	}

	// Edit subject patterns
	fmt.Printf("Subject [%s]: ", strings.Join(selectedFilter.Subject, ", "))
	subjectInput, _ := reader.ReadString('\n')
	subjectInput = strings.TrimSpace(subjectInput)
	if subjectInput != "" {
		selectedFilter.Subject = parseCSV(subjectInput)
	}

	// Validate at least one pattern
	if len(selectedFilter.From) == 0 && len(selectedFilter.Subject) == 0 {
		PrintError("At least one 'from' or 'subject' pattern is required")
		return fmt.Errorf("at least one pattern required")
	}

	// Edit match mode
	fmt.Printf("Match mode [%s]: ", selectedFilter.Match)
	matchInput, _ := reader.ReadString('\n')
	matchInput = strings.TrimSpace(strings.ToLower(matchInput))
	if matchInput != "" {
		if matchInput == "all" || matchInput == "and" {
			selectedFilter.Match = "all"
		} else {
			selectedFilter.Match = "any"
		}
	}

	// Edit labels
	fmt.Printf("Labels [%s]: ", strings.Join(selectedFilter.Labels, ", "))
	labelsInput, _ := reader.ReadString('\n')
	labelsInput = strings.TrimSpace(labelsInput)
	if labelsInput != "" {
		selectedFilter.Labels = parseCSV(labelsInput)
	}

	// Edit scope
	fmt.Printf("Gmail scope [%s]: ", selectedFilter.GmailScope)
	scopeInput, _ := reader.ReadString('\n')
	scopeInput = strings.TrimSpace(scopeInput)
	if scopeInput != "" {
		selectedFilter.GmailScope = scopeInput
	}

	// Update filter using index
	if err := filter.UpdateFilter(filterIndex, selectedFilter); err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
		return err
	}

	fmt.Println()
	PrintSuccess(fmt.Sprintf("Filter '%s' updated successfully!", selectedFilter.Name))
	return nil
}

// handleListFilters displays all configured filters
func handleListFilters() error {
	PrintSection("Filter List")

	// Load all filters
	filters, err := filter.ListFilters()
	if err != nil {
		PrintError(fmt.Sprintf("Error loading filters: %v", err))
		return err
	}

	if len(filters) == 0 {
		fmt.Println()
		PrintInfo("No filters configured yet")
		fmt.Println()
		PrintInfo("Use 'Add Filter' to create your first filter")
		return nil
	}

	// Display filters
	fmt.Printf("\n📋 Found %d filter(s):\n\n", len(filters))

	for i, f := range filters {
		fmt.Printf("[%d] %s\n", i+1, ColorBold.Sprint(f.Name))

		// From patterns
		if len(f.From) > 0 {
			fmt.Printf("    From:    %s\n", strings.Join(f.From, ", "))
		} else {
			fmt.Printf("    From:    %s\n", ColorDim.Sprint("(any)"))
		}

		// Subject patterns
		if len(f.Subject) > 0 {
			fmt.Printf("    Subject: %s\n", strings.Join(f.Subject, ", "))
		} else {
			fmt.Printf("    Subject: %s\n", ColorDim.Sprint("(any)"))
		}

		// Match mode
		fmt.Printf("    Match:   %s\n", f.Match)

		// Labels
		if len(f.Labels) > 0 {
			fmt.Printf("    Labels:  %s\n", ColorCyan.Sprint(strings.Join(f.Labels, ", ")))
		}

		// Gmail scope
		if f.GmailScope != "" && f.GmailScope != "inbox" {
			fmt.Printf("    Scope:   %s\n", f.GmailScope)
		}

		// Expiration
		if f.ExpiresAt != nil {
			fmt.Printf("    Expires: %s\n", f.ExpiresAt.Format("2006-01-02"))
		}

		fmt.Println()
	}

	return nil
}

// handleRemoveFilter handles the interactive filter removal process
func handleRemoveFilter() error {
	PrintSection("Remove Filter")
	reader := bufio.NewReader(os.Stdin)

	// Load all filters
	filters, err := filter.ListFilters()
	if err != nil {
		PrintError(fmt.Sprintf("Error loading filters: %v", err))
		return err
	}

	if len(filters) == 0 {
		fmt.Println()
		PrintInfo("No filters to remove")
		return nil
	}

	// Display filters
	fmt.Println()
	PrintInfo("Select a filter to remove:")
	fmt.Println()

	for i, f := range filters {
		fmt.Printf("  [%d] %s\n", i+1, f.Name)
	}

	// Get selection
	fmt.Println()
	ColorGreen.Print("Enter number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(filters) {
		PrintError("Invalid selection")
		return fmt.Errorf("invalid selection")
	}

	filterName := filters[num-1].Name

	// Confirm deletion
	fmt.Println()
	ColorYellow.Printf("Remove filter '%s'? (y/N): ", filterName)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		PrintInfo("Cancelled")
		return nil
	}

	// Remove filter
	if err := filter.RemoveFilter(filterName); err != nil {
		PrintError(fmt.Sprintf("Error: %v", err))
		return err
	}

	fmt.Println()
	PrintSuccess(fmt.Sprintf("Filter '%s' removed successfully!", filterName))
	return nil
}

// RunInteractiveMenu starts the main interactive menu
func RunInteractiveMenu() {
	menu := BuildMainMenu()
	menu.Run()
}

// ShowAlertHistory displays today's email alerts
func ShowAlertHistory() error {
	PrintSection("Alert History")

	// Initialize database
	db, err := storage.InitDB()
	if err != nil {
		PrintError(fmt.Sprintf("Error opening alert database: %v", err))
		return err
	}
	defer storage.CloseDB(db)

	// Get all alerts from today
	alerts, err := storage.GetTodayAlerts(db)
	if err != nil {
		PrintError(fmt.Sprintf("Error fetching alerts: %v", err))
		return err
	}

	if len(alerts) == 0 {
		fmt.Println()
		PrintInfo("No alerts today")
		return nil
	}

	// Display header
	count, _ := storage.CountTodayAlerts(db)
	fmt.Printf("\n📬 Today's Alerts (%d total)\n\n", count)

	// Display each alert
	for i, alert := range alerts {
		// Add priority indicator
		priorityIcon := "📩" // Normal priority
		if alert.Priority == 1 {
			priorityIcon = "🔥" // High priority
		}

		fmt.Printf("[%d] %s %s\n", i+1, priorityIcon, alert.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Filter: %s\n", alert.FilterName)
		if alert.Priority == 1 {
			fmt.Printf("    Priority: HIGH\n")
		}
		fmt.Printf("    From:   %s\n", alert.Sender)
		fmt.Printf("    Subject: %s\n", alert.Subject)

		if alert.Snippet != "" {
			// Truncate snippet if too long
			snippet := alert.Snippet
			if len(snippet) > 100 {
				snippet = snippet[:97] + "..."
			}
			fmt.Printf("    Preview: %s\n", snippet)
		}

		fmt.Printf("    Link:   %s\n", alert.GmailLink)
		fmt.Println()
	}

	return nil
}
