package colors

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	// Primary
	PrimaryColor   = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor = lipgloss.Color("#10B981") // Green
	AccentColor    = lipgloss.Color("#F59E0B") // Amber
	ErrorColor     = lipgloss.Color("#EF4444") // Red
	MutedColor     = lipgloss.Color("#6B7280") // Gray

	// Additional
	CyanColor   = lipgloss.Color("14")      // Cyan
	BrightGreen = lipgloss.Color("10")      // Bright green
	LightGray   = lipgloss.Color("#E5E7EB") // Light gray
	LightPurple = lipgloss.Color("#A78BFA") // Light purple
)

// Styled components for consistent usage across commands
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	EntryStyle = lipgloss.NewStyle().
			Foreground(LightGray)

	CommandStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	QueryStyle = lipgloss.NewStyle().
			Foreground(LightGray).
			Italic(true)

	TimeStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	SelectedOptionStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	DetailBoxStyle = lipgloss.NewStyle().
			Foreground(LightGray)

	ThinkingStyle = lipgloss.NewStyle().
			Foreground(CyanColor).
			Bold(true)

	RunningStyle = lipgloss.NewStyle().
			Foreground(BrightGreen).
			Bold(true)

	FullResponseStyle = lipgloss.NewStyle().
				Foreground(LightPurple).
				Bold(true)

	SettingsLabelStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true)

	SettingsValueStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor)

	SettingsConfiguredStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true)

	SettingsNotSetStyle = lipgloss.NewStyle().
				Foreground(ErrorColor).
				Bold(true)

	SettingsOptionStyle = lipgloss.NewStyle().
				Foreground(LightGray)

	SettingsWarningStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true)

	SettingsInfoStyle = lipgloss.NewStyle().
				Foreground(CyanColor)
)
