// Package change provides a master password change form for the password manager.
// It uses Bubble Tea for TUI functionality and maintains consistent styling with the rest of the application.
package change

import (
	"fmt"
	"strings"

	"github.com/Fozzyack/password-manager/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChangePasswordField represents the different input fields
type ChangePasswordField int

const (
	CurrentPasswordField ChangePasswordField = iota
	NewPasswordField
	ConfirmPasswordField
)

// ChangeModel represents the state of the password change form
type ChangeModel struct {
	inputs        []textinput.Model
	currentField  int
	submitted     bool
	cancelled     bool
	options       *types.Options
	currentPass   string
	newPass       string
	confirmPass   string
}

// Form styling
var (
	changeTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)

	changeContainerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Margin(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(70).
		Align(lipgloss.Left)

	changeFieldLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Margin(0, 0, 0, 1)

	changeRequiredStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true)

	changeHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	changeErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true).
		Align(lipgloss.Left).
		Margin(0, 0, 1, 1)

	changeSuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#90EE90")).
		Bold(true).
		Align(lipgloss.Center).
		Margin(1, 0)
)

// NewChangePasswordForm creates a new master password change form
func NewChangePasswordForm(options *types.Options) ChangeModel {
	// Clear screen for clean form display
	fmt.Print("\033[2J\033[H")

	inputs := make([]textinput.Model, 3)

	// Current password field
	inputs[CurrentPasswordField] = textinput.New()
	inputs[CurrentPasswordField].Placeholder = "Enter current master password"
	inputs[CurrentPasswordField].EchoMode = textinput.EchoPassword
	inputs[CurrentPasswordField].EchoCharacter = '•'
	inputs[CurrentPasswordField].CharLimit = 200
	inputs[CurrentPasswordField].Width = 50
	inputs[CurrentPasswordField].Focus()

	// New password field
	inputs[NewPasswordField] = textinput.New()
	inputs[NewPasswordField].Placeholder = "Enter new master password (8+ chars)"
	inputs[NewPasswordField].EchoMode = textinput.EchoPassword
	inputs[NewPasswordField].EchoCharacter = '•'
	inputs[NewPasswordField].CharLimit = 200
	inputs[NewPasswordField].Width = 50

	// Confirm password field
	inputs[ConfirmPasswordField] = textinput.New()
	inputs[ConfirmPasswordField].Placeholder = "Confirm new master password"
	inputs[ConfirmPasswordField].EchoMode = textinput.EchoPassword
	inputs[ConfirmPasswordField].EchoCharacter = '•'
	inputs[ConfirmPasswordField].CharLimit = 200
	inputs[ConfirmPasswordField].Width = 50

	// Style all inputs
	for i := range inputs {
		inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
		inputs[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		inputs[i].PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true)
	}

	return ChangeModel{
		inputs:       inputs,
		currentField: 0,
		submitted:    false,
		cancelled:    false,
		options:      options,
	}
}

// Init implements the tea.Model interface
func (m ChangeModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles user input and form navigation
func (m ChangeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			m.options.Quit = false // Don't quit the entire app, just cancel the form
			return m, tea.Quit

		case "enter":
			// Move to next field or submit if on last field
			if m.currentField < len(m.inputs)-1 {
				m.inputs[m.currentField].Blur()
				m.currentField++
				m.inputs[m.currentField].Focus()
				return m, m.inputs[m.currentField].Cursor.BlinkCmd()
			} else {
				// Validate and submit
				if m.validateForm() {
					// Store values
					m.currentPass = m.inputs[CurrentPasswordField].Value()
					m.newPass = m.inputs[NewPasswordField].Value()
					m.confirmPass = m.inputs[ConfirmPasswordField].Value()
					m.submitted = true
					return m, tea.Quit
				}
				// If validation fails, stay on current field
				return m, nil
			}

		case "tab", "shift+tab", "up", "down":
			// Navigate between fields
			if msg.String() == "up" || msg.String() == "shift+tab" {
				if m.currentField > 0 {
					m.inputs[m.currentField].Blur()
					m.currentField--
					m.inputs[m.currentField].Focus()
					return m, m.inputs[m.currentField].Cursor.BlinkCmd()
				}
			} else {
				if m.currentField < len(m.inputs)-1 {
					m.inputs[m.currentField].Blur()
					m.currentField++
					m.inputs[m.currentField].Focus()
					return m, m.inputs[m.currentField].Cursor.BlinkCmd()
				}
			}
		}
	}

	// Update the current input field
	var cmd tea.Cmd
	m.inputs[m.currentField], cmd = m.inputs[m.currentField].Update(msg)
	return m, cmd
}

