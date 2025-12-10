package notify

import (
	"fmt"
	"sync"
	"time"
)

// NotificationHealth tracks the health status of notification delivery
type NotificationHealth struct {
	desktopFailures int
	mobileFailures  int
	lastDesktopOK   time.Time
	lastMobileOK    time.Time
	mu              sync.RWMutex
}

var (
	health         = &NotificationHealth{}
	healthMu       sync.Mutex
	lastHealthWarn time.Time
)

const (
	// After 5 consecutive failures, warn the user
	failureThreshold = 5
	// Only warn once per hour to avoid spam
	warnCooldown = 1 * time.Hour
)

// RecordDesktopSuccess records a successful desktop notification
func RecordDesktopSuccess() {
	health.mu.Lock()
	defer health.mu.Unlock()

	health.desktopFailures = 0
	health.lastDesktopOK = time.Now()
}

// RecordDesktopFailure records a failed desktop notification
func RecordDesktopFailure() {
	health.mu.Lock()
	defer health.mu.Unlock()

	health.desktopFailures++
	checkAndWarnDesktop()
}

// RecordMobileSuccess records a successful mobile notification
func RecordMobileSuccess() {
	health.mu.Lock()
	defer health.mu.Unlock()

	health.mobileFailures = 0
	health.lastMobileOK = time.Now()
}

// RecordMobileFailure records a failed mobile notification
func RecordMobileFailure() {
	health.mu.Lock()
	defer health.mu.Unlock()

	health.mobileFailures++
	checkAndWarnMobile()
}

// checkAndWarnDesktop checks if desktop notifications are persistently failing
func checkAndWarnDesktop() {
	if health.desktopFailures >= failureThreshold {
		healthMu.Lock()
		defer healthMu.Unlock()

		// Only warn once per hour
		if time.Since(lastHealthWarn) > warnCooldown {
			fmt.Printf("\n⚠️  WARNING: Desktop notifications have failed %d times in a row\n", health.desktopFailures)
			fmt.Println("   This may indicate a system notification issue.")
			fmt.Println("   Check your OS notification settings.")
			lastHealthWarn = time.Now()
		}
	}
}

// checkAndWarnMobile checks if mobile notifications are persistently failing
func checkAndWarnMobile() {
	if health.mobileFailures >= failureThreshold {
		healthMu.Lock()
		defer healthMu.Unlock()

		// Only warn once per hour
		if time.Since(lastHealthWarn) > warnCooldown {
			fmt.Printf("\n⚠️  WARNING: Mobile notifications have failed %d times in a row\n", health.mobileFailures)
			fmt.Println("   This may indicate a network or ntfy.sh connectivity issue.")
			fmt.Println("   Check your internet connection and ntfy topic configuration.")
			lastHealthWarn = time.Now()
		}
	}
}

// GetHealthStatus returns the current notification health status
func GetHealthStatus() (desktopOK bool, mobileOK bool, desktopFailCount int, mobileFailCount int) {
	health.mu.RLock()
	defer health.mu.RUnlock()

	desktopOK = health.desktopFailures < failureThreshold
	mobileOK = health.mobileFailures < failureThreshold
	desktopFailCount = health.desktopFailures
	mobileFailCount = health.mobileFailures

	return
}

// ResetHealth resets all health counters (useful for testing or after fixing issues)
func ResetHealth() {
	health.mu.Lock()
	defer health.mu.Unlock()

	health.desktopFailures = 0
	health.mobileFailures = 0
	health.lastDesktopOK = time.Time{}
	health.lastMobileOK = time.Time{}
}
