// Package form provides a multi-field input form for adding new password entries.
// It uses Bubble Tea for TUI functionality and maintains consistent styling with the rest of the application.
package form

import (
	"fmt"
	"strings"

	"github.com/Fozzyack/password-manager/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents a single input field in the form
type FormField struct {
	Label       string
	Placeholder string
	Required    bool
	Masked      bool
	Value       string
}

// FormModel represents the state of the multi-field form
type FormModel struct {
	fields       []FormField
	inputs       []textinput.Model
	currentField int
	submitted    bool
	cancelled    bool
	options      *types.Options
}

// Form styling
var (
	formTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)

	formContainerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Margin(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(70).
		Align(lipgloss.Left)

	fieldLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Margin(0, 0, 0, 1)

	requiredStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true).
		Align(lipgloss.Left).
		Margin(0, 0, 1, 1)
)

// NewPasswordForm creates a new password entry form with predefined fields
func NewPasswordForm(options *types.Options) FormModel {
	// Clear screen for clean form display
	fmt.Print("\033[2J\033[H")

	fields := []FormField{
		{
			Label:       "Site/Service Name",
			Placeholder: "e.g., Gmail, GitHub, Banking",
			Required:    true,
			Masked:      false,
		},
		{
			Label:       "Username",
			Placeholder: "your_username",
			Required:    false,
			Masked:      false,
		},
		{
			Label:       "Email",
			Placeholder: "user@example.com",
			Required:    false,
			Masked:      false,
		},
		{
			Label:       "URL",
			Placeholder: "https://example.com",
			Required:    false,
			Masked:      false,
		},
		{
			Label:       "Password",
			Placeholder: "Enter password or generate one",
			Required:    true,
			Masked:      true,
		},
	}

	inputs := make([]textinput.Model, len(fields))
	for i := range inputs {
		ti := textinput.New()
		ti.Placeholder = fields[i].Placeholder
		ti.CharLimit = 200
		ti.Width = 50

		// Style the textinput
		ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
		ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true)

		// Set password masking for password field
		if fields[i].Masked {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}

		// Focus on the first field
		if i == 0 {
			ti.Focus()
		}

		inputs[i] = ti
	}

	return FormModel{
		fields:       fields,
		inputs:       inputs,
		currentField: 0,
		submitted:    false,
		cancelled:    false,
		options:      options,
	}
}

// Init implements the tea.Model interface
func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles user input and form navigation
func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				// Validate required fields before submitting
				if m.validateForm() {
					// Update field values
					for i := range m.fields {
						m.fields[i].Value = m.inputs[i].Value()
					}
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

// View renders the form interface
func (m FormModel) View() string {
	var content strings.Builder

	// Title
	title := formTitleStyle.Render("➕ Add New Password Entry")
	content.WriteString(title + "\n\n")

	// Form fields
	formContent := ""
	for i, field := range m.fields {
		// Field label
		label := field.Label
		if field.Required {
			label += requiredStyle.Render(" *")
		}
		formContent += fieldLabelStyle.Render(label) + "\n"

		// Input field with focus styling
		if i == m.currentField {
			formContent += "  " + m.inputs[i].View() + "\n"
		} else {
			formContent += "  " + m.inputs[i].View() + "\n"
		}

		formContent += "\n"
	}

	// Validation errors
	errorMsg := ""
	if !m.validateForm() && m.currentField == len(m.inputs)-1 {
		errorMsg = m.getValidationError()
		if errorMsg != "" {
			formContent += errorStyle.Render("❌ " + errorMsg) + "\n\n"
		}
	}

	content.WriteString(formContainerStyle.Render(formContent))

	// Help text
	help := helpStyle.Render("Tab/Enter: Next field • ↑↓: Navigate • Enter on last field: Save • Esc: Cancel")
	content.WriteString(help)

	return content.String()
}

// validateForm checks if all required fields are filled
func (m FormModel) validateForm() bool {
	for i, field := range m.fields {
		if field.Required && strings.TrimSpace(m.inputs[i].Value()) == "" {
			return false
		}
	}
	return true
}

// getValidationError returns a validation error message
func (m FormModel) getValidationError() string {
	for i, field := range m.fields {
		if field.Required && strings.TrimSpace(m.inputs[i].Value()) == "" {
			return fmt.Sprintf("'%s' is required", field.Label)
		}
	}
	return ""
}

// GetFormData returns the form data as a map
func (m FormModel) GetFormData() map[string]string {
	data := make(map[string]string)
	for _, field := range m.fields {
		// Convert field label to lowercase and replace special characters/spaces with underscores
		key := strings.ToLower(field.Label)
		key = strings.ReplaceAll(key, "/", "_")
		key = strings.ReplaceAll(key, " ", "_")
		data[key] = field.Value
	}
	return data
}

// IsSubmitted returns whether the form was successfully submitted
func (m FormModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns whether the form was cancelled
func (m FormModel) IsCancelled() bool {
	return m.cancelled
}
