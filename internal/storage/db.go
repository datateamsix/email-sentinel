package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/datateamsix/email-sentinel/internal/config"

	_ "modernc.org/sqlite"
)

// Alert represents an email notification stored in the database
type Alert struct {
	ID           int64
	Timestamp    time.Time
	Sender       string
	Subject      string
	Snippet      string
	Labels       string   // Gmail labels
	MessageID    string
	GmailLink    string
	FilterName   string
	FilterLabels []string // Filter categories (not stored in DB, populated at runtime)
	Priority     int
}

// OTPAlert represents an OTP code extracted from an email
type OTPAlert struct {
	ID          int64
	Timestamp   time.Time
	ExpiresAt   time.Time
	Sender      string
	Subject     string
	OTPCode     string
	Confidence  float64
	Source      string
	PatternName string
	MessageID   string
	GmailLink   string
	FilterName  string
	IsActive    bool
	CopiedAt    *time.Time // Nullable timestamp
}

const schema = `
CREATE TABLE IF NOT EXISTS alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    sender TEXT NOT NULL,
    subject TEXT NOT NULL,
    snippet TEXT,
    labels TEXT,
    message_id TEXT NOT NULL UNIQUE,
    gmail_link TEXT NOT NULL,
    filter_name TEXT NOT NULL,
    priority INTEGER DEFAULT 0 CHECK(priority IN (0, 1))
);

CREATE INDEX IF NOT EXISTS idx_timestamp ON alerts(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_message_id ON alerts(message_id);

CREATE TABLE IF NOT EXISTS filter_labels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    label TEXT NOT NULL UNIQUE COLLATE NOCASE,
    created_at INTEGER NOT NULL,
    last_used INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_label ON filter_labels(label COLLATE NOCASE);
`

// InitDB initializes the SQLite database and creates tables if needed
func InitDB() (*sql.DB, error) {
	configDir, err := config.EnsureConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	dbPath := filepath.Join(configDir, "history.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Limit to 1 connection to prevent write conflicts
	db.SetMaxOpenConns(1)

	// Create tables and indexes
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// Run migrations for new features (like OTP alerts)
	if err := RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}

// InsertAlert saves a new alert to the database
// If the message_id already exists, it returns an error (duplicate)
func InsertAlert(db *sql.DB, a *Alert) error {
	query := `
		INSERT INTO alerts (timestamp, sender, subject, snippet, labels, message_id, gmail_link, filter_name, priority)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(
		query,
		a.Timestamp.Unix(),
		a.Sender,
		a.Subject,
		a.Snippet,
		a.Labels,
		a.MessageID,
		a.GmailLink,
		a.FilterName,
		a.Priority,
	)

	if err != nil {
		return fmt.Errorf("failed to insert alert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get insert ID: %w", err)
	}

	a.ID = id
	return nil
}

// GetTodayAlerts returns all alerts from today (since midnight)
func GetTodayAlerts(db *sql.DB) ([]Alert, error) {
	// Get today's midnight
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return getAlertsSince(db, midnight)
}

// GetRecentAlerts returns the N most recent alerts
func GetRecentAlerts(db *sql.DB, limit int) ([]Alert, error) {
	query := `
		SELECT id, timestamp, sender, subject, snippet, labels, message_id, gmail_link, filter_name, priority
		FROM alerts
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent alerts: %w", err)
	}
	defer rows.Close()

	alerts, err := scanAlerts(rows)
	if err != nil {
		return nil, err
	}

	// Populate FilterLabels from filter configuration
	if err := PopulateFilterLabels(alerts); err != nil {
		// Log error but don't fail - alerts can still be shown
		fmt.Printf("Warning: Could not populate filter labels: %v\n", err)
	}

	return alerts, nil
}

// getAlertsSince returns all alerts since the given time
func getAlertsSince(db *sql.DB, since time.Time) ([]Alert, error) {
	query := `
		SELECT id, timestamp, sender, subject, snippet, labels, message_id, gmail_link, filter_name, priority
		FROM alerts
		WHERE timestamp >= ?
		ORDER BY timestamp DESC
	`

	rows, err := db.Query(query, since.Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %w", err)
	}
	defer rows.Close()

	return scanAlerts(rows)
}

// CountTodayAlerts returns the count of alerts since midnight
func CountTodayAlerts(db *sql.DB) (int, error) {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	query := "SELECT COUNT(*) FROM alerts WHERE timestamp >= ?"
	var count int
	err := db.QueryRow(query, midnight.Unix()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count alerts: %w", err)
	}

	return count, nil
}

