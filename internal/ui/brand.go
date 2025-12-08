package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// Constants
const (
	AppName    = "Email Sentinel"
	AppTagline = "Real-time Gmail Notification System"

	// Width constraints
	BannerWidth = 58
	DividerChar = "═"
)

// Version information (injected at build time via ldflags)
var (
	AppVersion = "dev"      // Overridden by -ldflags "-X ..." at build time
	BuildTime  = "unknown"  // Build timestamp
	GitCommit  = "none"     // Git commit SHA
)

// Color definitions (gracefully degrade if terminal doesn't support)
var (
	// Primary colors
	ColorCyan    = color.New(color.FgCyan)
	ColorBlue    = color.New(color.FgBlue)
	ColorGreen   = color.New(color.FgGreen)
	ColorYellow  = color.New(color.FgYellow)
	ColorRed     = color.New(color.FgRed)
	ColorGray    = color.New(color.FgHiBlack)

	// Styles
	ColorBold    = color.New(color.Bold)
	ColorDim     = color.New(color.Faint)

	// Composite styles
	HeaderStyle  = color.New(color.FgCyan, color.Bold)
	SuccessStyle = color.New(color.FgGreen, color.Bold)
	ErrorStyle   = color.New(color.FgRed, color.Bold)
	WarningStyle = color.New(color.FgYellow, color.Bold)
	InfoStyle    = color.New(color.FgBlue, color.Bold)
)

// Symbols
const (
	SymbolCheck   = "✓"
	SymbolCross   = "✗"
	SymbolWarning = "⚠"
	SymbolInfo    = "ℹ"
	SymbolArrow   = "→"
	SymbolBullet  = "•"
)

// Banner design (exactly 58 chars wide)
const bannerTemplate = `
╔════════════════════════════════════════════════════════╗
║                                                        ║
║     ███████╗███╗   ███╗ █████╗ ██╗██╗                 ║
║     ██╔════╝████╗ ████║██╔══██╗██║██║                 ║
║     █████╗  ██╔████╔██║███████║██║██║                 ║
║     ██╔══╝  ██║╚██╔╝██║██╔══██║██║██║                 ║
║     ███████╗██║ ╚═╝ ██║██║  ██║██║███████╗            ║
║                                                        ║
║       ███████╗███████╗███╗   ██╗████████╗██╗          ║
║       ██╔════╝██╔════╝████╗  ██║╚══██╔══╝██║          ║
║       ███████╗█████╗  ██╔██╗ ██║   ██║   ██║          ║
║       ╚════██║██╔══╝  ██║╚██╗██║   ██║   ██║          ║
║       ███████║███████╗██║ ╚████║   ██║   ██║          ║
║                                                        ║
║           📧  Real-time Gmail Monitoring               ║
║                  Version %-18s            ║
╚════════════════════════════════════════════════════════╝
`

// Compact banner for limited space
const compactBannerTemplate = `
┌──────────────────────────────────────────────────────┐
│  📧 Email Sentinel  •  v%-18s           │
│      Real-time Gmail Notification System             │
└──────────────────────────────────────────────────────┘
`

// PrintBanner displays the branded ASCII banner
func PrintBanner(version string) {
	if version == "" {
		version = AppVersion
	}
	banner := fmt.Sprintf(bannerTemplate, version)
	ColorCyan.Print(banner)
}

// PrintCompactBanner displays a smaller version of the banner
func PrintCompactBanner(version string) {
	if version == "" {
		version = AppVersion
	}
	banner := fmt.Sprintf(compactBannerTemplate, version)
	ColorCyan.Print(banner)
}

// PrintSimpleBanner displays a minimal text banner
func PrintSimpleBanner(version string) {
	if version == "" {
		version = AppVersion
	}
	ColorBold.Printf("📧 %s", AppName)
	fmt.Printf(" v%s\n", version)
}

// PrintSection displays a section header
// Example: ════════════ FILTER MANAGEMENT ════════════
func PrintSection(title string) {
	title = strings.ToUpper(title)
	titleLen := len(title)

	// Calculate padding for centering (total width 58)
	totalPadding := BannerWidth - titleLen - 2 // -2 for spaces around title
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	fmt.Println()
	HeaderStyle.Printf("%s %s %s\n",
		strings.Repeat(DividerChar, leftPadding),
		title,
		strings.Repeat(DividerChar, rightPadding),
	)
	fmt.Println()
}

// PrintSubsection displays a subsection header
func PrintSubsection(title string) {
	fmt.Println()
	ColorCyan.Printf("▸ %s\n", title)
}

// PrintSuccess displays a success message with checkmark
func PrintSuccess(message string) {
	SuccessStyle.Printf("%s ", SymbolCheck)
	fmt.Println(message)
}

// PrintError displays an error message with X
func PrintError(message string) {
	ErrorStyle.Printf("%s ", SymbolCross)
	fmt.Println(message)
}

// PrintWarning displays a warning message
func PrintWarning(message string) {
	WarningStyle.Printf("%s ", SymbolWarning)
	fmt.Println(message)
}

