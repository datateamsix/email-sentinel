package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks user for yes/no confirmation
// Returns true if user confirms (y/yes), false otherwise
func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(ColorYellow.Sprint("? "))
	fmt.Print(message)
	fmt.Print(" ")
	ColorDim.Print("[y/N]: ")

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes"
}

// ConfirmYes asks for confirmation with default Yes
// Returns true if user confirms (y/yes) or presses Enter, false otherwise
func ConfirmYes(message string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(ColorYellow.Sprint("? "))
	fmt.Print(message)
	fmt.Print(" ")
	ColorDim.Print("[Y/n]: ")

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	// Default to yes if Enter is pressed
	return response == "" || response == "y" || response == "yes"
}

// ConfirmDangerous asks for confirmation with warning styling
// Used for destructive operations (delete, reset, etc.)
func ConfirmDangerous(message string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	PrintWarning(message)
	fmt.Println()

	fmt.Print(ColorRed.Sprint("⚠ "))
	ColorBold.Print("This action cannot be undone. Are you sure? ")
	ColorDim.Print("[y/N]: ")

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes"
}

// ConfirmWithOptions asks user to choose from multiple options
// Returns the selected option (0-indexed) or -1 if cancelled
func ConfirmWithOptions(message string, options []string) int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Print(ColorYellow.Sprint("? "))
	fmt.Println(message)
	fmt.Println()

	for i, option := range options {
		fmt.Printf("  [%d] %s\n", i+1, option)
	}
	fmt.Println()
	fmt.Print(ColorGreen.Sprintf("Select option [1-%d] or [q] to cancel: ", len(options)))

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "q" || response == "quit" || response == "cancel" {
		return -1
	}

	// Try to parse as number
	var choice int
	_, err := fmt.Sscanf(response, "%d", &choice)
	if err != nil || choice < 1 || choice > len(options) {
		return -1
	}

	return choice - 1 // Return 0-indexed
}

// AskInput prompts for text input with optional default value
func AskInput(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(ColorGreen.Sprint("› "))
	fmt.Print(prompt)

	if defaultValue != "" {
		ColorDim.Printf(" [%s]", defaultValue)
	}
	fmt.Print(": ")

	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "" && defaultValue != "" {
		return defaultValue
	}

	return response
}

// AskPassword prompts for password input (masked)
// Note: This is a simple implementation. For production, consider using
// a library like golang.org/x/term for proper password input
func AskPassword(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(ColorGreen.Sprint("› "))
	fmt.Print(prompt)
	fmt.Print(": ")

	// Simple implementation - shows characters
	// For production, use terminal package to hide input
	password, _ := reader.ReadString('\n')
	return strings.TrimSpace(password)
}

// PressEnterToContinue waits for user to press Enter
func PressEnterToContinue() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()
	ColorDim.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

// PressAnyKeyToContinue waits for any key press (actually waits for Enter)
func PressAnyKeyToContinue() {
	PressEnterToContinue()
}
