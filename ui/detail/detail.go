// Package detail provides a detailed view for displaying individual password entries.
// It uses Bubble Tea for TUI functionality and maintains consistent styling with the rest of the application.
package detail

import (
	"fmt"
	"strings"

	"github.com/Fozzyack/password-manager/encryption"
	"github.com/Fozzyack/password-manager/types"
	"github.com/Fozzyack/password-manager/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModel represents the state of the password detail view
type DetailModel struct {
	entry           encryption.Data
	filename        string
	siteName        string
	showPassword    bool
	deleteRequested bool
	options         *types.Options
}

// Detail view styling
var (
	detailTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)

	detailContainerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Margin(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(70).
		Align(lipgloss.Left)

	fieldLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Width(15).
		Align(lipgloss.Right)

	fieldValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	passwordHiddenStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Padding(0, 1)

	passwordVisibleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Bold(true)

	strengthStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true)

	detailHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		PaddingLeft(4).
		Italic(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	timestampStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Padding(0, 1)
)

// NewPasswordDetail creates a new password detail view
func NewPasswordDetail(entry encryption.Data, filename, siteName string, options *types.Options) DetailModel {
	// Clear screen for clean detail display
	fmt.Print("\033[2J\033[H")

	return DetailModel{
		entry:           entry,
		filename:        filename,
		siteName:        siteName,
		showPassword:    false,
		deleteRequested: false,
		options:         options,
	}
}

// Init implements the tea.Model interface
func (m DetailModel) Init() tea.Cmd {
	return nil
}

// Update handles user input for the detail view
func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q", "backspace":
			// Return to password list
			return m, tea.Quit

		case "v", " ":
			// Toggle password visibility
			m.showPassword = !m.showPassword

		case "d", "D":
			// Request deletion
			m.deleteRequested = true
			return m, tea.Quit

		case "enter":
			// Return to list (same as escape)
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the password detail interface
func (m DetailModel) View() string {
	var content strings.Builder

	// Title
	title := detailTitleStyle.Render("🔍 Password Details")
	content.WriteString(title + "\n\n")

	// Detail content
	detailContent := ""

	// Site/Service Name
	detailContent += fieldLabelStyle.Render("Site/Service:") + 
		fieldValueStyle.Render(m.siteName) + "\n\n"

	// Username
	if m.entry.Username != "" {
		detailContent += fieldLabelStyle.Render("Username:") + 
			fieldValueStyle.Render(m.entry.Username) + "\n\n"
	}

	// Email
	if m.entry.Email != "" {
		detailContent += fieldLabelStyle.Render("Email:") + 
			fieldValueStyle.Render(m.entry.Email) + "\n\n"
	}

	// URL
	if m.entry.URL != "" {
		detailContent += fieldLabelStyle.Render("URL:") + 
			fieldValueStyle.Render(m.entry.URL) + "\n\n"
	}

	// Password
	passwordLabel := fieldLabelStyle.Render("Password:")
	if m.showPassword {
		passwordValue := passwordVisibleStyle.Render(m.entry.Password)
		detailContent += passwordLabel + passwordValue + "\n"
		
		// Show password strength
		strength, description := utils.EvaluatePasswordStrength(m.entry.Password)
		var strengthColor string
		switch strength {
		case 0, 1:
			strengthColor = "#FF5F87" // Red
		case 2:
			strengthColor = "#FFD700" // Yellow
		case 3:
			strengthColor = "#87CEEB" // Light Blue
		case 4:
			strengthColor = "#90EE90" // Light Green
		}
		
		strengthText := strengthStyle.Copy().
			Foreground(lipgloss.Color(strengthColor)).
			Render(fmt.Sprintf("Strength: %s", description))
		detailContent += fieldLabelStyle.Render("") + strengthText + "\n\n"
	} else {
		passwordValue := passwordHiddenStyle.Render("••••••••••••••••")
		detailContent += passwordLabel + passwordValue + "\n"
		detailContent += fieldLabelStyle.Render("") + 
			passwordHiddenStyle.Render("Press 'v' or Space to reveal password") + "\n\n"
	}

	// File information
	detailContent += "─" + strings.Repeat("─", 60) + "\n\n"
	
	detailContent += fieldLabelStyle.Render("Filename:") + 
		fieldValueStyle.Render(m.filename + ".gpg") + "\n\n"

	// Timestamps
	detailContent += fieldLabelStyle.Render("Created:") + 
		timestampStyle.Render(m.entry.CreatedAt.Format("Monday, January 2, 2006 at 3:04 PM")) + "\n\n"

	if !m.entry.UpdatedAt.Equal(m.entry.CreatedAt) {
		detailContent += fieldLabelStyle.Render("Updated:") + 
			timestampStyle.Render(m.entry.UpdatedAt.Format("Monday, January 2, 2006 at 3:04 PM")) + "\n\n"
	}

	content.WriteString(detailContainerStyle.Render(detailContent))

	// Help text
	var helpText string
	if m.showPassword {
		helpText = "v/Space: Hide Password • d: Delete • Esc/q/Backspace: Back to List • Enter: Back to List"
	} else {
		helpText = "v/Space: Show Password • d: Delete • Esc/q/Backspace: Back to List • Enter: Back to List"
	}
	
	help := detailHelpStyle.Render(helpText)
	content.WriteString(help)

	return content.String()
}

// IsPasswordVisible returns whether the password is currently visible
func (m DetailModel) IsPasswordVisible() bool {
	return m.showPassword
}

// IsDeletionRequested returns whether deletion was requested for this entry
func (m DetailModel) IsDeletionRequested() bool {
	return m.deleteRequested
}
