/*
Copyright © 2025 DATATEAMSIX <research@dt6.io>
*/
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/datateamsix/email-sentinel/internal/filter"
	"github.com/datateamsix/email-sentinel/internal/storage"
)

var (
	filterName    string
	filterFrom    string
	filterSubject string
	filterMatch   string
	filterLabels  string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new email filter",
	Long: `Add a new filter to match incoming emails.

You can filter by sender (from), subject line keywords, or both.
When using both, choose whether ALL conditions must match (AND) 
or ANY condition triggers a match (OR).

Examples:
  # Interactive mode
  email-sentinel filter add

  # Filter by sender only
  email-sentinel filter add --name "From Boss" --from "boss@company.com"

  # Filter by subject keywords
  email-sentinel filter add --name "Urgent" --subject "urgent,asap,important"

  # Both sender AND subject must match
  email-sentinel filter add --name "Job Alerts" --from "linkedin.com" --subject "interview" --match all

  # Either sender OR subject matches (default)
  email-sentinel filter add --name "Recruiter" --from "greenhouse.io,lever.co" --subject "opportunity" --match any`,
	Run: runFilterAdd,
}

func init() {
	filterCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&filterName, "name", "n", "", "Filter name")
	addCmd.Flags().StringVarP(&filterFrom, "from", "f", "", "Sender patterns (comma-separated)")
	addCmd.Flags().StringVarP(&filterSubject, "subject", "s", "", "Subject patterns (comma-separated)")
	addCmd.Flags().StringVarP(&filterMatch, "match", "m", "any", "Match mode: 'any' (OR) or 'all' (AND)")
	addCmd.Flags().StringVarP(&filterLabels, "labels", "l", "", "Labels/categories (comma-separated, e.g., work,urgent)")
}

func runFilterAdd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	interactive := !cmd.Flags().Changed("name")

	if interactive {
		fmt.Println("\n📧 Add New Email Filter")
		fmt.Println(strings.Repeat("━", 40))
	}

	// Get name (required)
	if filterName == "" {
		fmt.Print("\nFilter name: ")
		filterName, _ = reader.ReadString('\n')
		filterName = strings.TrimSpace(filterName)
	}

	if filterName == "" {
		fmt.Println("❌ Filter name is required")
		os.Exit(1)
	}

	// Get from patterns
	if !cmd.Flags().Changed("from") && interactive {
		fmt.Println("\n📤 Sender Filter (From)")
		fmt.Println("   Match emails from specific senders.")
		fmt.Println("   Examples: boss@company.com, @linkedin.com, greenhouse.io")
		fmt.Print("\nFrom contains (comma-separated, or blank to skip): ")
		filterFrom, _ = reader.ReadString('\n')
		filterFrom = strings.TrimSpace(filterFrom)
	}

	// Get subject patterns
	if !cmd.Flags().Changed("subject") && interactive {
		fmt.Println("\n📝 Subject Filter")
		fmt.Println("   Match emails with specific words in subject line.")
		fmt.Println("   Examples: interview, urgent, invoice")
		fmt.Print("\nSubject contains (comma-separated, or blank to skip): ")
		filterSubject, _ = reader.ReadString('\n')
		filterSubject = strings.TrimSpace(filterSubject)
	}

	// Validate at least one pattern
	if filterFrom == "" && filterSubject == "" {
		fmt.Println("\n❌ At least one 'from' or 'subject' pattern is required")
		os.Exit(1)
	}

	// Parse comma-separated values
	fromPatterns := parseCSV(filterFrom)
	subjectPatterns := parseCSV(filterSubject)

	// Get match mode (only ask if both from and subject are specified)
	if !cmd.Flags().Changed("match") && len(fromPatterns) > 0 && len(subjectPatterns) > 0 && interactive {
		fmt.Println("\n🔀 Match Mode")
		fmt.Println("   You specified both sender and subject filters.")
		fmt.Println()
		fmt.Println("   ANY (OR): Notify if sender matches OR subject matches")
		fmt.Println("             → More notifications, broader matching")
		fmt.Println()
		fmt.Println("   ALL (AND): Notify only if sender AND subject both match")
		fmt.Println("              → Fewer notifications, precise matching")
		fmt.Print("\nMatch mode [any/all] (default: any): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "all" || input == "and" {
			filterMatch = "all"
		} else {
			filterMatch = "any"
		}
	}

	// Get labels/categories
	if !cmd.Flags().Changed("labels") && interactive {
		// Try to load existing labels from database
		db, _ := getDB()
		var existingLabels []string
		if db != nil {
			existingLabels, _ = getExistingLabels(db)
			db.Close()
		}

		fmt.Println("\n🏷️  Labels/Categories (Optional)")
		fmt.Println("   Organize filters by category (e.g., work, personal, urgent)")

		if len(existingLabels) > 0 {
			fmt.Printf("   Existing labels: %s\n", strings.Join(existingLabels, ", "))
		}

		fmt.Print("\nLabels (comma-separated, or blank to skip): ")
		filterLabels, _ = reader.ReadString('\n')
		filterLabels = strings.TrimSpace(filterLabels)
	}

	// Parse labels
	labelsList := parseCSV(filterLabels)

	// Create filter
	f := filter.Filter{
		Name:    filterName,
		From:    fromPatterns,
		Subject: subjectPatterns,
		Match:   filterMatch,
		Labels:  labelsList,
	}

	// Save filter
	if err := filter.AddFilter(f); err != nil {
		fmt.Printf("\n❌ Error adding filter: %v\n", err)
		os.Exit(1)
	}

	// Save labels to database for reuse
	if len(labelsList) > 0 {
		db, err := getDB()
		if err == nil && db != nil {
			saveLabelsToDatabase(db, labelsList)
			db.Close()
		}
	}

	fmt.Println("\n✅ Filter added successfully!")
	fmt.Println()
	printFilter(f)

	// Reset flags for next use
	filterName = ""
	filterFrom = ""
	filterSubject = ""
	filterMatch = "any"
	filterLabels = ""
}

func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func printFilter(f filter.Filter) {
	fmt.Printf("  Name:    %s\n", f.Name)
	if len(f.From) > 0 {
		fmt.Printf("  From:    %s\n", strings.Join(f.From, ", "))
	}
	if len(f.Subject) > 0 {
		fmt.Printf("  Subject: %s\n", strings.Join(f.Subject, ", "))
	}
	if len(f.Labels) > 0 {
		fmt.Printf("  Labels:  %s\n", strings.Join(f.Labels, ", "))
	}

	matchDesc := "any (OR - either condition triggers)"
	if f.Match == "all" {
		matchDesc = "all (AND - all conditions must match)"
	}
	fmt.Printf("  Match:   %s\n", matchDesc)
}

// getDB initializes and returns a database connection
func getDB() (*sql.DB, error) {
	return storage.InitDB()
}

// getExistingLabels retrieves all existing labels from the database
func getExistingLabels(db *sql.DB) ([]string, error) {
	return storage.GetAllLabels(db)
}

// saveLabelsToDatabase saves labels to the database for reuse
func saveLabelsToDatabase(db *sql.DB, labels []string) {
	storage.SaveLabels(db, labels)
}
