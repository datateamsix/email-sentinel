package ui

import (
	"fmt"
	"strings"
)

// PrintKeyboardShortcuts shows available keyboard shortcuts
func PrintKeyboardShortcuts() {
	width := 58

	fmt.Println()
	fmt.Println(ColorCyan.Sprint("╔" + strings.Repeat("═", width-2) + "╗"))
	printHelpCenteredRow("KEYBOARD SHORTCUTS", width)
	fmt.Println(ColorCyan.Sprint("╠" + strings.Repeat("═", width-2) + "╣"))
	printHelpEmptyRow(width)

	// Menu navigation
	printHelpRow("  1-9           Select menu item by number", width)
	printHelpRow("  Enter         Confirm selection", width)
	printHelpRow("  b             Back to previous menu", width)
	printHelpRow("  q             Quit application", width)
	printHelpEmptyRow(width)

	// Dashboard specific
	printHelpRow("  r             Refresh (in dashboard)", width)
	printHelpRow("  s             Start/Stop service (where available)", width)
	printHelpEmptyRow(width)

	// General
	printHelpRow("  ?             Show this help", width)
	printHelpRow("  Ctrl+C        Force quit", width)
	printHelpEmptyRow(width)

	fmt.Println(ColorCyan.Sprint("╚" + strings.Repeat("═", width-2) + "╝"))
	fmt.Println()
}

// PrintQuickHelp shows a compact help message
func PrintQuickHelp() {
	fmt.Println()
	ColorDim.Println("  Quick Help:")
	ColorDim.Println("    [1-9] Select item  [b] Back  [q] Quit  [?] Full help")
	fmt.Println()
}

// printHelpRow prints a row with borders for help display
func printHelpRow(content string, width int) {
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

// printHelpCenteredRow prints centered text with borders for help display
func printHelpCenteredRow(text string, width int) {
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

// printHelpEmptyRow prints an empty row with borders for help display
func printHelpEmptyRow(width int) {
	fmt.Printf("%s%s%s\n",
		ColorCyan.Sprint("║"),
		strings.Repeat(" ", width-2),
		ColorCyan.Sprint("║"),
	)
}
