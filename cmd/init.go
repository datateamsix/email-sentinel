/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/config"
	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/gmail"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize email-sentinel with Gmail authentication",
	Long: `Initialize email-sentinel by authenticating with Gmail.
	
This command will:
1. Read your credentials.json file
2. Open a browser for Google OAuth authorization
3. Save your authentication token for future use

You must have a credentials.json file from Google Cloud Console.`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Initializing email-sentinel...")

	// Check if already initialized
	if gmail.TokenExists() {
		fmt.Println("\n⚠️  Already initialized! Token exists.")
		fmt.Print("Do you want to re-authenticate? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Keeping existing authentication.")
			showPostInitMenu()
			return
		}
	}

	// Look for credentials.json
	credPath := findCredentials()
	if credPath == "" {
		fmt.Println("\n❌ Error: credentials.json not found")
		fmt.Println("\nPlease ensure credentials.json is in one of these locations:")
		fmt.Println("  - Current directory: ./credentials.json")
		configDir, _ := config.ConfigDir()
		fmt.Printf("  - Config directory: %s/credentials.json\n", configDir)
		os.Exit(1)
	}
	fmt.Printf("✓ Found credentials: %s\n", credPath)

	// Load OAuth config
	oauthConfig, err := gmail.LoadCredentials(credPath)
	if err != nil {
		fmt.Printf("\n❌ Error loading credentials: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Credentials loaded")

	// Run OAuth flow
	token, err := gmail.GetTokenFromWeb(oauthConfig)
	if err != nil {
		fmt.Printf("\n❌ Error during authentication: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Authentication successful")

	// Save token
	if err := gmail.SaveToken(token); err != nil {
		fmt.Printf("\n❌ Error saving token: %v\n", err)
		os.Exit(1)
	}

	tokenPath, _ := config.TokenPath()
	fmt.Printf("✓ Token saved to: %s\n", tokenPath)

	fmt.Println("\n✅ Initialization complete!")
	showPostInitMenu()
}

func showPostInitMenu() {
	// Show existing filters if any
	filters, err := filter.ListFilters()
	if err == nil && len(filters) > 0 {
		fmt.Printf("\n📋 Existing Filters (%d)\n", len(filters))
		fmt.Println(strings.Repeat("━", 40))
		for i, f := range filters {
			fmt.Printf("[%d] %s\n", i+1, f.Name)
			if len(f.From) > 0 {
				fmt.Printf("    From:    %s\n", strings.Join(f.From, ", "))
			}
			if len(f.Subject) > 0 {
				fmt.Printf("    Subject: %s\n", strings.Join(f.Subject, ", "))
			}
			fmt.Printf("    Match:   %s\n", f.Match)
		}
	}

	fmt.Println("\n" + strings.Repeat("━", 40))
	fmt.Println("What would you like to do?")
	fmt.Println()
	fmt.Println("  1. Add a new filter      email-sentinel filter add")
	fmt.Println("  2. Edit a filter         email-sentinel filter edit")
	fmt.Println("  3. List all filters      email-sentinel filter list")
	fmt.Println("  4. Remove a filter       email-sentinel filter remove")
	fmt.Println("  5. Start watching        email-sentinel start")
	fmt.Println()
}

// findCredentials looks for credentials.json in common locations
func findCredentials() string {
	locations := []string{
		"credentials.json",
		"./credentials.json",
	}

	// Also check config directory
	if configDir, err := config.ConfigDir(); err == nil {
		locations = append(locations, filepath.Join(configDir, "credentials.json"))
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	return ""
}
