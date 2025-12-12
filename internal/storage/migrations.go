/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// GetSchemaVersion retrieves the current schema version from the database
// Returns 0 if the schema_version table doesn't exist yet
func GetSchemaVersion(db *sql.DB) (int, error) {
	// Check if schema_version table exists
	query := `
		SELECT name FROM sqlite_master
		WHERE type='table' AND name='schema_version'
	`
	var tableName string
	err := db.QueryRow(query).Scan(&tableName)
	if err == sql.ErrNoRows {
		// Table doesn't exist, this is version 0
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to check schema_version table: %w", err)
	}

	// Get current version
	var version int
	err = db.QueryRow("SELECT version FROM schema_version ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get schema version: %w", err)
	}

	return version, nil
}

// setSchemaVersion updates the schema version in the database
func setSchemaVersion(tx *sql.Tx, version int) error {
	// Create schema_version table if it doesn't exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER NOT NULL,
			applied_at INTEGER NOT NULL
		)
	`
	if _, err := tx.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create schema_version table: %w", err)
	}

	// Insert new version
	insertSQL := "INSERT INTO schema_version (version, applied_at) VALUES (?, ?)"
	if _, err := tx.Exec(insertSQL, version, time.Now().Unix()); err != nil {
		return fmt.Errorf("failed to insert schema version: %w", err)
	}

	return nil
}

// RunMigrations executes all pending database migrations
// Each migration is run in a transaction and rolled back on failure
func RunMigrations(db *sql.DB) error {
	currentVersion, err := GetSchemaVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current schema version: %w", err)
	}

	// Define all migrations in order
	migrations := []struct {
		version int
		name    string
		migrate func(*sql.Tx) error
	}{
		{1, "Add OTP alerts table", Migration_001_AddOTPTable},
		{2, "Add AI summaries table", Migration_002_AddAISummariesTable},
		{3, "Add digital accounts table", Migration_003_AddAccountsTable},
	}

	// Run each pending migration
	for _, m := range migrations {
		if currentVersion >= m.version {
			// Migration already applied
			continue
		}

		log.Printf("Running migration %d: %s", m.version, m.name)

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.version, err)
		}

		// Run migration using transaction
		if err := m.migrate(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d failed: %w", m.version, err)
		}

		// Update schema version
		if err := setSchemaVersion(tx, m.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update schema version for migration %d: %w", m.version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.version, err)
		}

		log.Printf("Migration %d completed successfully", m.version)
	}

	return nil
}

// Migration_001_AddOTPTable creates the otp_alerts table with indexes
// This migration is idempotent - safe to run multiple times
func Migration_001_AddOTPTable(tx *sql.Tx) error {
	schema := `
		CREATE TABLE IF NOT EXISTS otp_alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER NOT NULL,
			expires_at INTEGER NOT NULL,
			sender TEXT NOT NULL,
			subject TEXT NOT NULL,
			otp_code TEXT NOT NULL,
			confidence REAL NOT NULL,
			source TEXT NOT NULL,
			pattern_name TEXT NOT NULL,
			message_id TEXT NOT NULL,
			gmail_link TEXT NOT NULL,
			filter_name TEXT,
			is_active INTEGER DEFAULT 1,
			copied_at INTEGER
		);

		CREATE INDEX IF NOT EXISTS idx_otp_timestamp ON otp_alerts(timestamp DESC);
		CREATE INDEX IF NOT EXISTS idx_otp_expires ON otp_alerts(expires_at);
		CREATE INDEX IF NOT EXISTS idx_otp_active ON otp_alerts(is_active);
		CREATE INDEX IF NOT EXISTS idx_otp_message_id ON otp_alerts(message_id);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create otp_alerts table: %w", err)
	}

	return nil
}

// Migration_002_AddAISummariesTable creates the ai_summaries table with indexes
// This migration is idempotent - safe to run multiple times
func Migration_002_AddAISummariesTable(tx *sql.Tx) error {
	schema := `
		CREATE TABLE IF NOT EXISTS ai_summaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message_id TEXT NOT NULL UNIQUE,
			summary TEXT NOT NULL,
			questions TEXT,
			action_items TEXT,
			provider TEXT NOT NULL,
			model TEXT NOT NULL,
			generated_at INTEGER NOT NULL,
			tokens_used INTEGER DEFAULT 0
		);

		CREATE INDEX IF NOT EXISTS idx_summary_message_id ON ai_summaries(message_id);
		CREATE INDEX IF NOT EXISTS idx_summary_generated_at ON ai_summaries(generated_at DESC);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create ai_summaries table: %w", err)
	}

	return nil
}

// Migration_003_AddAccountsTable creates the accounts table for tracking digital accounts
// Supports subscriptions, trials, and free accounts
// This migration is idempotent - safe to run multiple times
func Migration_003_AddAccountsTable(tx *sql.Tx) error {
	schema := `
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service_name TEXT NOT NULL,
			email_address TEXT NOT NULL,
			account_type TEXT,
			status TEXT DEFAULT 'active',
			price_monthly REAL,
			trial_end_date INTEGER,
			gmail_message_id TEXT,
			detected_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			confidence REAL DEFAULT 0.0,
			cancel_url TEXT,
			category TEXT
		);

		CREATE INDEX IF NOT EXISTS idx_accounts_service ON accounts(service_name);
		CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email_address);
		CREATE INDEX IF NOT EXISTS idx_accounts_type ON accounts(account_type);
		CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status);
		CREATE INDEX IF NOT EXISTS idx_accounts_trial_end ON accounts(trial_end_date);
		CREATE INDEX IF NOT EXISTS idx_accounts_detected ON accounts(detected_at DESC);

		CREATE TABLE IF NOT EXISTS account_alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER NOT NULL,
			alert_type TEXT NOT NULL,
			alert_date INTEGER NOT NULL,
			sent_at INTEGER,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_account_alerts_account ON account_alerts(account_id);
		CREATE INDEX IF NOT EXISTS idx_account_alerts_date ON account_alerts(alert_date);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create accounts tables: %w", err)
	}

	return nil
}
