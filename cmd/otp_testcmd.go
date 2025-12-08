/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/otp"
	"github.com/datateamsix/email-sentinel/internal/ui"
)

// otpTestCmd represents the otp test command
var otpTestCmd = &cobra.Command{
	Use:   "test <text>",
	Short: "Test OTP extraction on sample text",
	Long: `Test the OTP extraction algorithm on sample text.

This is useful for testing whether your OTP codes will be properly
extracted from emails. Pass the email subject or body text as an argument.

Examples:
  email-sentinel otp test "Your verification code is 123456"
  email-sentinel otp test "OTP: 456789"
  email-sentinel otp test "Your code 1234"`,
	Args: cobra.ExactArgs(1),
	Run:  runOTPTest,
}

func init() {
	otpCmd.AddCommand(otpTestCmd)
}

func runOTPTest(cmd *cobra.Command, args []string) {
	text := args[0]

	fmt.Printf("Testing text: %s\n\n", ui.ColorDim.Sprint(text))

	// Use default rules for testing
	rules := otp.DefaultOTPRules()

	// Extract OTP - pass text as both subject and body to maximize detection
	result := otp.DetectOTP(text, text, text, "test@example.com", rules)

	if result == nil {
		fmt.Println("❌ No OTP code detected")
		fmt.Println("\nTip: OTP extraction looks for patterns like:")
		fmt.Println("  - 'code: 123456'")
		fmt.Println("  - 'your verification code is 123456'")
		fmt.Println("  - 'OTP: 123456'")
		os.Exit(1)
	}

	// Display result
	fmt.Printf("✅ OTP Detected: %s\n", ui.ColorBold.Sprint(result.Code))
	fmt.Printf("Confidence: %.2f\n", result.Confidence)
	fmt.Printf("Pattern: %s\n", result.Pattern)
	fmt.Println("Source: text")

	if result.Confidence < 0.8 {
		fmt.Println("\n⚠️  Low confidence - code may be incorrectly extracted")
	}
}
