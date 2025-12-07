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
