/*
Copyright © 2025 Datateamsix <research@dt6.io>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/filter"
)

var removeCmd = &cobra.Command{
	Use:   "remove [filter-name]",
	Short: "Remove a filter",
	Long: `Remove an email filter by name.

If no name is provided, you'll be shown a list to choose from.

Examples:
  email-sentinel filter remove
  email-sentinel filter remove "Job Alerts"`,
	Run: runFilterRemove,
}

func init() {
	filterCmd.AddCommand(removeCmd)
}

func runFilterRemove(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var filterName string

	if len(args) > 0 {
		filterName = args[0]
	} else {
		// Interactive selection
		filters, err := filter.ListFilters()
		if err != nil {
			fmt.Printf("❌ Error loading filters: %v\n", err)
			os.Exit(1)
		}

		if len(filters) == 0 {
			fmt.Println("No filters to remove.")
			return
		}

		fmt.Println("\n🗑️  Select a filter to remove:")
		fmt.Println(strings.Repeat("━", 40))

		for i, f := range filters {
			fmt.Printf("[%d] %s\n", i+1, f.Name)
		}

		fmt.Print("\nEnter number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(filters) {
			fmt.Println("❌ Invalid selection")
			os.Exit(1)
		}

		filterName = filters[num-1].Name
	}

	// Confirm deletion
	fmt.Printf("\nRemove filter '%s'? (y/N): ", filterName)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("Cancelled.")
		return
	}

	if err := filter.RemoveFilter(filterName); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Filter '%s' removed.\n", filterName)
}
