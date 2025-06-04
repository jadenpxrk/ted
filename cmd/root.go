package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Version information
const Version = "v0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ted",
	Version: Version,
	Short:   "AI-powered command line assistant",
	Long: `Ted is the fastest way to get answers in the terminal.

Available commands:
  agent    - Generate a single command from natural language and optionally execute it
  ask      - Get multiple command suggestions for a question  
  history  - View your command history with an interactive interface
  settings - Configure API keys and preferences
  version  - Show version information

Examples:
  ted agent how to make a python virtual environment
  ted ask how to find large files
  ted history
  ted settings
  ted version`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
