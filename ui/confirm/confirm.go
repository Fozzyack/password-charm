// Package confirm provides a confirmation dialog for destructive actions.
// It uses Bubble Tea for TUI functionality and maintains consistent styling with the rest of the application.
package confirm

import (
	"fmt"
	"strings"

	"github.com/Fozzyack/password-manager/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmModel represents the state of the confirmation dialog
type ConfirmModel struct {
	siteName  string
	filename  string
	action    string  // e.g., "delete", "remove"
	confirmed bool
	cancelled bool
	cursor    int     // 0 for No, 1 for Yes
	options   *types.Options
}

// Confirmation dialog styling
var (
	confirmTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF5F87")).
		Align(lipgloss.Center)

	confirmContainerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Margin(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF5F87")).
		Width(60).
		Align(lipgloss.Center)

	warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	entryInfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#626262")).
		Align(lipgloss.Center)

	buttonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#626262")).
		Padding(0, 3).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#626262"))

	selectedButtonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(0, 3).
		Margin(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF5F87")).
		Bold(true)

	confirmHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Margin(1, 0)
)

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(siteName, filename, action string, options *types.Options) ConfirmModel {
	// Clear screen for clean dialog display
	fmt.Print("\033[2J\033[H")

	return ConfirmModel{
		siteName:  siteName,
		filename:  filename,
		action:    action,
		confirmed: false,
		cancelled: false,
		cursor:    0, // Default to "No" for safety
		options:   options,
	}
}

// Init implements the tea.Model interface
func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

// Update handles user input for the confirmation dialog
func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "n", "N":
			// Cancel the action
			m.cancelled = true
			return m, tea.Quit

		case "y", "Y":
			// Confirm the action
			m.confirmed = true
			return m, tea.Quit

		case "enter", " ":
			// Confirm based on cursor position
			if m.cursor == 1 { // Yes is selected
				m.confirmed = true
			} else { // No is selected
				m.cancelled = true
			}
			return m, tea.Quit

		case "up", "k":
			// Move cursor to No
			m.cursor = 0

		case "down", "j":
			// Move cursor to Yes
			m.cursor = 1

		case "tab":
			// Toggle cursor position
			m.cursor = 1 - m.cursor
		}
	}

	return m, nil
}

// View renders the confirmation dialog interface
func (m ConfirmModel) View() string {
	var content strings.Builder

	// Title with warning color
	title := confirmTitleStyle.Render(fmt.Sprintf("⚠️  Confirm %s", strings.Title(m.action)))
	content.WriteString(title + "\n\n")

	// Dialog content
	dialogContent := ""

	// Warning message
	dialogContent += warningStyle.Render(fmt.Sprintf("Are you sure you want to %s this password entry?", m.action)) + "\n\n"
	dialogContent += warningStyle.Render("This action cannot be undone!") + "\n\n"

	// Entry information
	entryInfo := fmt.Sprintf("Site: %s\nFile: %s.gpg", m.siteName, m.filename)
	dialogContent += entryInfoStyle.Render(entryInfo) + "\n\n"

	// Buttons
	var noButton, yesButton string
	if m.cursor == 0 {
		noButton = selectedButtonStyle.Render("No")
		yesButton = buttonStyle.Render("Yes")
	} else {
		noButton = buttonStyle.Render("No")
		yesButton = selectedButtonStyle.Render("Yes")
	}

	buttons := fmt.Sprintf("    %s    %s", noButton, yesButton)
	dialogContent += buttons + "\n\n"

	content.WriteString(confirmContainerStyle.Render(dialogContent))

	// Help text
	help := confirmHelpStyle.Render("↑↓/j/k: Navigate • Enter/Space: Confirm • y: Yes • n/Esc: No")
	content.WriteString(help)

	return content.String()
}

// IsConfirmed returns whether the action was confirmed
func (m ConfirmModel) IsConfirmed() bool {
	return m.confirmed
}

// IsCancelled returns whether the action was cancelled
func (m ConfirmModel) IsCancelled() bool {
	return m.cancelled
}
