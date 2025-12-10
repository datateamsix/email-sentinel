/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/datateamsix/email-sentinel/internal/storage"
	"github.com/spf13/cobra"
)

// dbBackupCmd represents the db backup command
var dbBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of the database",
	Long: `Creates a backup of the Email Sentinel database.

The backup is created using SQLite's VACUUM INTO command, which produces
a clean, defragmented copy of the database.

Backups are stored in:
  - Windows: %APPDATA%\email-sentinel\backups\
  - macOS: ~/Library/Application Support/email-sentinel/backups/
  - Linux: ~/.config/email-sentinel/backups/

The application automatically keeps the last 5 backups and removes older ones.

Example:
  email-sentinel db backup`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📦 Creating database backup...")

		// Initialize database connection
		db, err := storage.InitDB()
		if err != nil {
			fmt.Printf("❌ Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer storage.CloseDB(db)

		// Create backup
		if err := storage.BackupDatabase(db); err != nil {
			fmt.Printf("❌ Backup failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Backup completed successfully!")
		fmt.Println("💡 Tip: Backups are created automatically on application startup")
	},
}

func init() {
	dbCmd.AddCommand(dbBackupCmd)
}
