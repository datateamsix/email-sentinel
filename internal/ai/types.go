package ai

import "time"

// EmailSummary represents an AI-generated summary of an email
type EmailSummary struct {
	ID          int64
	MessageID   string    // Gmail message ID (foreign key to alerts)
	Summary     string    // Brief overview (max 500 chars)
	Questions   []string  // Key questions detected
	ActionItems []string  // Required actions
	Provider    string    // AI provider used (claude/openai/gemini)
	Model       string    // Model name
	GeneratedAt time.Time // When summary was created
	TokensUsed  int       // API tokens consumed
}

// SummaryRequest represents a request to generate a summary
type SummaryRequest struct {
	Sender  string
	Subject string
	Body    string
	Snippet string
	MaxLength int
}

// SummaryResponse represents the AI provider's response
type SummaryResponse struct {
	Summary     string   `json:"summary"`
	Questions   []string `json:"questions"`
	ActionItems []string `json:"action_items"`
}
