/*
Copyright © 2025 DataTeamSix <research@dt6.io>
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

var editCmd = &cobra.Command{
	Use:   "edit [filter-name]",
	Short: "Edit an existing filter",
	Long: `Edit an existing email filter.

If no filter name is provided, you'll be shown a list to choose from.

Examples:
  email-sentinel filter edit
  email-sentinel filter edit "Job Alerts"`,
	Run: runFilterEdit,
}

func init() {
	filterCmd.AddCommand(editCmd)
}

func runFilterEdit(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	// Load existing filters
	filters, err := filter.ListFilters()
	if err != nil {
		fmt.Printf("❌ Error loading filters: %v\n", err)
		os.Exit(1)
	}

	if len(filters) == 0 {
		fmt.Println("No filters to edit.")
		fmt.Println("\nAdd one with: email-sentinel filter add")
		return
	}

	// Determine which filter to edit
	var selectedFilter *filter.Filter
	var selectedIndex int

	if len(args) > 0 {
		// Find by name
		name := args[0]
		for i, f := range filters {
			if strings.EqualFold(f.Name, name) {
				selectedFilter = &filters[i]
				selectedIndex = i
				break
			}
		}
		if selectedFilter == nil {
			fmt.Printf("❌ Filter '%s' not found\n", name)
			os.Exit(1)
		}
	} else {
		// Interactive selection
		fmt.Println("\n📋 Select a filter to edit:")
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

		selectedIndex = num - 1
		selectedFilter = &filters[selectedIndex]
	}

	fmt.Printf("\n✏️  Editing: %s\n", selectedFilter.Name)
	fmt.Println(strings.Repeat("━", 40))
	fmt.Println("Press Enter to keep current value, or type new value.")

	// Edit name
	fmt.Printf("\nName [%s]: ", selectedFilter.Name)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		selectedFilter.Name = input
	}

	// Edit from patterns
	currentFrom := strings.Join(selectedFilter.From, ", ")
	if currentFrom == "" {
		currentFrom = "(none)"
	}
	fmt.Printf("\nFrom contains [%s]: ", currentFrom)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		if input == "-" || input == "none" {
			selectedFilter.From = []string{}
		} else {
			selectedFilter.From = parseCSV(input)
		}
	}

	// Edit subject patterns
	currentSubject := strings.Join(selectedFilter.Subject, ", ")
	if currentSubject == "" {
		currentSubject = "(none)"
	}
	fmt.Printf("\nSubject contains [%s]: ", currentSubject)
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		if input == "-" || input == "none" {
			selectedFilter.Subject = []string{}
		} else {
			selectedFilter.Subject = parseCSV(input)
		}
	}

	// Validate at least one pattern
	if len(selectedFilter.From) == 0 && len(selectedFilter.Subject) == 0 {
		fmt.Println("\n❌ At least one 'from' or 'subject' pattern is required")
		os.Exit(1)
	}

	// Edit match mode (only if both from and subject exist)
	if len(selectedFilter.From) > 0 && len(selectedFilter.Subject) > 0 {
		fmt.Printf("\nMatch mode - 'any' (OR) or 'all' (AND) [%s]: ", selectedFilter.Match)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "any" || input == "all" || input == "and" || input == "or" {
			if input == "and" {
				input = "all"
			} else if input == "or" {
				input = "any"
			}
			selectedFilter.Match = input
		}
	}

	// Update the filter
	if err := filter.UpdateFilter(selectedIndex, *selectedFilter); err != nil {
		fmt.Printf("\n❌ Error updating filter: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Filter updated successfully!")
	fmt.Println()
	printFilter(*selectedFilter)
}
