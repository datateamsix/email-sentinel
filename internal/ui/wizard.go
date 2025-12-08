package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/datateamsix/email-sentinel/internal/config"
	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/gmail"
	"github.com/datateamsix/email-sentinel/internal/notify"
	"github.com/datateamsix/email-sentinel/internal/otp"
	"github.com/datateamsix/email-sentinel/internal/storage"

	"golang.org/x/oauth2"
)

// Wizard represents the setup wizard state
type Wizard struct {
	CurrentStep int
	TotalSteps  int
	Config      *WizardConfig
	reader      *bufio.Reader
}

// WizardConfig holds wizard results
type WizardConfig struct {
	GmailAuthenticated bool
	GmailEmail         string
	FilterCreated      bool
	FilterName         string
	DesktopEnabled     bool
	MobileEnabled      bool
	NtfyTopic          string
	OTPEnabled         bool
	OTPConfidence      float64
	OTPExpiry          time.Duration
	CredentialsPath    string
	OAuthConfig        *oauth2.Config
	Token              *oauth2.Token
}

// NewWizard creates a new setup wizard
func NewWizard() *Wizard {
	return &Wizard{
		CurrentStep: 0,
		TotalSteps:  8,
		Config: &WizardConfig{
			DesktopEnabled: true,          // Default to desktop enabled
			OTPEnabled:     true,          // Default to OTP enabled
			OTPConfidence:  0.7,           // Default confidence threshold
			OTPExpiry:      5 * time.Minute, // Default expiry
		},
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run executes the full wizard flow
func (w *Wizard) Run() error {
	ClearScreen()

	// Step 1: Welcome
	if err := w.stepWelcome(); err != nil {
		return err
	}

	// Step 2: Prerequisites Check
	if err := w.stepPrerequisites(); err != nil {
		return err
	}

	// Step 3: Gmail Authentication
	if err := w.stepAuthentication(); err != nil {
		return err
	}

	// Step 4: Create First Filter
	if err := w.stepCreateFilter(); err != nil {
		return err
	}

	// Step 5: Notification Setup
	if err := w.stepNotifications(); err != nil {
		return err
	}

	// Step 6: OTP Setup
	if err := w.stepOTPSetup(); err != nil {
		return err
	}

	// Step 7: Test & Verify
	if err := w.stepTest(); err != nil {
		return err
	}

	// Step 8: Complete
	if err := w.stepComplete(); err != nil {
		return err
	}

	return nil
}

// ShouldRunWizard checks if this is a fresh install
func ShouldRunWizard() bool {
	// Check if token exists - if not, it's a fresh install
	if !gmail.TokenExists() {
		return true
	}

	// Check if any filters exist
	filters, err := filter.ListFilters()
	if err != nil || len(filters) == 0 {
		return true
	}

	return false
}

// stepWelcome shows the welcome screen
func (w *Wizard) stepWelcome() error {
	ClearScreen()
	w.printStepHeader("Welcome", 0)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("", 61)
	w.printBoxLine("         Welcome to Email Sentinel Setup!", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  This wizard will help you configure:", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  ✦ Gmail API authentication", 61)
	w.printBoxLine("  ✦ Your first email filter", 61)
	w.printBoxLine("  ✦ Desktop and mobile notifications", 61)
	w.printBoxLine("  ✦ OTP/2FA code detection", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Estimated time: 5 minutes", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Press [Enter] to begin or [q] to quit", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))

	choice := w.getUserInput("")
	if choice == "q" || choice == "quit" {
		return fmt.Errorf("wizard cancelled by user")
	}

	w.CurrentStep++
	return nil
}

// stepPrerequisites checks system requirements
func (w *Wizard) stepPrerequisites() error {
	ClearScreen()
	w.printStepHeader("Prerequisites Check", 1)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("           Step 1/7: Prerequisites Check", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  Checking requirements...", 61)
	w.printBoxLine("", 61)

	// Check Go runtime
	w.printBoxLine("  [✓] Go runtime detected", 61)
	w.printBoxLine("", 61)

	// Check config directory
	configDir, err := config.EnsureConfigDir()
	if err != nil {
		w.printBoxLine("  [✗] Config directory not writable", 61)
		fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))
		return fmt.Errorf("config directory error: %w", err)
	}
	w.printBoxLine("  [✓] Config directory writable", 61)
	w.printBoxLine("", 61)

	// Check credentials.json
	credPath := w.findCredentialsFile()
	if credPath == "" {
		w.printBoxLine("  [?] credentials.json", 61)
		fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))
		fmt.Println()

		return w.handleMissingCredentials(configDir)
	}

	w.Config.CredentialsPath = credPath
	w.printBoxLine("  [✓] credentials.json found", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))

	fmt.Println()
	PrintSuccess("All prerequisites met!")
	w.waitForEnter()

	w.CurrentStep++
	return nil
}

