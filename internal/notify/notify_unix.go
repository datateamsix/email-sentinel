//go:build !windows
// +build !windows

package notify

import (
	"fmt"

	"github.com/datateamsix/email-sentinel/internal/storage"
)

// SendAlertNotification sends a desktop notification for an email alert
// On Linux/macOS, this uses the beeep library for cross-platform notifications
//
// Behavior:
//   - Title: Email subject with priority indicator
//   - Body: "From: <sender>" + AI summary (if available)
//   - Priority 1 emails show 🔥 HIGH PRIORITY indicator
//   - AI-summarized emails show 🤖 icon and summary
func SendAlertNotification(a storage.Alert) error {
	// Build message with filter labels if present
	message := fmt.Sprintf("From: %s", a.Sender)
	if len(a.FilterLabels) > 0 {
		labelsStr := ""
		for _, label := range a.FilterLabels {
			labelsStr += "🏷️ " + label + " "
		}
		message = labelsStr + "\n" + message
	}

	// Prioritize AI summary over snippet if available
	if a.AISummary != nil && a.AISummary.Summary != "" {
		// Use AI summary instead of snippet
		aiMessage := message + "\n\n🤖 " + a.AISummary.Summary

		// Add questions if present (max 2 for space)
		if len(a.AISummary.Questions) > 0 {
			aiMessage += "\n\n❓ "
			if len(a.AISummary.Questions) == 1 {
				aiMessage += a.AISummary.Questions[0]
			} else {
				aiMessage += fmt.Sprintf("%s (+ %d more)", a.AISummary.Questions[0], len(a.AISummary.Questions)-1)
			}
		}

		// Add action items if present (max 2 for space)
		if len(a.AISummary.ActionItems) > 0 {
			aiMessage += "\n✅ "
			if len(a.AISummary.ActionItems) == 1 {
				aiMessage += a.AISummary.ActionItems[0]
			} else {
				aiMessage += fmt.Sprintf("%s (+ %d more)", a.AISummary.ActionItems[0], len(a.AISummary.ActionItems)-1)
			}
		}

		message = aiMessage
	} else if a.Snippet != "" {
		// Fall back to snippet if no AI summary
		snippet := a.Snippet
		// Truncate snippet if too long
		if len(snippet) > 120 {
			snippet = snippet[:117] + "..."
		}
		// Append snippet to message
		message = message + "\n\n" + snippet
	}

	// Build title with priority indicator
	var title string
	if a.Priority == 1 {
		title = "🔥 HIGH PRIORITY: " + a.Subject
	} else {
		title = "📧 " + a.Subject
	}

	// Send using cross-platform desktop notification
	return SendDesktopNotification(title, message)
}

// SendTestNotification sends a test desktop notification to verify notifications work
func SendTestNotification() error {
	testAlert := storage.Alert{
		Subject:    "Email Sentinel Test",
		Sender:     "test@example.com",
		Snippet:    "If you can see this, desktop notifications are working! ✅",
		GmailLink:  "https://mail.google.com",
		FilterName: "Test",
		Priority:   0,
	}

	return SendAlertNotification(testAlert)
}

// SendPriorityTestNotification sends a test high-priority notification
func SendPriorityTestNotification() error {
	testAlert := storage.Alert{
		Subject:    "URGENT: This is a high priority test",
		Sender:     "boss@company.com",
		Snippet:    "This demonstrates how urgent alerts appear with different styling.",
		GmailLink:  "https://mail.google.com",
		FilterName: "VIP Sender",
		Priority:   1,
	}

	return SendAlertNotification(testAlert)
}
