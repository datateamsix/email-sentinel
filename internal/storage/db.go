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
	ID         int64
	Timestamp  time.Time
	Sender     string
	Subject    string
	Snippet    string
	Labels     string
	MessageID  string
	GmailLink  string
	FilterName string
	Priority   int
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

	return scanAlerts(rows)
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