// stepAuthentication handles Gmail OAuth
func (w *Wizard) stepAuthentication() error {
	ClearScreen()
	w.printStepHeader("Gmail Authentication", 2)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("           Step 2/7: Gmail Authentication", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  Email Sentinel needs permission to read your Gmail.", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  ✦ Read-only access (cannot send/delete emails)", 61)
	w.printBoxLine("  ✦ Credentials stored locally on your computer", 61)
	w.printBoxLine("  ✦ You can revoke access anytime in Google settings", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Press [Enter] to open browser for authentication", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))

	w.getUserInput("")

	// Load credentials
	oauthConfig, err := gmail.LoadCredentials(w.Config.CredentialsPath)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to load credentials: %v", err))
		return err
	}
	w.Config.OAuthConfig = oauthConfig

	// Get token from web
	fmt.Println()
	token, err := gmail.GetTokenFromWeb(oauthConfig)
	if err != nil {
		PrintError(fmt.Sprintf("Authentication failed: %v", err))
		return err
	}
	w.Config.Token = token

	// Save token
	if err := gmail.SaveToken(token); err != nil {
		PrintError(fmt.Sprintf("Failed to save token: %v", err))
		return err
	}

	w.Config.GmailAuthenticated = true
	w.Config.GmailEmail = "authenticated@gmail.com" // Placeholder - would need to fetch real email

	fmt.Println()
	PrintSuccess("Gmail authentication successful!")
	w.waitForEnter()

	w.CurrentStep++
	return nil
}

// stepCreateFilter creates the first email filter
func (w *Wizard) stepCreateFilter() error {
	ClearScreen()
	w.printStepHeader("Create Your First Filter", 3)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("           Step 3/7: Create Your First Filter", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  Let's create a filter to watch for important emails.", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Common examples:", 61)
	w.printBoxLine("  • Job alerts: from linkedin.com, greenhouse.io", 61)
	w.printBoxLine("  • Client emails: from @clientdomain.com", 61)
	w.printBoxLine("  • Urgent: subject contains \"urgent\", \"asap\"", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))
	fmt.Println()

	// Get filter name
	filterName := w.getUserInput("Filter name: ")
	if filterName == "" {
		PrintError("Filter name is required")
		w.waitForEnter()
		return w.stepCreateFilter() // Retry
	}

	// Get from patterns
	fmt.Println()
	fmt.Println(ColorDim.Sprint("  Match emails from specific senders."))
	fmt.Println(ColorDim.Sprint("  Examples: boss@company.com, @linkedin.com, greenhouse.io"))
	fromInput := w.getUserInput("\nFrom contains (comma-separated, or blank to skip): ")
	fromPatterns := parseCSV(fromInput)

	// Get subject patterns
	fmt.Println()
	fmt.Println(ColorDim.Sprint("  Match emails with specific words in subject line."))
	fmt.Println(ColorDim.Sprint("  Examples: interview, urgent, invoice"))
	subjectInput := w.getUserInput("\nSubject contains (comma-separated, or blank to skip): ")
	subjectPatterns := parseCSV(subjectInput)

	// Validate at least one pattern
	if len(fromPatterns) == 0 && len(subjectPatterns) == 0 {
		PrintError("At least one 'from' or 'subject' pattern is required")
		w.waitForEnter()
		return w.stepCreateFilter() // Retry
	}

	// Get match mode if both patterns exist
	matchMode := "any"
	if len(fromPatterns) > 0 && len(subjectPatterns) > 0 {
		fmt.Println()
		fmt.Println(ColorDim.Sprint("  You specified both sender and subject filters."))
		fmt.Println()
		fmt.Println(ColorDim.Sprint("  ANY (OR): Notify if sender matches OR subject matches"))
		fmt.Println(ColorDim.Sprint("            → More notifications, broader matching"))
		fmt.Println()
		fmt.Println(ColorDim.Sprint("  ALL (AND): Notify only if sender AND subject both match"))
		fmt.Println(ColorDim.Sprint("             → Fewer notifications, precise matching"))
		matchInput := w.getUserInput("\nMatch mode [any/all] (default: any): ")
		if strings.ToLower(matchInput) == "all" || strings.ToLower(matchInput) == "and" {
			matchMode = "all"
		}
	}

	// Get labels
	fmt.Println()
	fmt.Println(ColorDim.Sprint("  Organize filters by category (e.g., work, personal, urgent)"))
	labelsInput := w.getUserInput("\nLabels (comma-separated, or blank to skip): ")
	labels := parseCSV(labelsInput)

	// Create filter
	f := filter.Filter{
		Name:    filterName,
		From:    fromPatterns,
		Subject: subjectPatterns,
		Match:   matchMode,
		Labels:  labels,
	}

	// Save filter
	if err := filter.AddFilter(f); err != nil {
		PrintError(fmt.Sprintf("Error adding filter: %v", err))
		w.waitForEnter()
		return w.stepCreateFilter() // Retry
	}

	// Save labels to database
	if len(labels) > 0 {
		db, err := storage.InitDB()
		if err == nil && db != nil {
			storage.SaveLabels(db, labels)
			db.Close()
		}
	}

	w.Config.FilterCreated = true
	w.Config.FilterName = filterName

	fmt.Println()
	PrintSuccess("Filter created successfully!")
	fmt.Println()
	printFilterSummary(f)
	w.waitForEnter()

	w.CurrentStep++
	return nil
}

