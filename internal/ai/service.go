package ai

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/datateamsix/email-sentinel/internal/storage"
)

// Service handles AI summary generation with caching and rate limiting
type Service struct {
	provider    Provider
	config      *Config
	db          *sql.DB
	rateLimiter *RateLimiter
	mu          sync.Mutex
}

// RateLimiter tracks API usage to enforce rate limits
type RateLimiter struct {
	hourlyCount int
	dailyCount  int
	hourReset   time.Time
	dayReset    time.Time
	mu          sync.Mutex
}

// NewService creates a new AI summary service
func NewService(cfg *Config, db *sql.DB) (*Service, error) {
	if !cfg.AISummary.Enabled {
		return nil, fmt.Errorf("AI summary is not enabled in configuration")
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	now := time.Now()
	return &Service{
		provider: provider,
		config:   cfg,
		db:       db,
		rateLimiter: &RateLimiter{
			hourReset: now.Add(1 * time.Hour),
			dayReset:  now.Add(24 * time.Hour),
		},
	}, nil
}

// GenerateSummary generates an AI summary for an email
// Returns cached summary if available, otherwise calls the AI provider
func (s *Service) GenerateSummary(messageID, sender, subject, body, snippet string, priority int) (*storage.EmailSummary, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we should skip based on priority
	if s.config.AISummary.Behavior.PriorityOnly && priority != 1 {
		return nil, nil // Skip non-priority emails
	}

	// Check cache first if enabled
	if s.config.AISummary.Behavior.EnableCache {
		cached, err := storage.GetAISummaryByMessageID(s.db, messageID)
		if err != nil {
			log.Printf("⚠️  Error checking cache: %v", err)
		} else if cached != nil {
			log.Printf("🤖 Using cached AI summary for message %s", messageID)
			return cached, nil
		}
	}

	// Check rate limits
	if !s.rateLimiter.CanProceed(s.config.AISummary.RateLimit) {
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// Generate summary
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.AISummary.Behavior.TimeoutSeconds)*time.Second)
	defer cancel()

	req := SummaryRequest{
		Sender:    sender,
		Subject:   subject,
		Body:      body,
		Snippet:   snippet,
		MaxLength: s.config.AISummary.Behavior.MaxSummaryLength,
	}

	log.Printf("🤖 Generating AI summary for: %s", subject)

	var resp *SummaryResponse
	var tokens int
	var err error

	// Retry logic
	maxRetries := s.config.AISummary.Behavior.RetryAttempts
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, tokens, err = s.provider.GenerateSummary(ctx, req)
		if err == nil {
			break
		}

		if attempt < maxRetries {
			log.Printf("⚠️  AI API error (attempt %d/%d): %v", attempt+1, maxRetries+1, err)
			time.Sleep(time.Duration(attempt+1) * time.Second) // Exponential backoff
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries+1, err)
	}

	// Truncate summary if too long
	if len(resp.Summary) > s.config.AISummary.Behavior.MaxSummaryLength {
		resp.Summary = resp.Summary[:s.config.AISummary.Behavior.MaxSummaryLength-3] + "..."
	}

	// Save to database
	summary := &storage.EmailSummary{
		MessageID:   messageID,
		Summary:     resp.Summary,
		Questions:   resp.Questions,
		ActionItems: resp.ActionItems,
		Provider:    s.provider.Name(),
		Model:       s.getModelName(),
		GeneratedAt: time.Now(),
		TokensUsed:  tokens,
	}

	if err := storage.InsertAISummary(s.db, summary); err != nil {
		log.Printf("⚠️  Failed to save AI summary: %v", err)
		// Don't fail - we still return the summary
	}

	// Update rate limiter
	s.rateLimiter.Increment()

	log.Printf("✅ AI summary generated (%d tokens)", tokens)
	return summary, nil
}

// getModelName returns the model name for the current provider
func (s *Service) getModelName() string {
	switch s.provider.Name() {
	case "claude":
		return s.config.AISummary.API.Claude.Model
	case "openai":
		return s.config.AISummary.API.OpenAI.Model
	case "gemini":
		return s.config.AISummary.API.Gemini.Model
	default:
		return "unknown"
	}
}

// CanProceed checks if we can make another API call within rate limits
func (rl *RateLimiter) CanProceed(limits RateLimitConfig) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Reset counters if time windows have passed
	if now.After(rl.hourReset) {
		rl.hourlyCount = 0
		rl.hourReset = now.Add(1 * time.Hour)
	}
	if now.After(rl.dayReset) {
		rl.dailyCount = 0
		rl.dayReset = now.Add(24 * time.Hour)
	}

	// Check limits (0 means unlimited)
	if limits.MaxPerHour > 0 && rl.hourlyCount >= limits.MaxPerHour {
		return false
	}
	if limits.MaxPerDay > 0 && rl.dailyCount >= limits.MaxPerDay {
		return false
	}

	return true
}

// Increment increments the rate limit counters
func (rl *RateLimiter) Increment() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.hourlyCount++
	rl.dailyCount++
}

// GetStats returns current rate limit stats
func (rl *RateLimiter) GetStats() (hourly, daily int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.hourlyCount, rl.dailyCount
}
