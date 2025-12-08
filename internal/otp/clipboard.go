/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package otp

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

var (
	lastCopiedCode  string
	autoClearTimer  *time.Timer
	clipboardActive bool
)

// CopyToClipboard copies an OTP code to the system clipboard
func CopyToClipboard(code string) error {
	err := clipboard.WriteAll(code)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	lastCopiedCode = code
	clipboardActive = true

	return nil
}

// ScheduleAutoClear schedules the clipboard to be cleared after the given duration
func ScheduleAutoClear(duration time.Duration) {
	// Cancel any existing timer
	if autoClearTimer != nil {
		autoClearTimer.Stop()
	}

	autoClearTimer = time.AfterFunc(duration, func() {
		if clipboardActive {
			// Clear clipboard
			clipboard.WriteAll("")

			// Zero out the last copied code
			SecureZeroString(&lastCopiedCode)

			clipboardActive = false
		}
	})
}

// SecureZeroString overwrites a string in memory (best effort)
// Note: Go strings are immutable, so this works on the backing array
func SecureZeroString(s *string) {
	if s == nil {
		return
	}

	// Convert to byte slice and zero it
	bytes := []byte(*s)
	for i := range bytes {
		bytes[i] = 0
	}

	*s = ""
}

// GetClipboard retrieves the current clipboard content
func GetClipboard() (string, error) {
	return clipboard.ReadAll()
}

// ClearClipboard immediately clears the clipboard
func ClearClipboard() error {
	err := clipboard.WriteAll("")
	if err != nil {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}

	SecureZeroString(&lastCopiedCode)
	clipboardActive = false

	return nil
}

// IsClipboardActive returns whether the clipboard contains an OTP code
func IsClipboardActive() bool {
	return clipboardActive
}