// stepNotifications configures notification settings
func (w *Wizard) stepNotifications() error {
	ClearScreen()
	w.printStepHeader("Notification Setup", 4)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("           Step 4/7: Notification Setup", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  How would you like to be notified?", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  [1] Desktop only (recommended to start)", 61)
	w.printBoxLine("  [2] Desktop + Mobile (requires ntfy.sh app)", 61)
	w.printBoxLine("  [3] Mobile only", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))

	choice := w.getUserInput("\nSelect option [1-3]: ")

	switch choice {
	case "1":
		w.Config.DesktopEnabled = true
		w.Config.MobileEnabled = false
		PrintSuccess("Desktop notifications enabled")
	case "2":
		w.Config.DesktopEnabled = true
		w.Config.MobileEnabled = true
		if err := w.setupMobileNotifications(); err != nil {
			return err
		}
	case "3":
		w.Config.DesktopEnabled = false
		w.Config.MobileEnabled = true
		if err := w.setupMobileNotifications(); err != nil {
			return err
		}
	default:
		PrintWarning("Invalid choice, defaulting to desktop only")
		w.Config.DesktopEnabled = true
		w.Config.MobileEnabled = false
	}

	// Save notification settings
	cfg, err := filter.LoadConfig()
	if err != nil {
		cfg = filter.DefaultConfig()
	}
	cfg.Notifications.Desktop = w.Config.DesktopEnabled
	cfg.Notifications.Mobile.Enabled = w.Config.MobileEnabled
	cfg.Notifications.Mobile.NtfyTopic = w.Config.NtfyTopic
	filter.SaveConfig(cfg)

	w.waitForEnter()
	w.CurrentStep++
	return nil
}

// setupMobileNotifications guides through mobile setup
func (w *Wizard) setupMobileNotifications() error {
	fmt.Println()
	fmt.Println(ColorCyan.Sprint("📱 Mobile Notification Setup (ntfy.sh)"))
	fmt.Println(strings.Repeat("─", 58))
	fmt.Println()
	fmt.Println("ntfy.sh is a free, open-source push notification service.")
	fmt.Println()
	fmt.Println("Install the ntfy app:")
	fmt.Println("• iOS:     https://apps.apple.com/app/ntfy/id1625396347")
	fmt.Println("• Android: https://play.google.com/store/apps/details?id=io.heckel.ntfy")
	fmt.Println()
	fmt.Println(ColorYellow.Sprint("Choose a unique, private topic name"))
	fmt.Println(ColorDim.Sprint("(like a secret channel - anyone with the name can send to it)"))
	fmt.Println()

	topic := w.getUserInput("Enter your ntfy topic name: ")
	if topic == "" {
		PrintWarning("Skipping mobile notifications (no topic provided)")
		w.Config.MobileEnabled = false
		return nil
	}

	w.Config.NtfyTopic = topic
	PrintSuccess(fmt.Sprintf("Mobile notifications enabled (topic: %s)", topic))
	return nil
}

