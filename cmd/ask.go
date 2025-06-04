package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"ted/internal/colors"
	"ted/internal/config"
	"ted/internal/gemini"
	"ted/internal/history"

	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Get multiple command suggestions for a question",
	Long: `Get top 3 command suggestions for a question.

Example:
  ted ask how to make a python virtual environment
  ted ask how to find large files
  ted ask how to check disk usage`,
	RunE: runAsk,
}

func runAsk(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a question. Example: ted ask how to make a python venv")
	}

	question := strings.Join(args, " ")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if cfg.GeminiAPIKey == "" {
		return fmt.Errorf("Gemini API key not configured. Run 'ted settings' to set it up")
	}

	client, err := gemini.NewClient(cfg.GeminiAPIKey, cfg.Model, cfg.Temperature)
	if err != nil {
		return fmt.Errorf("error creating Gemini client: %w", err)
	}
	defer client.Close()

	fmt.Printf("%s\n\n", colors.ThinkingStyle.Render("Thinking..."))

	ctx := context.Background()
	response, err := client.GenerateAskCommands(ctx, question)
	if err != nil {
		return fmt.Errorf("error generating commands: %w", err)
	}

	for i, option := range response.Commands {
		coloredCommand := colors.CommandStyle.Render(fmt.Sprintf("`%s`", option.Command))
		fmt.Printf("%d. %s - %s\n", i+1, coloredCommand, option.Description)
	}

	fmt.Print("\nSelect an option (1-3) or press Enter to exit: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return nil
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(response.Commands) {
		return fmt.Errorf("invalid selection")
	}

	selectedCommand := response.Commands[choice-1].Command

	fmt.Printf("%s\n", colors.RunningStyle.Render(fmt.Sprintf("Running `%s`", selectedCommand)))

	execCmd := exec.Command("sh", "-c", selectedCommand)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	execCmdErr := execCmd.Run()

	hist, err := history.Load()
	if err == nil {
		responseText := ""
		for i, option := range response.Commands {
			if i > 0 {
				responseText += "\n"
			}
			responseText += fmt.Sprintf("%d. `%s` - %s", i+1, option.Command, option.Description)
		}
		if err := hist.AddEntry("ask", question, responseText, &selectedCommand); err != nil {
			fmt.Printf("Warning: Failed to save to history: %v\n", err)
		}
		hist.Close()
	}

	if execCmdErr != nil {
		return fmt.Errorf("command failed: %w", execCmdErr)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(askCmd)
}