// DeleteAlertsBefore deletes all alerts older than the given time
func DeleteAlertsBefore(db *sql.DB, cutoff time.Time) (int64, error) {
	query := "DELETE FROM alerts WHERE timestamp < ?"
	result, err := db.Exec(query, cutoff.Unix())
	if err != nil {
		return 0, fmt.Errorf("failed to delete old alerts: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get deleted count: %w", err)
	}

	return deleted, nil
}

// CleanupDailyAlerts deletes all alerts from before today (midnight)
// This is called at 12:00 AM daily to wipe yesterday's alerts
func CleanupDailyAlerts(db *sql.DB) (int64, error) {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	deleted, err := DeleteAlertsBefore(db, midnight)
	if err != nil {
		return 0, fmt.Errorf("daily cleanup failed: %w", err)
	}

	return deleted, nil
}

// DeleteAllAlerts deletes all alerts from the database
// Returns the number of alerts deleted
func DeleteAllAlerts(db *sql.DB) (int64, error) {
	query := "DELETE FROM alerts"
	result, err := db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete all alerts: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get deleted count: %w", err)
	}

	return deleted, nil
}

// DeleteAlerts24HoursOld deletes alerts older than 24 hours
// Returns the number of alerts deleted
func DeleteAlerts24HoursOld(db *sql.DB) (int64, error) {
	cutoff := time.Now().Add(-24 * time.Hour)
	deleted, err := DeleteAlertsBefore(db, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to delete 24-hour-old alerts: %w", err)
	}
	return deleted, nil
}