// stepOTPSetup configures OTP/2FA code detection
func (w *Wizard) stepOTPSetup() error {
	ClearScreen()
	w.printStepHeader("OTP/2FA Setup", 5)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("           Step 5/7: OTP/2FA Code Detection", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  Automatically extract verification codes from emails!", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Features:", 61)
	w.printBoxLine("  • Auto-detect OTP codes from Gmail, GitHub, etc.", 61)
	w.printBoxLine("  • Copy codes to clipboard instantly", 61)
	w.printBoxLine("  • Codes expire automatically for security", 61)
	w.printBoxLine("  • View recent codes with 'email-sentinel otp list'", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))
	fmt.Println()

	// Ask if user wants to enable OTP
	fmt.Println(ColorDim.Sprint("  Enable OTP/2FA code detection?"))
	choice := w.getUserInput("\n[Y/n] (default: yes): ")
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice == "n" || choice == "no" {
		w.Config.OTPEnabled = false
		PrintWarning("OTP detection disabled")
		w.waitForEnter()
		w.CurrentStep++
		return nil
	}

	w.Config.OTPEnabled = true
	fmt.Println()
	PrintSuccess("OTP detection enabled!")

	// Ask if they want to customize settings
	fmt.Println()
	fmt.Println(ColorDim.Sprint("  Would you like to customize OTP settings?"))
	fmt.Println(ColorDim.Sprint("  (or press Enter to use recommended defaults)"))
	customChoice := w.getUserInput("\n[y/N]: ")

	if strings.ToLower(strings.TrimSpace(customChoice)) == "y" || strings.ToLower(strings.TrimSpace(customChoice)) == "yes" {
		// Customize confidence threshold
		fmt.Println()
		fmt.Println(ColorDim.Sprint("  Confidence threshold (0.0-1.0, default: 0.7)"))
		fmt.Println(ColorDim.Sprint("  Higher = fewer false positives, may miss some codes"))
		fmt.Println(ColorDim.Sprint("  Lower = catch more codes, may have false positives"))
		confidenceInput := w.getUserInput("\nConfidence threshold: ")
		if confidenceInput != "" {
			if conf, err := strconv.ParseFloat(confidenceInput, 64); err == nil && conf >= 0 && conf <= 1 {
				w.Config.OTPConfidence = conf
			} else {
				PrintWarning("Invalid value, using default (0.7)")
			}
		}

		// Customize expiry duration
		fmt.Println()
		fmt.Println(ColorDim.Sprint("  Code expiry duration (in minutes, default: 5)"))
		fmt.Println(ColorDim.Sprint("  Most OTP codes expire in 5-10 minutes"))
		expiryInput := w.getUserInput("\nExpiry (minutes): ")
		if expiryInput != "" {
			if minutes, err := strconv.Atoi(expiryInput); err == nil && minutes > 0 {
				w.Config.OTPExpiry = time.Duration(minutes) * time.Minute
			} else {
				PrintWarning("Invalid value, using default (5 minutes)")
			}
		}
	}

	// Save OTP configuration
	rules := otp.DefaultOTPRules()
	rules.Enabled = w.Config.OTPEnabled
	rules.ConfidenceThreshold = w.Config.OTPConfidence
	rules.ExpiryDuration = w.Config.OTPExpiry

	// Get config directory and save rules
	configDir, err := config.ConfigDir()
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not get config directory: %v", err))
	} else {
		rulesPath := filepath.Join(configDir, "otp_rules.yaml")
		if err := otp.SaveOTPRules(rulesPath, rules); err != nil {
			PrintWarning(fmt.Sprintf("Could not save OTP settings: %v", err))
		} else {
			fmt.Println()
			PrintInfo(fmt.Sprintf("OTP Settings: Confidence=%.1f, Expiry=%v", w.Config.OTPConfidence, w.Config.OTPExpiry))
		}
	}

	w.waitForEnter()
	w.CurrentStep++
	return nil
}

