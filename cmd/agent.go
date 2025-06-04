package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ted/internal/colors"
	"ted/internal/config"
	"ted/internal/gemini"
	"ted/internal/history"
	"ted/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent [query]",
	Short: "Generate a command from natural language",
	Long: `Generate a command from natural language using AI.

Example:
  ted agent how to make a python virtual environment
  ted agent list all files in current directory
  ted agent compress a folder into a zip file`,
	RunE: runAgent,
}

func runAgent(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a query. Example: ted agent how to make a python3 venv")
	}

	query := strings.Join(args, " ")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if cfg.GeminiAPIKey == "" {
		return fmt.Errorf("gemini API key not configured. Run 'ted settings' to set it up")
	}

	client, err := gemini.NewClient(cfg.GeminiAPIKey, cfg.Model, cfg.Temperature)
	if err != nil {
		return fmt.Errorf("error creating Gemini client: %w", err)
	}
	defer client.Close()

	fmt.Printf("%s\n", colors.ThinkingStyle.Render("Thinking..."))

	ctx := context.Background()
	response, err := client.GenerateAgentCommand(ctx, query)
	if err != nil {
		return fmt.Errorf("error generating command: %w", err)
	}

	confirmModel := ui.NewConfirmModel(response.Command, response.Explanation)
	p := tea.NewProgram(confirmModel)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running UI: %w", err)
	}

	confirm := finalModel.(ui.ConfirmModel)
	if confirm.WasCancelled() {
		fmt.Println("Command execution cancelled.")
		return nil
	}

	if confirm.ShouldExecute() {
		if err := executeCommand(response.Command); err != nil {
			return err
		}

		if err := saveToHistory(query, response); err != nil {
			fmt.Printf("Warning: Failed to save to history: %v\n", err)
		}
	}

	return nil
}

func executeCommand(command string) error {
	fmt.Printf("%s\n", colors.RunningStyle.Render(fmt.Sprintf("Running `%s`", command)))

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func saveToHistory(query string, response *gemini.AgentResponse) error {
	hist, err := history.Load()
	if err != nil {
		return err
	}
	defer hist.Close()

	selected := response.Command
	return hist.AddEntry("agent", query, response.Explanation, &selected)
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
