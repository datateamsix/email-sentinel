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
//   - Body: "From: <sender>"
//   - Clicking opens the Gmail link in default browser
//   - Priority 1 emails use an urgent visual style
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

	// Add snippet as additional context if available
	if a.Snippet != "" {
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
		return fmt.Errorf("failed to send Windows toast notification: %w", err)
	}

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
