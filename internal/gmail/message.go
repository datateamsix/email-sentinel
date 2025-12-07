package gmail

import (
	"strings"

	"google.golang.org/api/gmail/v1"
)

// EmailMessage represents a parsed email message
type EmailMessage struct {
	ID      string
	From    string
	Subject string
	Snippet string
	Date    string
}

// ParseMessage extracts relevant fields from a Gmail API message
func ParseMessage(msg *gmail.Message) *EmailMessage {
	email := &EmailMessage{
		ID:      msg.Id,
		Snippet: msg.Snippet,
	}

	// Extract headers
	for _, header := range msg.Payload.Headers {
		switch strings.ToLower(header.Name) {
		case "from":
			email.From = header.Value
		case "subject":
			email.Subject = header.Value
		case "date":
			email.Date = header.Value
		}
	}

	return email
}

// GetFromAddress extracts just the email address from a "From" header
// Example: "John Doe <john@example.com>" -> "john@example.com"
func GetFromAddress(from string) string {
	// Check if email is in format "Name <email@domain.com>"
	if start := strings.Index(from, "<"); start != -1 {
		if end := strings.Index(from, ">"); end != -1 {
			return strings.TrimSpace(from[start+1 : end])
		}
	}

	// Otherwise return the whole string (it's probably just the email)
	return strings.TrimSpace(from)
}

// GetFromDomain extracts the domain from an email address
// Example: "john@example.com" -> "example.com"
func GetFromDomain(email string) string {
	address := GetFromAddress(email)
	if idx := strings.Index(address, "@"); idx != -1 {
		return address[idx+1:]
	}
	return address
}

// BuildGmailLink generates a stable Gmail permalink from a message ID
// The link uses #all/ to work regardless of which label/folder the email is in
// Example: BuildGmailLink("abc123") -> "https://mail.google.com/mail/u/0/#all/abc123"
func BuildGmailLink(messageID string) string {
	return "https://mail.google.com/mail/u/0/#all/" + messageID
}
