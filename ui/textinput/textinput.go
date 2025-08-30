// Package textinput provides a styled terminal input component for the password manager.
// It uses Bubble Tea for TUI functionality and Lipgloss for styling, with support for
// password masking, error display, and clean visual design.
package textinput

import (
	"fmt"

	"github.com/Fozzyack/password-manager/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type (
	errMsg error
)

var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)

	containerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Margin(1, 2).
		Align(lipgloss.Center)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87")).
		Bold(true).
		Align(lipgloss.Center)
)

type model struct {
	textInput textinput.Model
	err       error
	header    string
	output    *string
	options *types.Options
}


func InitialModel(header string, placeholder string, output *string, options *types.Options) model {
	return InitialModelWithMasking(header, placeholder, output, options, true)
}

// InitialModelWithMasking creates a textinput model with optional password masking control.
// When maskPassword is false, the input will be visible even for password fields.
func InitialModelWithMasking(header string, placeholder string, output *string, options *types.Options, maskPassword bool) model {
	// Clear the screen when starting a new input session
	fmt.Print("\033[2J\033[H")
	
	ti := textinput.New()
	ti.Placeholder = placeholder 
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	
	// Style the textinput
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Italic(true)
	
	// Set password mode if this is a password field and masking is enabled
	if placeholder == "Password" && maskPassword {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	return model{
		textInput: ti,
		err:       nil,
		output:    output,
		header:    header,
		options:      options,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			*m.output = m.textInput.Value()
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.options.Quit = true
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var content string
	
	// Create the header
	header := headerStyle.Render(m.header)
	
	// Create the input field with some spacing
	input := fmt.Sprintf("\n%s\n", m.textInput.View())
	
	// Create the help text
	help := helpStyle.Render("Press Enter to continue • Esc to quit")
	
	// Handle error display
	errorMsg := ""
	if m.err != nil {
		errorMsg = errorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\n"
	} else if m.options.ErrorMessage != "" {
		errorMsg = errorStyle.Render(m.options.ErrorMessage) + "\n\n"
	}
	
	// Combine all elements
	content = fmt.Sprintf("%s%s\n\n%s%s\n\n%s", 
		errorMsg,
		header, 
		input, 
		help,
		"\n",
	)
	
	// Wrap in container for final styling
	return containerStyle.Render(content)
}
