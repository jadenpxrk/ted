package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmModel struct {
	command     string
	explanation string
	confirmed   bool
	cancelled   bool
}

func NewConfirmModel(command, explanation string) ConfirmModel {
	return ConfirmModel{
		command:     command,
		explanation: explanation,
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.confirmed = true
			return m, tea.Quit
		case "n", "N", "q", "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	if m.confirmed {
		return ""
	}
	if m.cancelled {
		return "Cancelled.\n"
	}

	// Simple styling without borders
	commandStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))

	content := fmt.Sprintf("%s\nCommand: %s\n%s",
		m.explanation,
		commandStyle.Render(m.command),
		promptStyle.Render("Execute this command? (y/N): "))

	return content
}

func (m ConfirmModel) ShouldExecute() bool {
	return m.confirmed
}

func (m ConfirmModel) WasCancelled() bool {
	return m.cancelled
}