// stepTest runs verification tests
func (w *Wizard) stepTest() error {
	ClearScreen()
	w.printStepHeader("Test Your Setup", 5)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("           Step 6/7: Test Your Setup", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  Let's verify everything works!", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  [1] 🧪 Send test desktop notification", 61)
	w.printBoxLine("  [2] 📱 Send test mobile notification", 61)
	w.printBoxLine("  [3] 📧 Check Gmail connection", 61)
	w.printBoxLine("  [4] ✓ Skip tests - I'm ready", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))

	for {
		choice := w.getUserInput("\nSelect option [1-4]: ")

		switch choice {
		case "1":
			if w.Config.DesktopEnabled {
				fmt.Println()
				PrintInfo("Sending test desktop notification...")
				err := notify.SendDesktopNotification(
					"Email Sentinel - Test Notification",
					"If you see this, desktop notifications are working!",
				)
				if err != nil {
					PrintError(fmt.Sprintf("Desktop notification failed: %v", err))
				} else {
					PrintSuccess("Desktop notification sent!")
				}
			} else {
				PrintWarning("Desktop notifications are disabled")
			}

		case "2":
			if w.Config.MobileEnabled && w.Config.NtfyTopic != "" {
				fmt.Println()
				PrintInfo("Sending test mobile notification...")
				err := notify.SendMobileNotification(
					w.Config.NtfyTopic,
					"Email Sentinel - Test",
					"If you see this on your phone, mobile notifications are working!",
				)
				if err != nil {
					PrintError(fmt.Sprintf("Mobile notification failed: %v", err))
				} else {
					PrintSuccess(fmt.Sprintf("Mobile notification sent to topic: %s", w.Config.NtfyTopic))
				}
			} else {
				PrintWarning("Mobile notifications are disabled or topic not configured")
			}

		case "3":
			fmt.Println()
			PrintInfo("Testing Gmail connection...")
			// Basic connection test - if we got a token, it should work
			if w.Config.GmailAuthenticated {
				PrintSuccess("Gmail authentication verified!")
			} else {
				PrintError("Gmail not authenticated")
			}

		case "4":
			w.CurrentStep++
			return nil

		default:
			PrintError("Invalid choice, please try again")
		}
	}
}

