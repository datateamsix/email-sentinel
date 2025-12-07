package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

// SendDesktopNotification sends a native OS notification
func SendDesktopNotification(title, message string) error {
	// Use beeep to send cross-platform notification
	err := beeep.Notify(title, message, "")
	if err != nil {
		return fmt.Errorf("failed to send desktop notification: %w", err)
	}
	return nil
}

// SendEmailAlert sends a desktop notification for a matched email
func SendEmailAlert(filterName, from, subject string) error {
	title := fmt.Sprintf("📧 Email Match: %s", filterName)
	message := fmt.Sprintf("From: %s\nSubject: %s", from, subject)

	return SendDesktopNotification(title, message)
}

// SendEmailAlertWithLabels sends a desktop notification for a matched email with labels
func SendEmailAlertWithLabels(filterName string, labels []string, from, subject string) error {
	title := fmt.Sprintf("📧 Email Match: %s", filterName)
	message := fmt.Sprintf("From: %s\nSubject: %s", from, subject)

	if len(labels) > 0 {
		labelsStr := ""
		for _, label := range labels {
			labelsStr += "🏷️ " + label + " "
		}
		message = fmt.Sprintf("%s\n%s", labelsStr, message)
	}

	return SendDesktopNotification(title, message)
}
