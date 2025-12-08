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

// StartOTPCleanup runs OTP cleanup every 1 minute
// It expires inactive OTP codes and deletes old ones
// Runs in a goroutine until stopChan is closed
func StartOTPCleanup(db *sql.DB, stopChan <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("🔐 OTP cleanup scheduler started (runs every 1 minute)")

	// Run immediately on start
	runOTPCleanup(db)

	for {
		select {
		case <-ticker.C:
			runOTPCleanup(db)

		case <-stopChan:
			log.Println("🛑 OTP cleanup scheduler stopped")
			return
		}
	}
}

// runOTPCleanup executes the OTP cleanup tasks
func runOTPCleanup(db *sql.DB) {
	// Mark expired codes as inactive
	expired, err := ExpireOTPAlerts(db)
	if err != nil {
		log.Printf("❌ Failed to expire OTP alerts: %v", err)
	} else if expired > 0 {
		log.Printf("🔐 Expired %d OTP alert(s)", expired)
	}

	// Delete old codes (older than 24h)
	deleted, err := DeleteExpiredOTPAlerts(db)
	if err != nil {
		log.Printf("❌ Failed to delete old OTP alerts: %v", err)
	} else if deleted > 0 {
		log.Printf("🔐 Deleted %d old OTP alert(s)", deleted)
	}
}
