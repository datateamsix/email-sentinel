//go:build windows
// +build windows

package notify

import (
	"fmt"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/go-toast/toast"
)

const (
	// AppID is the identifier for Windows toast notifications
	// This makes notifications appear under "Email Sentinel" in Action Center
	AppID = "DataTeamSix.EmailSentinel"

	// Icon paths - using built-in Windows icons
	// For custom icons, you would use absolute paths to .png or .ico files
	IconNormal = "" // Empty uses default system icon
	IconUrgent = "" // Empty uses default system icon
)

// SendAlertNotification sends a Windows toast notification for an email alert
// The notification appears in the Windows Action Center and is clickable
//
// Behavior:
//   - Title: Email subject
//   - Body: "From: <sender>" + AI summary (if available)
//   - Clicking opens the Gmail link in default browser
//   - Priority 1 emails use an urgent visual style
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

	notification := toast.Notification{
		AppID:   AppID,
		Title:   a.Subject,
		Message: message,
		Actions: []toast.Action{
			{
				Type:      "protocol",
				Label:     "Open Email",
				Arguments: a.GmailLink,
			},
		},
		Audio: toast.Default, // System default notification sound
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

		notification.Message = aiMessage
	} else if a.Snippet != "" {
		// Fall back to snippet if no AI summary
		snippet := a.Snippet
		// Truncate snippet if too long (Windows toast has character limits)
		if len(snippet) > 120 { // Reduced from 150 to account for labels
			snippet = snippet[:117] + "..."
		}
		// Append snippet to message
		notification.Message = message + "\n\n" + snippet
	}

	// For priority alerts, use different audio and visual cues
	if a.Priority == 1 {
		// Use reminder audio for urgent alerts (more attention-grabbing)
		notification.Audio = toast.Reminder

		// Add priority indicator to title and message
		notification.Title = "🔥 HIGH PRIORITY: " + a.Subject
	} else {
		// Normal priority - use standard audio
		notification.Audio = toast.Default

		// Add email icon to normal notifications
		notification.Title = "📧 " + a.Subject
	}

	// Push the notification
	err := notification.Push()
	if err != nil {
		RecordDesktopFailure()
		return fmt.Errorf("failed to send Windows toast notification: %w", err)
	}

	RecordDesktopSuccess()
	return nil
}

// SendTestNotification sends a test toast notification to verify Windows notifications work
func SendTestNotification() error {
	testAlert := storage.Alert{
		Subject:    "Email Sentinel Test",
		Sender:     "test@example.com",
		Snippet:    "If you can see this, Windows toast notifications are working! ✅",
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
