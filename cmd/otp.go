/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/datateamsix/email-sentinel/internal/ui"
)

// otpCmd represents the otp command
var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Manage OTP codes from emails",
	Long: `Manage one-time password (OTP) codes extracted from emails.

Email Sentinel can extract OTP codes from incoming emails and store them
for quick access. This is useful for two-factor authentication codes.

Available Commands:
  list    List recent OTP codes
  get     Get the most recent OTP and copy to clipboard
  clear   Clear expired OTP codes
  test    Test OTP extraction on sample text

Examples:
  email-sentinel otp list
  email-sentinel otp get
  email-sentinel otp clear`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(otpCmd)
}

// formatTimestamp formats a timestamp as relative time ("2 minutes ago")
func formatTimestamp(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds == 1 {
			return "1 second ago"
		}
		return fmt.Sprintf("%d seconds ago", seconds)
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

// formatExpiry formats expiry time as "in 3 minutes" or "expired"
func formatExpiry(t time.Time) string {
	if time.Now().After(t) {
		return ui.ColorRed.Sprint("expired")
	}

	duration := time.Until(t)

	if duration < time.Minute {
		seconds := int(duration.Seconds())
		if seconds == 1 {
			return "in 1 second"
		}
		return fmt.Sprintf("in %d seconds", seconds)
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", minutes)
	} else {
		hours := int(duration.Hours())
		if hours == 1 {
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	}
}

// formatOTPAlert formats an OTP alert for display
func formatOTPAlert(otp storage.OTPAlert, index int) string {
	var sb strings.Builder

	// Header with code and confidence
	sb.WriteString(fmt.Sprintf("[%d] %s %s (Confidence: %.2f)\n",
		index,
		ui.ColorCyan.Sprint("🔐"),
		ui.ColorBold.Sprint(otp.OTPCode),
		otp.Confidence,
	))

	// Details
	sb.WriteString(fmt.Sprintf("    From: %s\n", otp.Sender))
	sb.WriteString(fmt.Sprintf("    Received: %s\n", formatTimestamp(otp.Timestamp)))
	sb.WriteString(fmt.Sprintf("    Expires: %s\n", formatExpiry(otp.ExpiresAt)))
	sb.WriteString(fmt.Sprintf("    Source: %s\n", otp.Source))

	if otp.CopiedAt != nil {
		sb.WriteString(fmt.Sprintf("    %s\n", ui.ColorGray.Sprint("✓ Copied to clipboard")))
	}

	return sb.String()
}
