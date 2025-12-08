/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
)

// otpClearCmd represents the otp clear command
var otpClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear expired OTP codes",
	Long: `Delete all expired OTP codes from the database.

This helps keep the database clean by removing old codes that can
no longer be used. You'll be prompted for confirmation.

Examples:
  email-sentinel otp clear`,
	Run: runOTPClear,
}

func init() {
	otpCmd.AddCommand(otpClearCmd)
}

func runOTPClear(cmd *cobra.Command, args []string) {
	// Open database
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("❌ Error opening database: %v\n", err)
		fmt.Println("   Tip: Database may not exist. Start monitoring with 'email-sentinel start' first.")
		os.Exit(1)
	}
	defer storage.CloseDB(db)

	// First expire codes, then count them
	storage.ExpireOTPAlerts(db)

	// Count expired (inactive) codes
	otps, err := storage.GetRecentOTPAlerts(db, 1000)
	if err != nil {
		fmt.Printf("❌ Error counting codes: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, otp := range otps {
		if !otp.IsActive {
			count++
		}
	}

	if count == 0 {
		fmt.Println("✨ No expired OTP codes to clear")
		return
	}

	// Prompt for confirmation
	fmt.Printf("Found %d expired OTP code(s)\n", count)
	fmt.Print("Delete these codes? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ Error reading input: %v\n", err)
		os.Exit(1)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Cancelled.")
		return
	}

	// Delete expired codes
	deleted, err := storage.DeleteExpiredOTPAlerts(db)
	if err != nil {
		fmt.Printf("❌ Error deleting codes: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔐 Cleared %d expired OTP code(s)\n", deleted)
}