// PrintInfo displays an info message
func PrintInfo(message string) {
	InfoStyle.Printf("%s ", SymbolInfo)
	fmt.Println(message)
}

// PrintDivider prints a horizontal line
func PrintDivider() {
	ColorGray.Println(strings.Repeat("─", BannerWidth))
}

// PrintBullet prints a bulleted list item
func PrintBullet(text string) {
	ColorCyan.Printf("  %s ", SymbolBullet)
	fmt.Println(text)
}

// PrintKeyValue prints a key-value pair
func PrintKeyValue(key, value string) {
	ColorDim.Printf("  %-20s ", key+":")
	fmt.Println(value)
}

// PrintTable prints a simple table
func PrintTable(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print header
	fmt.Print("  ")
	for i, header := range headers {
		ColorBold.Printf("%-*s  ", colWidths[i], header)
	}
	fmt.Println()

	// Print separator
	fmt.Print("  ")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("─", width) + "  ")
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		fmt.Print("  ")
		for i, cell := range row {
			if i < len(colWidths) {
				fmt.Printf("%-*s  ", colWidths[i], cell)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// PrintBox draws a box around lines of text
func PrintBox(lines []string) {
	if len(lines) == 0 {
		return
	}

	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	fmt.Printf("  ┌─%s─┐\n", strings.Repeat("─", maxLen))
	for _, line := range lines {
		padding := strings.Repeat(" ", maxLen-len(line))
		fmt.Printf("  │ %s%s │\n", line, padding)
	}
	fmt.Printf("  └─%s─┘\n", strings.Repeat("─", maxLen))
}

// PrintCommandExample prints a command example
func PrintCommandExample(description, command string) {
	ColorGray.Printf("  # %s\n", description)
	ColorGreen.Print("  $ ")
	fmt.Println(command)
}

// PrintVersionInfo displays detailed version information
func PrintVersionInfo(version, buildTime, gitCommit string) {
	if version == "" {
		version = AppVersion
	}

	fmt.Println()
	ColorBold.Printf("%s\n", AppName)
	PrintKeyValue("Version", version)
	PrintKeyValue("Build Time", buildTime)
	PrintKeyValue("Git Commit", gitCommit)
	PrintKeyValue("Go Version", runtime.Version())
	PrintKeyValue("Platform", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	fmt.Println()
}

// PrintWelcome displays a welcome message for first-time users
func PrintWelcome() {
	fmt.Println()
	PrintDivider()
	fmt.Println()

	ColorBold.Printf("  Welcome to %s!\n\n", AppName)

	fmt.Println("  Your personal Gmail monitoring assistant. Get instant")
	fmt.Println("  notifications when important emails arrive.")
	fmt.Println()

	ColorGreen.Println("  Quick Start:")
	PrintBullet("Initialize:     email-sentinel init")
	PrintBullet("Add filter:     email-sentinel filter add")
	PrintBullet("Start watching: email-sentinel start --tray")
	fmt.Println()

	ColorYellow.Println("  Need help?")
	PrintBullet("Run: email-sentinel --help")
	PrintBullet("Docs: https://github.com/datateamsix/email-sentinel")
	fmt.Println()

	PrintDivider()
	fmt.Println()
}

// ClearScreen clears the terminal (cross-platform)
func ClearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

// DisableColors disables all color output (for CI/CD or when piping)
func DisableColors() {
	color.NoColor = true
}

// EnableColors enables color output
func EnableColors() {
	color.NoColor = false
}

// IsColorEnabled returns whether colors are currently enabled
func IsColorEnabled() bool {
	return !color.NoColor
}

// Colorize returns text with color (if colors enabled)
func Colorize(text string, c *color.Color) string {
	return c.Sprint(text)
}

// Bold returns bold text
func Bold(text string) string {
	return ColorBold.Sprint(text)
}

// Dim returns dimmed text
func Dim(text string) string {
	return ColorDim.Sprint(text)
}

// Highlight returns highlighted text (cyan)
func Highlight(text string) string {
	return ColorCyan.Sprint(text)
}

// Success returns success-styled text (green)
func Success(text string) string {
	return SuccessStyle.Sprint(text)
}

// Error returns error-styled text (red)
func Error(text string) string {
	return ErrorStyle.Sprint(text)
}

// Warning returns warning-styled text (yellow)
func Warning(text string) string {
	return WarningStyle.Sprint(text)
}

// Info returns info-styled text (blue)
func Info(text string) string {
	return InfoStyle.Sprint(text)
}

// ProgressBar creates a simple progress bar
func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)
	empty := width - filled

	bar := fmt.Sprintf("[%s%s] %d/%d",
		ColorGreen.Sprint(strings.Repeat("█", filled)),
		ColorGray.Sprint(strings.Repeat("░", empty)),
		current,
		total,
	)

	return bar
}

// PrintProgressBar prints a progress bar
func PrintProgressBar(current, total int, width int) {
	fmt.Println(ProgressBar(current, total, width))
}

// Spinner characters for loading animation
var spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// GetSpinnerChar returns a spinner character based on step
func GetSpinnerChar(step int) string {
	return spinnerChars[step%len(spinnerChars)]
}