// scanAlerts is a helper function to scan rows into Alert structs
func scanAlerts(rows *sql.Rows) ([]Alert, error) {
	var alerts []Alert

	for rows.Next() {
		var a Alert
		var timestamp int64

		err := rows.Scan(
			&a.ID,
			&timestamp,
			&a.Sender,
			&a.Subject,
			&a.Snippet,
			&a.Labels,
			&a.MessageID,
			&a.GmailLink,
			&a.FilterName,
			&a.Priority,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		a.Timestamp = time.Unix(timestamp, 0)
		alerts = append(alerts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating alerts: %w", err)
	}

	return alerts, nil
}

// SaveLabel saves or updates a label in the database
// If the label already exists, updates its last_used timestamp
func SaveLabel(db *sql.DB, label string) error {
	now := time.Now().Unix()

	query := `
		INSERT INTO filter_labels (label, created_at, last_used)
		VALUES (?, ?, ?)
		ON CONFLICT(label) DO UPDATE SET last_used = ?
	`

	_, err := db.Exec(query, label, now, now, now)
	if err != nil {
		return fmt.Errorf("failed to save label: %w", err)
	}

	return nil
}

// GetAllLabels returns all labels ordered by most recently used
func GetAllLabels(db *sql.DB) ([]string, error) {
	query := `SELECT label FROM filter_labels ORDER BY last_used DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query labels: %w", err)
	}
	defer rows.Close()

	var labels []string
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return nil, fmt.Errorf("failed to scan label: %w", err)
		}
		labels = append(labels, label)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating labels: %w", err)
	}

	return labels, nil
}

// SaveLabels saves multiple labels at once
func SaveLabels(db *sql.DB, labels []string) error {
	for _, label := range labels {
		if err := SaveLabel(db, label); err != nil {
			return err
		}
	}
	return nil
}

// ======================================
// OTP Alert Functions
// ======================================

// InsertOTPAlert saves a new OTP alert to the database
func InsertOTPAlert(db *sql.DB, otp *OTPAlert) error {
	query := `
		INSERT INTO otp_alerts (
			timestamp, expires_at, sender, subject, otp_code, confidence,
			source, pattern_name, message_id, gmail_link, filter_name, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(
		query,
		otp.Timestamp.Unix(),
		otp.ExpiresAt.Unix(),
		otp.Sender,
		otp.Subject,
		otp.OTPCode,
		otp.Confidence,
		otp.Source,
		otp.PatternName,
		otp.MessageID,
		otp.GmailLink,
		otp.FilterName,
		boolToInt(otp.IsActive),
	)

	if err != nil {
		return fmt.Errorf("failed to insert OTP alert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get insert ID: %w", err)
	}

	otp.ID = id
	return nil
}

// GetActiveOTPAlerts returns all OTP alerts that are active and not expired
func GetActiveOTPAlerts(db *sql.DB) ([]OTPAlert, error) {
	query := `
		SELECT
			id, timestamp, expires_at, sender, subject, otp_code, confidence,
			source, pattern_name, message_id, gmail_link, filter_name, is_active, copied_at
		FROM otp_alerts
		WHERE is_active = 1 AND expires_at > ?
		ORDER BY timestamp DESC
	`

	now := time.Now().Unix()
	rows, err := db.Query(query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to query active OTP alerts: %w", err)
	}
	defer rows.Close()

	return scanOTPAlerts(rows)
}

// GetRecentOTPAlerts returns the N most recent OTP alerts regardless of status
func GetRecentOTPAlerts(db *sql.DB, limit int) ([]OTPAlert, error) {
	query := `
		SELECT
			id, timestamp, expires_at, sender, subject, otp_code, confidence,
			source, pattern_name, message_id, gmail_link, filter_name, is_active, copied_at
		FROM otp_alerts
		ORDER BY timestamp DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent OTP alerts: %w", err)
	}
	defer rows.Close()

	return scanOTPAlerts(rows)
}

// MarkOTPAsCopied updates the copied_at timestamp for an OTP alert
func MarkOTPAsCopied(db *sql.DB, id int64) error {
	query := "UPDATE otp_alerts SET copied_at = ? WHERE id = ?"

	result, err := db.Exec(query, time.Now().Unix(), id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as copied: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("OTP alert with ID %d not found", id)
	}

	return nil
}

// ExpireOTPAlerts marks all expired OTP alerts as inactive
// Returns the number of alerts that were expired
func ExpireOTPAlerts(db *sql.DB) (int64, error) {
	query := "UPDATE otp_alerts SET is_active = 0 WHERE expires_at <= ? AND is_active = 1"

	result, err := db.Exec(query, time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("failed to expire OTP alerts: %w", err)
	}

	expired, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get expired count: %w", err)
	}

	return expired, nil
}

// DeleteExpiredOTPAlerts deletes OTP alerts older than 24 hours
// Returns the number of alerts that were deleted
func DeleteExpiredOTPAlerts(db *sql.DB) (int64, error) {
	cutoff := time.Now().Add(-24 * time.Hour).Unix()
	query := "DELETE FROM otp_alerts WHERE timestamp < ?"

	result, err := db.Exec(query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired OTP alerts: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get deleted count: %w", err)
	}

	return deleted, nil
}

// scanOTPAlerts is a helper function to scan rows into OTPAlert structs
func scanOTPAlerts(rows *sql.Rows) ([]OTPAlert, error) {
	var alerts []OTPAlert

	for rows.Next() {
		var otp OTPAlert
		var timestamp, expiresAt int64
		var copiedAt sql.NullInt64
		var isActive int

		err := rows.Scan(
			&otp.ID,
			&timestamp,
			&expiresAt,
			&otp.Sender,
			&otp.Subject,
			&otp.OTPCode,
			&otp.Confidence,
			&otp.Source,
			&otp.PatternName,
			&otp.MessageID,
			&otp.GmailLink,
			&otp.FilterName,
			&isActive,
			&copiedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan OTP alert: %w", err)
		}

		otp.Timestamp = time.Unix(timestamp, 0)
		otp.ExpiresAt = time.Unix(expiresAt, 0)
		otp.IsActive = isActive == 1

		if copiedAt.Valid {
			t := time.Unix(copiedAt.Int64, 0)
			otp.CopiedAt = &t
		}

		alerts = append(alerts, otp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating OTP alerts: %w", err)
	}

	return alerts, nil
}

// boolToInt converts a boolean to an integer (0 or 1) for SQLite storage
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// PopulateFilterLabels populates the FilterLabels field for alerts by loading
// the filter configuration and matching filter names
func PopulateFilterLabels(alerts []Alert) error {
	// Import filter package to load config
	// We need to do this dynamically to avoid import cycles
	// For now, we'll use a simpler approach: check the filter name for common patterns

	for i := range alerts {
		// For now, use a simple heuristic: check if filter name contains "otp"
		// This can be enhanced later to load actual filter config
		filterNameLower := ""
		for _, ch := range alerts[i].FilterName {
			if ch >= 'A' && ch <= 'Z' {
				filterNameLower += string(ch + 32)
			} else {
				filterNameLower += string(ch)
			}
		}

		// Check if filter name suggests OTP
		if containsSubstring(filterNameLower, "otp") ||
		   containsSubstring(filterNameLower, "code") ||
		   containsSubstring(filterNameLower, "verification") ||
		   containsSubstring(filterNameLower, "2fa") ||
		   containsSubstring(filterNameLower, "authentication") {
			alerts[i].FilterLabels = []string{"otp"}
		}
	}

	return nil
}

// containsSubstring checks if a string contains a substring (simple implementation)
func containsSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
