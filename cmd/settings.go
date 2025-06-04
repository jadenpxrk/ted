package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ted/internal/colors"
	"ted/internal/config"

	"github.com/spf13/cobra"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Configure ted settings",
	Long: `Configure your ted CLI settings including API keys and model preferences.

This will guide you through setting up:
- Gemini API key
- AI model selection
- Temperature setting for AI response determinism`,
	RunE: runSettings,
}

func runSettings(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s ", colors.SettingsLabelStyle.Render("Current Gemini API key:"))
	if cfg.GeminiAPIKey != "" {
		fmt.Printf("%s\n", colors.SettingsConfiguredStyle.Render("[CONFIGURED]"))
	} else {
		fmt.Printf("%s\n", colors.SettingsNotSetStyle.Render("[NOT SET]"))
	}
	fmt.Printf("%s ", colors.PromptStyle.Render("Enter new Gemini API key (or press Enter to keep current):"))

	if scanner.Scan() {
		apiKey := strings.TrimSpace(scanner.Text())
		if apiKey != "" {
			cfg.GeminiAPIKey = apiKey
			fmt.Printf("%s\n", colors.SuccessStyle.Render("✓ Gemini API key updated"))
		}
	}

	fmt.Printf("\n%s %s\n", colors.SettingsLabelStyle.Render("Current model:"), colors.SettingsValueStyle.Render(cfg.Model))
	fmt.Printf("%s\n", colors.HeaderStyle.Render("Available models:"))
	models := []string{
		"gemini-2.0-flash",
		"gemini-2.0-flash-lite",
		"gemini-2.5-pro-preview-05-06",
		"gemini-2.5-flash-preview-05-20",
	}

	for i, model := range models {
		fmt.Printf("  %s", colors.SettingsOptionStyle.Render(fmt.Sprintf("%d. %s", i+1, model)))
		if model == "gemini-2.0-flash" {
			fmt.Printf(" %s", colors.SettingsConfiguredStyle.Render("(default)"))
		}
		fmt.Println()
	}
	fmt.Printf("%s ", colors.PromptStyle.Render("Enter number (1-4) or press Enter to keep current:"))

	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			if num, err := strconv.Atoi(input); err == nil && num >= 1 && num <= len(models) {
				cfg.Model = models[num-1]
				fmt.Printf("%s\n", colors.SuccessStyle.Render("✓ Model updated"))
			} else {
				fmt.Printf("%s\n", colors.SettingsWarningStyle.Render(fmt.Sprintf("⚠️  Invalid selection '%s'. Please enter a number between 1 and %d.", input, len(models))))
			}
		}
	}

	fmt.Printf("\n%s %.2f\n", colors.SettingsLabelStyle.Render("Current temperature:"), cfg.Temperature)
	fmt.Printf("%s\n", colors.HeaderStyle.Render("Temperature affects AI response determinism:"))
	fmt.Printf("  %s\n", colors.SettingsInfoStyle.Render("• Lower values (0.0-0.3): More deterministic, consistent responses"))
	fmt.Printf("  %s\n", colors.SettingsInfoStyle.Render("• Higher values (0.4-1.0): More creative, varied responses"))
	fmt.Printf("  %s\n", colors.SettingsInfoStyle.Render("• Default: 0.3 (recommended for command generation)"))
	fmt.Printf("%s ", colors.PromptStyle.Render("Enter temperature (0.0-1.0) or press Enter to keep current:"))

	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			if temp, err := strconv.ParseFloat(input, 32); err == nil && temp >= 0.0 && temp <= 1.0 {
				cfg.Temperature = float32(temp)
				fmt.Printf("%s\n", colors.SuccessStyle.Render("✓ Temperature updated"))
			} else {
				fmt.Printf("%s\n", colors.SettingsWarningStyle.Render("⚠️  Invalid temperature. Please enter a value between 0.0 and 1.0."))
			}
		}
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	fmt.Printf("\n%s\n", colors.SuccessStyle.Render("✓ Settings saved successfully!"))
	fmt.Printf("%s %s\n", colors.SettingsInfoStyle.Render("Config location:"), colors.SettingsValueStyle.Render(config.GetConfigPath()))

	fmt.Printf("\n%s\n", colors.HeaderStyle.Render("Current configuration:"))
	fmt.Printf("  %s ", colors.SettingsLabelStyle.Render("Gemini API Key:"))
	if cfg.GeminiAPIKey != "" {
		fmt.Printf("%s\n", colors.SettingsConfiguredStyle.Render("[CONFIGURED]"))
	} else {
		fmt.Printf("%s\n", colors.SettingsNotSetStyle.Render("[NOT SET]"))
	}
	fmt.Printf("  %s %s\n", colors.SettingsLabelStyle.Render("Model:"), colors.SettingsValueStyle.Render(cfg.Model))
	fmt.Printf("  %s %.2f\n", colors.SettingsLabelStyle.Render("Temperature:"), cfg.Temperature)

	if cfg.GeminiAPIKey == "" {
		fmt.Printf("\n%s\n", colors.SettingsWarningStyle.Render("⚠️  Warning: Gemini API key is not set. You'll need to configure it to use ted."))
		fmt.Printf("%s\n", colors.SettingsInfoStyle.Render("Get your API key from: https://makersuite.google.com/app/apikey"))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
