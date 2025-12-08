/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
)

var (
	activeOnly bool
	limitOTP   int
)

// otpListCmd represents the otp list command
var otpListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent OTP codes",
	Long: `List recent OTP codes extracted from emails.

By default, shows the 10 most recent OTP codes. Use --active to show
only codes that haven't expired yet.

Examples:
  # List last 10 OTP codes
  email-sentinel otp list

  # List only active (non-expired) codes
  email-sentinel otp list --active

  # List last 20 codes
  email-sentinel otp list --limit 20`,
	Run: runOTPList,
}

func init() {
	otpCmd.AddCommand(otpListCmd)
	otpListCmd.Flags().BoolVarP(&activeOnly, "active", "a", false, "Show only active (non-expired) codes")
	otpListCmd.Flags().IntVarP(&limitOTP, "limit", "l", 10, "Maximum number of codes to show")
}

func runOTPList(cmd *cobra.Command, args []string) {
	// Open database
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("❌ Error opening database: %v\n", err)
		fmt.Println("   Tip: Database may not exist. Start monitoring with 'email-sentinel start' first.")
		os.Exit(1)
	}
	defer storage.CloseDB(db)

	// Fetch OTP codes
	var otps []storage.OTPAlert
	if activeOnly {
		otps, err = storage.GetActiveOTPAlerts(db)
		// Limit results
		if len(otps) > limitOTP {
			otps = otps[:limitOTP]
		}
	} else {
		otps, err = storage.GetRecentOTPAlerts(db, limitOTP)
	}

	if err != nil {
		fmt.Printf("❌ Error fetching OTP codes: %v\n", err)
		os.Exit(1)
	}

	if len(otps) == 0 {
		if activeOnly {
			fmt.Println("📭 No active OTP codes found")
		} else {
			fmt.Println("📭 No OTP codes found")
		}
		fmt.Println("   Tip: OTP codes are automatically extracted from matching emails.")
		return
	}

	// Display header
	if activeOnly {
		fmt.Printf("🔐 Active OTP Codes (%d)\n\n", len(otps))
	} else {
		fmt.Printf("🔐 Recent OTP Codes (%d)\n\n", len(otps))
	}

	// Display each OTP
	for i, otp := range otps {
		fmt.Println(formatOTPAlert(otp, i+1))
	}
}
