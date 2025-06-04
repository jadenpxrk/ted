package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ted/internal/colors"
	"ted/internal/history"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// highlightCommands highlights commands in backticks with accent color
func highlightCommands(text string) string {
	// Create a style for highlighting commands
	commandHighlight := lipgloss.NewStyle().Foreground(colors.AccentColor).Bold(true)

	// Find all text within backticks and style them with AccentColor
	result := strings.Builder{}
	inBackticks := false
	currentWord := strings.Builder{}

	for _, char := range text {
		if char == '`' {
			if inBackticks {
				// End of command - render it with accent color
				commandText := currentWord.String()
				result.WriteString(commandHighlight.Render("`" + commandText + "`"))
				currentWord.Reset()
				inBackticks = false
			} else {
				// Start of command
				inBackticks = true
			}
		} else if inBackticks {
			currentWord.WriteRune(char)
		} else {
			result.WriteRune(char)
		}
	}

	// Handle case where backtick wasn't closed
	if inBackticks {
		result.WriteString("`")
		result.WriteString(currentWord.String())
	}

	return result.String()
}

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View command history",
	Long: `View your past ted command history with an interactive interface.

Use arrow keys to navigate, Enter to select an entry, and q to quit.`,
	RunE: runHistory,
}

func runHistory(cmd *cobra.Command, args []string) error {
	hist, err := history.Load()
	if err != nil {
		return fmt.Errorf("error loading history: %w", err)
	}
	defer hist.Close()

	entries, err := hist.GetEntries()
	if err != nil {
		return fmt.Errorf("error retrieving history entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("%s\n", colors.TitleStyle.Render("Ted Command History"))
		fmt.Printf("%s\n", colors.ErrorStyle.Render("No command history found."))
		fmt.Printf("%s\n", colors.QueryStyle.Render("Try running 'ted ask [question]' or 'ted agent [query]' first to build up some history."))
		return nil
	}

	fmt.Printf("%s\n", colors.TitleStyle.Render("Ted Command History"))
	fmt.Printf("%s\n", colors.HeaderStyle.Render(fmt.Sprintf("Total: %d entries", len(entries))))

	// Display all entries (max 5)
	for i, entry := range entries {
		// Format:
		// 1.
		// [query type]
		// query
		// YEAR-MONTH-DAY HOUR:MINUTE
		// `command`

		fmt.Printf("%s\n", colors.EntryStyle.Render(fmt.Sprintf("%d.", i+1)))
		fmt.Printf("%s\n", colors.CommandStyle.Render(fmt.Sprintf("[%s]", entry.Command)))
		fmt.Printf("%s\n", colors.QueryStyle.Render(entry.Query))
		fmt.Printf("%s\n", colors.TimeStyle.Render(entry.Timestamp.Format("2006-01-02 15:04")))

		// Selected
		responseText := ""
		if entry.Selected != nil {
			responseText = *entry.Selected
		} else {
			responseText = entry.Response
		}
		fmt.Printf("%s\n", colors.SelectedOptionStyle.Render(fmt.Sprintf("`%s`", responseText)))
		fmt.Println()
	}

	// Actions menu
	fmt.Printf("\n%s\n", colors.PromptStyle.Render("Actions:"))
	fmt.Printf("• Enter a number (1-%d) to view details\n", len(entries))
	fmt.Printf("• Type 'delete' to delete most recent entry\n")
	fmt.Printf("• Type 'clear' to delete all history\n")
	fmt.Printf("• Press Enter to exit\n")
	fmt.Printf("\n%s ", colors.PromptStyle.Render("Choose an action:"))

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			fmt.Printf("%s\n", colors.SuccessStyle.Render("Exited"))
			return nil
		}

		switch input {
		case "delete":
			err := hist.DeleteMostRecent()
			if err != nil {
				fmt.Printf("%s\n", colors.ErrorStyle.Render(err.Error()))
				return nil
			}
			fmt.Printf("%s\n", colors.SuccessStyle.Render("Most recent entry deleted successfully!"))

		case "clear":
			err := hist.Clear()
			if err != nil {
				fmt.Printf("%s\n", colors.ErrorStyle.Render("Error clearing history: "+err.Error()))
				return nil
			}
			fmt.Printf("%s\n", colors.SuccessStyle.Render("All history cleared successfully!"))

		default:
			num, err := strconv.Atoi(input)
			if err != nil || num < 1 || num > len(entries) {
				fmt.Printf("%s\n", colors.ErrorStyle.Render(fmt.Sprintf("Invalid selection. Please enter a number between 1 and %d.", len(entries))))
				return nil
			}

			// Show detailed view
			entry := entries[num-1]

			fmt.Printf("\n%s\n", colors.TitleStyle.Render(fmt.Sprintf("Entry %d Details", num)))
			fmt.Printf("%s\n", colors.EntryStyle.Render(fmt.Sprintf("%d.", num)))
			fmt.Printf("%s\n", colors.CommandStyle.Render(fmt.Sprintf("[%s]", entry.Command)))
			fmt.Printf("%s\n", colors.QueryStyle.Render(entry.Query))
			fmt.Printf("%s\n", colors.TimeStyle.Render(entry.Timestamp.Format("2006-01-02 15:04")))

			// Selected
			responseText := ""
			if entry.Selected != nil {
				responseText = *entry.Selected
			} else {
				responseText = entry.Response
			}
			fmt.Printf("%s\n", colors.SelectedOptionStyle.Render(fmt.Sprintf("`%s`", responseText)))

			// Show full response if it's different from selected
			if entry.Selected != nil && entry.Response != *entry.Selected {
				// Clean up the response text - handle both old format (with literal \n) and new format
				cleanResponse := entry.Response
				// Replace literal \n with actual newlines for proper formatting
				cleanResponse = strings.ReplaceAll(cleanResponse, "\\n", "\n")
				// Trim any excess whitespace
				cleanResponse = strings.TrimSpace(cleanResponse)
				// Highlight commands in backticks
				highlightedResponse := highlightCommands(cleanResponse)
				fmt.Printf("\n%s:\n%s\n",
					colors.FullResponseStyle.Render("Full Response"),
					colors.DetailBoxStyle.Render(highlightedResponse))
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
