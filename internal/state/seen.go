package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/datateamsix/email-sentinel/internal/config"
)

// SeenMessages tracks which message IDs have been processed
type SeenMessages struct {
	mu       sync.RWMutex
	messages map[string]time.Time // message ID -> timestamp when seen
	filePath string
}

// State represents the persistent state file
type State struct {
	SeenMessages []SeenMessage `json:"seen_messages"`
}

// SeenMessage represents a seen message with timestamp
type SeenMessage struct {
	ID        string    `json:"id"`
	SeenAt    time.Time `json:"seen_at"`
}

// NewSeenMessages creates a new SeenMessages tracker
func NewSeenMessages() (*SeenMessages, error) {
	configDir, err := config.ConfigDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(configDir, "seen_messages.json")

	sm := &SeenMessages{
		messages: make(map[string]time.Time),
		filePath: filePath,
	}

	// Load existing state if it exists
	if err := sm.load(); err != nil {
		// If file doesn't exist, that's okay - we'll start fresh
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Cleanup old messages (older than 30 days)
	sm.CleanupOld(30 * 24 * time.Hour)

	return sm, nil
}

// IsSeen checks if a message ID has been seen
func (sm *SeenMessages) IsSeen(messageID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.messages[messageID]
	return exists
}

// MarkSeen marks a message ID as seen with current timestamp
func (sm *SeenMessages) MarkSeen(messageID string) error {
	sm.mu.Lock()
	sm.messages[messageID] = time.Now()
	sm.mu.Unlock()

	return sm.save()
}

// MarkMultipleSeen marks multiple message IDs as seen
func (sm *SeenMessages) MarkMultipleSeen(messageIDs []string) error {
	now := time.Now()
	sm.mu.Lock()
	for _, id := range messageIDs {
		sm.messages[id] = now
	}
	sm.mu.Unlock()

	return sm.save()
}

// Count returns the number of seen messages
func (sm *SeenMessages) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.messages)
}

// Clear removes all seen messages
func (sm *SeenMessages) Clear() error {
	sm.mu.Lock()
	sm.messages = make(map[string]time.Time)
	sm.mu.Unlock()

	return sm.save()
}

// CleanupOld removes message IDs older than the specified duration
func (sm *SeenMessages) CleanupOld(maxAge time.Duration) int {
	cutoff := time.Now().Add(-maxAge)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	cleaned := 0
	for id, seenAt := range sm.messages {
		if seenAt.Before(cutoff) {
			delete(sm.messages, id)
			cleaned++
		}
	}

	if cleaned > 0 {
		sm.save() // Save after cleanup
	}

	return cleaned
}

// load reads the state from disk
func (sm *SeenMessages) load() error {
	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, msg := range state.SeenMessages {
		sm.messages[msg.ID] = msg.SeenAt
	}

	return nil
}

// save writes the state to disk
func (sm *SeenMessages) save() error {
	sm.mu.RLock()
	messages := make([]SeenMessage, 0, len(sm.messages))
	for id, seenAt := range sm.messages {
		messages = append(messages, SeenMessage{
			ID:     id,
			SeenAt: seenAt,
		})
	}
	sm.mu.RUnlock()

	state := State{
		SeenMessages: messages,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if _, err := config.EnsureConfigDir(); err != nil {
		return err
	}

	return os.WriteFile(sm.filePath, data, 0600)
}