// View renders the password change form interface
func (m ChangeModel) View() string {
	var content strings.Builder

	// Title
	title := changeTitleStyle.Render("🔄 Change Master Password")
	content.WriteString(title + "\n\n")

	// Form content
	formContent := ""

	// Current password field
	formContent += changeFieldLabelStyle.Render("Current Master Password")
	formContent += changeRequiredStyle.Render(" *") + "\n"
	formContent += "  " + m.inputs[CurrentPasswordField].View() + "\n\n"

	// New password field
	formContent += changeFieldLabelStyle.Render("New Master Password")
	formContent += changeRequiredStyle.Render(" *") + "\n"
	formContent += "  " + m.inputs[NewPasswordField].View() + "\n\n"

	// Confirm password field
	formContent += changeFieldLabelStyle.Render("Confirm New Password")
	formContent += changeRequiredStyle.Render(" *") + "\n"
	formContent += "  " + m.inputs[ConfirmPasswordField].View() + "\n\n"

	// Validation errors
	if errorMsg := m.getValidationError(); errorMsg != "" {
		formContent += changeErrorStyle.Render("❌ " + errorMsg) + "\n\n"
	}

	content.WriteString(changeContainerStyle.Render(formContent))

	// Help text
	help := changeHelpStyle.Render("Tab/↑↓: Navigate • Enter: Next/Submit • Esc: Cancel")
	content.WriteString(help)

	return content.String()
}

// validateForm checks if the form is valid for submission
func (m ChangeModel) validateForm() bool {
	currentPass := strings.TrimSpace(m.inputs[CurrentPasswordField].Value())
	newPass := strings.TrimSpace(m.inputs[NewPasswordField].Value())
	confirmPass := strings.TrimSpace(m.inputs[ConfirmPasswordField].Value())

	// Check required fields
	if currentPass == "" || newPass == "" || confirmPass == "" {
		return false
	}

	// Check new password length
	if len(newPass) < 8 {
		return false
	}

	// Check passwords match
	if newPass != confirmPass {
		return false
	}

	return true
}

// getValidationError returns the current validation error message
func (m ChangeModel) getValidationError() string {
	currentPass := strings.TrimSpace(m.inputs[CurrentPasswordField].Value())
	newPass := strings.TrimSpace(m.inputs[NewPasswordField].Value())
	confirmPass := strings.TrimSpace(m.inputs[ConfirmPasswordField].Value())

	if currentPass == "" && m.currentField > int(CurrentPasswordField) {
		return "Current password is required"
	}

	if newPass == "" && m.currentField > int(NewPasswordField) {
		return "New password is required"
	}

	if len(newPass) > 0 && len(newPass) < 8 {
		return "New password must be at least 8 characters long"
	}

	if confirmPass == "" && m.currentField > int(ConfirmPasswordField) {
		return "Password confirmation is required"
	}

	if newPass != "" && confirmPass != "" && newPass != confirmPass {
		return "New passwords do not match"
	}

	return ""
}

// GetFormData returns the form data
func (m ChangeModel) GetFormData() (string, string, string) {
	return m.currentPass, m.newPass, m.confirmPass
}

// IsSubmitted returns whether the form was successfully submitted
func (m ChangeModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns whether the form was cancelled
func (m ChangeModel) IsCancelled() bool {
	return m.cancelled
}