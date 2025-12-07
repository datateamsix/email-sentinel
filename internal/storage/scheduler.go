package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// StartDailyCleanup runs a cleanup task at 12:00 AM every day
// It deletes all alerts from before today (midnight)
// Runs in a goroutine until stopChan is closed
func StartDailyCleanup(db *sql.DB, stopChan <-chan struct{}) {
	for {
		// Calculate time until next midnight
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilMidnight := nextMidnight.Sub(now)

		log.Printf("📅 Daily cleanup scheduled for %s (in %v)", nextMidnight.Format("2006-01-02 15:04:05"), durationUntilMidnight.Round(time.Second))

		select {
		case <-time.After(durationUntilMidnight):
			// It's midnight, run cleanup
			deleted, err := CleanupDailyAlerts(db)
			if err != nil {
				log.Printf("❌ Daily cleanup failed: %v", err)
			} else {
				log.Printf("✅ Daily cleanup completed: deleted %d alert(s) from previous days", deleted)
			}

		case <-stopChan:
			log.Println("🛑 Daily cleanup scheduler stopped")
			return
		}
	}
}

// RunCleanupNow immediately runs the cleanup (useful for testing/manual trigger)
func RunCleanupNow(db *sql.DB) error {
	deleted, err := CleanupDailyAlerts(db)
	if err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	log.Printf("🧹 Manual cleanup completed: deleted %d alert(s)", deleted)
	return nil
}
