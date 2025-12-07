package ui

import "fmt"

// Demo demonstrates all UI components
func Demo() {
	// Banner
	PrintBanner("1.0.0")

	// Section
	PrintSection("Filter Management")

	// Messages
	PrintSuccess("Filter added successfully")
	PrintError("Failed to connect to Gmail API")
	PrintWarning("No filters configured")
	PrintInfo("Found 3 matching emails")

	// Subsection
	PrintSubsection("Recent Alerts")

	// Bullets
	PrintBullet("Email from boss@company.com")
	PrintBullet("Interview invitation from LinkedIn")
	PrintBullet("Payment confirmation")

	// Key-Value
	fmt.Println()
	PrintKeyValue("Status", "Running")
	PrintKeyValue("Polling Interval", "45 seconds")
	PrintKeyValue("Filters", "5")

	// Table
	fmt.Println()
	headers := []string{"Filter", "From", "Priority"}
	rows := [][]string{
		{"Job Alerts", "linkedin.com", "Normal"},
		{"Boss", "boss@company.com", "High"},
		{"Client", "client@example.com", "Normal"},
	}
	PrintTable(headers, rows)

	// Box
	PrintBox([]string{
		"Email Sentinel is running in the background.",
		"Press Ctrl+C to stop monitoring.",
	})

	// Command Example
	fmt.Println()
	PrintCommandExample("Add a new filter", "email-sentinel filter add --name \"Jobs\" --from \"linkedin.com\"")

	// Progress Bar
	fmt.Println()
	PrintProgressBar(7, 10, 30)

	// Divider
	fmt.Println()
	PrintDivider()
}