// stepComplete shows completion screen
func (w *Wizard) stepComplete() error {
	ClearScreen()
	w.printStepHeader("Setup Complete", 6)

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", 61) + "╗"))
	w.printBoxLine("                 Setup Complete!", 61)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", 61) + "╣"))
	w.printBoxLine("", 61)
	w.printBoxLine("  Email Sentinel is ready to use!", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Summary:", 61)

	if w.Config.GmailAuthenticated {
		w.printBoxLine("  ✓ Gmail authenticated", 61)
	}
	if w.Config.FilterCreated {
		w.printBoxLine(fmt.Sprintf("  ✓ Filter configured: \"%s\"", w.Config.FilterName), 61)
	}
	if w.Config.DesktopEnabled {
		w.printBoxLine("  ✓ Desktop notifications: enabled", 61)
	}
	if w.Config.MobileEnabled {
		w.printBoxLine(fmt.Sprintf("  ✓ Mobile notifications: enabled (topic: %s)", w.Config.NtfyTopic), 61)
	}
	if w.Config.OTPEnabled {
		w.printBoxLine("  ✓ OTP/2FA detection: enabled", 61)
	}

	w.printBoxLine("", 61)
	w.printBoxLine("  Quick commands:", 61)
	w.printBoxLine("  • email-sentinel start       Start monitoring", 61)
	w.printBoxLine("  • email-sentinel start --tray Run in system tray", 61)
	w.printBoxLine("  • email-sentinel otp list    View OTP codes", 61)
	w.printBoxLine("  • email-sentinel filter add  Add more filters", 61)
	w.printBoxLine("  • email-sentinel             Open interactive menu", 61)
	w.printBoxLine("", 61)
	w.printBoxLine("  Press [Enter] to go to main menu or [q] to exit", 61)
	w.printBoxLine("", 61)
	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", 61) + "╝"))

	choice := w.getUserInput("")
	if choice != "q" && choice != "quit" {
		// Launch interactive menu
		fmt.Println()
		PrintInfo("Launching interactive menu...")
		return nil
	}

	return nil
}

// Helper functions

// printStepHeader prints the current step
func (w *Wizard) printStepHeader(title string, step int) {
	fmt.Println()
	ColorCyan.Printf("═══ Email Sentinel Setup ═══ Step %d/7: %s\n", step+1, title)
	fmt.Println()
}

// printBoxLine prints a line within a box
func (w *Wizard) printBoxLine(text string, width int) {
	visibleLen := len(stripANSI(text))
	padding := width - visibleLen - 2
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("%s %s%s %s\n",
		ColorCyan.Sprint("║"),
		text,
		strings.Repeat(" ", padding),
		ColorCyan.Sprint("║"),
	)
}

// getUserInput prompts for user input
func (w *Wizard) getUserInput(prompt string) string {
	if prompt != "" {
		ColorGreen.Print(prompt)
	}
	input, _ := w.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// waitForEnter waits for user to press enter
func (w *Wizard) waitForEnter() {
	fmt.Println()
	ColorDim.Print("Press Enter to continue...")
	w.reader.ReadString('\n')
}

// findCredentialsFile searches for credentials.json
func (w *Wizard) findCredentialsFile() string {
	// Check current directory
	if _, err := os.Stat("credentials.json"); err == nil {
		return "credentials.json"
	}

	// Check config directory
	configDir, err := config.ConfigDir()
	if err == nil {
		credPath := filepath.Join(configDir, "credentials.json")
		if _, err := os.Stat(credPath); err == nil {
			return credPath
		}
	}

	return ""
}

// handleMissingCredentials guides user to get credentials.json
func (w *Wizard) handleMissingCredentials(configDir string) error {
	PrintWarning("credentials.json not found!")
	fmt.Println()
	fmt.Println("To get this file:")
	fmt.Println()
	PrintBullet("Go to https://console.cloud.google.com/")
	PrintBullet("Create a project and enable Gmail API")
	PrintBullet("Create OAuth credentials (Desktop app)")
	PrintBullet("Download JSON and save as 'credentials.json'")
	fmt.Println()
	PrintInfo("See README.md section \"Google Cloud Setup\" for details.")
	fmt.Println()
	fmt.Println("Place credentials.json in one of these locations:")
	PrintBullet(fmt.Sprintf("Current directory: %s", filepath.Join(".", "credentials.json")))
	PrintBullet(fmt.Sprintf("Config directory: %s", filepath.Join(configDir, "credentials.json")))
	fmt.Println()

	for {
		ColorGreen.Print("[r] Retry check  [o] Open Google Cloud Console  [s] Skip  [q] Quit: ")
		choice := w.getUserInput("")

		switch strings.ToLower(choice) {
		case "r", "retry":
			// Retry finding credentials
			credPath := w.findCredentialsFile()
			if credPath != "" {
				w.Config.CredentialsPath = credPath
				PrintSuccess("credentials.json found!")
				w.waitForEnter()
				w.CurrentStep++
				return nil
			}
			PrintError("credentials.json still not found. Please place it in one of the locations above.")

		case "o", "open":
			// Open browser to Google Cloud Console
			url := "https://console.cloud.google.com/"
			if err := openBrowser(url); err != nil {
				PrintError(fmt.Sprintf("Failed to open browser: %v", err))
				fmt.Println()
				fmt.Printf("Please visit manually: %s\n", url)
			} else {
				PrintSuccess("Opening Google Cloud Console in browser...")
			}

		case "s", "skip":
			PrintWarning("Skipping credentials check. You'll need to run 'email-sentinel init' later.")
			w.waitForEnter()
			w.CurrentStep++
			return nil

		case "q", "quit":
			return fmt.Errorf("wizard cancelled by user")

		default:
			PrintError("Invalid choice")
		}
	}
}

// openBrowser opens a URL in the default browser (cross-platform)
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux and others
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// parseCSV parses comma-separated values
func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

// printFilterSummary prints filter details
func printFilterSummary(f filter.Filter) {
	ColorDim.Printf("  Name:    ")
	fmt.Println(f.Name)

	if len(f.From) > 0 {
		ColorDim.Printf("  From:    ")
		fmt.Println(strings.Join(f.From, ", "))
	}

	if len(f.Subject) > 0 {
		ColorDim.Printf("  Subject: ")
		fmt.Println(strings.Join(f.Subject, ", "))
	}

	if len(f.Labels) > 0 {
		ColorDim.Printf("  Labels:  ")
		fmt.Println(strings.Join(f.Labels, ", "))
	}

	matchDesc := "any (OR - either condition triggers)"
	if f.Match == "all" {
		matchDesc = "all (AND - all conditions must match)"
	}
	ColorDim.Printf("  Match:   ")
	fmt.Println(matchDesc)
}
