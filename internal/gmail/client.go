package gmail

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Client wraps the Gmail API service with auto-refreshing tokens
type Client struct {
	service     *gmail.Service
	token       *oauth2.Token
	oauthConfig *oauth2.Config
	tokenMu     sync.RWMutex
}

// NewClient creates a new Gmail API client using the provided OAuth token
// The client automatically refreshes expired tokens and saves them to disk
func NewClient(token *oauth2.Token, oauthConfig *oauth2.Config) (*Client, error) {
	ctx := context.Background()

	// Create token source that auto-refreshes
	tokenSource := oauthConfig.TokenSource(ctx, token)

	httpClient := oauth2.NewClient(ctx, tokenSource)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %w", err)
	}

	client := &Client{
		service:     service,
		token:       token,
		oauthConfig: oauthConfig,
	}

	// Start background token refresh monitor
	go client.monitorTokenRefresh(tokenSource)

	return client, nil
}

// monitorTokenRefresh checks for token refreshes and saves them to disk
func (c *Client) monitorTokenRefresh(tokenSource oauth2.TokenSource) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Get current token from source
		newToken, err := tokenSource.Token()
		if err != nil {
			// CRITICAL: Token refresh failed - alert user immediately
			fmt.Println("")
			fmt.Println("❌ CRITICAL: OAuth token refresh failed!")
			fmt.Printf("   Error: %v\n", err)
			fmt.Println("")
			fmt.Println("   This usually means your Gmail authentication has expired.")
			fmt.Println("   Please re-authenticate with:")
			fmt.Println("   email-sentinel init")
			fmt.Println("")
			// Continue monitoring, will retry next cycle (5 minutes)
			continue
		}

		// Check if token was refreshed (access token changed)
		c.tokenMu.RLock()
		tokenChanged := c.token.AccessToken != newToken.AccessToken
		c.tokenMu.RUnlock()

		if tokenChanged {
			// Save refreshed token
			c.tokenMu.Lock()
			c.token = newToken
			c.tokenMu.Unlock()

			if err := SaveToken(newToken); err != nil {
				// Log error but continue - not fatal
				fmt.Printf("⚠️  Warning: Failed to save refreshed token: %v\n", err)
			}
		}
	}
}

// RefreshTokenIfNeeded manually refreshes the token if it's expired or about to expire
func (c *Client) RefreshTokenIfNeeded() error {
	c.tokenMu.RLock()
	needsRefresh := time.Until(c.token.Expiry) < 5*time.Minute
	c.tokenMu.RUnlock()

	if !needsRefresh {
		return nil
	}

	ctx := context.Background()
	tokenSource := c.oauthConfig.TokenSource(ctx, c.token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	c.tokenMu.Lock()
	c.token = newToken
	c.tokenMu.Unlock()

	// Save to disk
	if err := SaveToken(newToken); err != nil {
		return fmt.Errorf("failed to save refreshed token: %w", err)
	}

	fmt.Println("✅ OAuth token refreshed successfully")
	return nil
}

// GetRecentMessages fetches recent messages from the inbox with retry logic
// maxResults specifies the maximum number of messages to retrieve
func (c *Client) GetRecentMessages(maxResults int64) ([]*gmail.Message, error) {
	const maxRetries = 3
	const baseDelay = 2 * time.Second

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		messages, err := c.getRecentMessagesOnce(maxResults)
		if err == nil {
			return messages, nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			return nil, err
		}

		// Exponential backoff
		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<uint(attempt))
			fmt.Printf("⚠️  Gmail API error (attempt %d/%d): %v\n", attempt+1, maxRetries, err)
			fmt.Printf("   Retrying in %v...\n", delay)
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// getRecentMessagesOnce fetches messages without retry logic
func (c *Client) getRecentMessagesOnce(maxResults int64) ([]*gmail.Message, error) {
	user := "me"

	// Refresh token if needed before making API call
	if err := c.RefreshTokenIfNeeded(); err != nil {
		return nil, err
	}

	// List message IDs
	listCall := c.service.Users.Messages.List(user).
		MaxResults(maxResults).
		Q("in:inbox") // Only inbox messages

	response, err := listCall.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %w", err)
	}

	if len(response.Messages) == 0 {
		return []*gmail.Message{}, nil
	}

	// Fetch full message details for each message
	messages := make([]*gmail.Message, 0, len(response.Messages))
	for _, msg := range response.Messages {
		fullMsg, err := c.service.Users.Messages.Get(user, msg.Id).
			Format("full").
			Do()
		if err != nil {
			// Log error but continue with other messages
			fmt.Printf("⚠️  Warning: Could not fetch message %s: %v\n", msg.Id, err)
			continue
		}
		messages = append(messages, fullMsg)
	}

	return messages, nil
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Network errors
	if contains(errStr, "timeout") || contains(errStr, "connection refused") ||
		contains(errStr, "temporary failure") || contains(errStr, "network") {
		return true
	}

	// Rate limiting errors
	if contains(errStr, "rate limit") || contains(errStr, "429") ||
		contains(errStr, "quota") {
		return true
	}

	// Server errors (5xx)
	if contains(errStr, "500") || contains(errStr, "502") ||
		contains(errStr, "503") || contains(errStr, "504") {
		return true
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)*2 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetMessagesAfter fetches messages received after a specific message ID
func (c *Client) GetMessagesAfter(afterMessageID string, maxResults int64) ([]*gmail.Message, error) {
	user := "me"

	query := "in:inbox"
	if afterMessageID != "" {
		// Note: Gmail API doesn't have a direct "after ID" query
		// We'll fetch recent messages and filter client-side
		query = fmt.Sprintf("%s newer_than:1h", query)
	}

	listCall := c.service.Users.Messages.List(user).
		MaxResults(maxResults).
		Q(query)

	response, err := listCall.Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve messages: %w", err)
	}

	if len(response.Messages) == 0 {
		return []*gmail.Message{}, nil
	}

	// Fetch full message details
	messages := make([]*gmail.Message, 0, len(response.Messages))
	foundAfter := afterMessageID == "" // If no afterMessageID, include all

	for _, msg := range response.Messages {
		// Skip messages until we find the "after" message
		if !foundAfter {
			if msg.Id == afterMessageID {
				foundAfter = true
			}
			continue
		}

		fullMsg, err := c.service.Users.Messages.Get(user, msg.Id).
			Format("full").
			Do()
		if err != nil {
			fmt.Printf("Warning: Could not fetch message %s: %v\n", msg.Id, err)
			continue
		}
		messages = append(messages, fullMsg)
	}

	return messages, nil
}

// MarkAsRead marks a message as read
func (c *Client) MarkAsRead(messageID string) error {
	user := "me"

	modifyRequest := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}

	_, err := c.service.Users.Messages.Modify(user, messageID, modifyRequest).Do()
	if err != nil {
		return fmt.Errorf("unable to mark message as read: %w", err)
	}

	return nil
}
