/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/otp"
	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/datateamsix/email-sentinel/internal/ui"
)

// otpGetCmd represents the otp get command
var otpGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the most recent OTP code",
	Long: `Get the most recent active OTP code and copy it to clipboard.

This command retrieves the newest OTP code that hasn't expired yet
and automatically copies it to your clipboard for easy pasting.

Examples:
  email-sentinel otp get`,
	Run: runOTPGet,
}

func init() {
	otpCmd.AddCommand(otpGetCmd)
}

func runOTPGet(cmd *cobra.Command, args []string) {
	// Open database
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("❌ Error opening database: %v\n", err)
		fmt.Println("   Tip: Database may not exist. Start monitoring with 'email-sentinel start' first.")
		os.Exit(1)
	}
	defer storage.CloseDB(db)

	// Get most recent active OTP
	otps, err := storage.GetActiveOTPAlerts(db)
	if err != nil {
		fmt.Printf("❌ Error fetching OTP: %v\n", err)
		os.Exit(1)
	}

	if len(otps) == 0 {
		fmt.Println("📭 No active OTP codes found")
		fmt.Println("   Tip: OTP codes expire after a few minutes.")
		return
	}

	otpAlert := otps[0] // Most recent

	// Display OTP
	fmt.Printf("🔐 Most Recent OTP: %s\n", ui.ColorBold.Sprint(otpAlert.OTPCode))

	// Copy to clipboard
	err = otp.CopyToClipboard(otpAlert.OTPCode)
	if err != nil {
		fmt.Printf("⚠️  Failed to copy to clipboard: %v\n", err)
		fmt.Printf("   Code: %s\n", otpAlert.OTPCode)
	} else {
		fmt.Println("✅ Copied to clipboard!")

		// Mark as copied in database
		if err := storage.MarkOTPAsCopied(db, otpAlert.ID); err != nil {
			// Non-fatal error, just log it
			fmt.Printf("   Warning: Failed to mark as copied: %v\n", err)
		}
	}

	// Display metadata
	fmt.Printf("From: %s\n", otpAlert.Sender)
	fmt.Printf("Expires %s\n", formatExpiry(otpAlert.ExpiresAt))
}
