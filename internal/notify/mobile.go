package notify

import (
	"bytes"
	"fmt"
	"net/http"
)

const ntfyBaseURL = "https://ntfy.sh"

// SendMobileNotification sends a push notification via ntfy.sh
func SendMobileNotification(topic, title, message string) error {
	if topic == "" {
		return fmt.Errorf("ntfy topic is empty")
	}

	url := fmt.Sprintf("%s/%s", ntfyBaseURL, topic)

	// Create request body with title and message
	body := fmt.Sprintf("%s\n\n%s", title, message)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Title", title)
	req.Header.Set("Priority", "high")
	req.Header.Set("Tags", "email,alert")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send mobile notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ntfy.sh returned status %d", resp.StatusCode)
	}

	return nil
}

// SendMobileEmailAlert sends a mobile notification for a matched email
func SendMobileEmailAlert(topic, filterName, from, subject string) error {
	title := fmt.Sprintf("📧 %s", filterName)
	message := fmt.Sprintf("From: %s\nSubject: %s", from, subject)

	return SendMobileNotification(topic, title, message)
}

// SendMobileEmailAlertWithLabels sends a mobile notification for a matched email with labels
func SendMobileEmailAlertWithLabels(topic, filterName string, labels []string, from, subject string) error {
	title := fmt.Sprintf("📧 %s", filterName)
	message := fmt.Sprintf("From: %s\nSubject: %s", from, subject)

	// Add labels as tags for ntfy.sh
	tags := "email,alert"
	if len(labels) > 0 {
		labelsStr := ""
		for _, label := range labels {
			labelsStr += "🏷️ " + label + " "
			tags += "," + label // Add labels as tags for better filtering
		}
		message = fmt.Sprintf("%s\n%s", labelsStr, message)
	}

	// Create custom request to include label tags
	url := fmt.Sprintf("%s/%s", ntfyBaseURL, topic)
	body := fmt.Sprintf("%s\n\n%s", title, message)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers with label tags
	req.Header.Set("Title", title)
	req.Header.Set("Priority", "high")
	req.Header.Set("Tags", tags)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send mobile notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ntfy.sh returned status %d", resp.StatusCode)
	}

	return nil
}
