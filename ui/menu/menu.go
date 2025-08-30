// Package menu provides the main menu interface for password management operations.
// It uses Bubble Tea for TUI functionality and allows users to navigate between
// different password management actions like viewing, adding, editing, and deleting passwords.
package menu

import (
	"fmt"
	"strings"

	"github.com/Fozzyack/password-manager/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a single menu option with its display text and action identifier
type MenuItem struct {
	Title       string // Display text for the menu item
	Description string // Brief description of what this option does
	Action      string // Unique identifier for the action
}

// MenuModel represents the state of the main menu interface
type MenuModel struct {
	choices      []MenuItem      // Available menu options
	cursor       int             // Currently selected menu item index
	selected     bool            // Whether an item has been selected
	selectedItem string          // The action identifier of the selected item
	options      *types.Options  // Shared application options
}

// Menu styling with Lipgloss
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)

	menuStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(60).
		Align(lipgloss.Left)

	selectedItemStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
		Padding(0, 1)

	helpTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	containerStyle = lipgloss.NewStyle().
		Padding(2, 4).
		Margin(1, 2).
		Align(lipgloss.Center)
)

// InitialMenuModel creates a new menu model with predefined password management options
func InitialMenuModel(options *types.Options) MenuModel {
	return MenuModel{
		choices: []MenuItem{
			{
				Title:       "📋 List All Passwords",
				Description: "View all stored password entries",
				Action:      "list",
			},
			{
				Title:       "➕ Add New Password",
				Description: "Create a new password entry",
				Action:      "add",
			},
			{
				Title:       "🔄 Change Master Password",
				Description: "Update your master password",
				Action:      "change_master",
			},
			{
				Title:       "📤 Export Passwords",
				Description: "Export passwords to file",
				Action:      "export",
			},
			{
				Title:       "🚪 Quit",
				Description: "Exit the password manager",
				Action:      "quit",
			},
		},
		cursor:       0,
		selected:     false,
		selectedItem: "",
		options:      options,
	}
}

// Init implements the tea.Model interface
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles user input and menu navigation
func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.options.Quit = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = true
			m.selectedItem = m.choices[m.cursor].Action
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the menu interface
func (m MenuModel) View() string {
	var content strings.Builder

	// Title
	title := titleStyle.Render("🔐 Password Manager - Main Menu")
	content.WriteString(title + "\n\n")

	// Menu items with consistent width to prevent shifting
	menuContent := ""
	itemWidth := 52 // Account for padding inside the border
	
	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
			titleLine := fmt.Sprintf("%s %s", cursor, choice.Title)
			menuContent += selectedItemStyle.Width(itemWidth).Render(titleLine)
			menuContent += "\n"
			descLine := fmt.Sprintf("   %s", choice.Description)
			menuContent += itemStyle.Width(itemWidth).Render(descLine)
		} else {
			titleLine := fmt.Sprintf("%s %s", cursor, choice.Title)
			menuContent += itemStyle.Width(itemWidth).Render(titleLine)
		}
		menuContent += "\n"
		if i < len(m.choices)-1 {
			menuContent += "\n" // spacing between items
		}
	}

	content.WriteString(menuStyle.Render(menuContent))
	content.WriteString("\n")

	// Help text
	help := helpTextStyle.Render("Use ↑↓ or j/k to navigate • Enter to select • Esc to quit")
	content.WriteString(help)

	// Wrap in container
	return containerStyle.Render(content.String())
}

// GetSelectedAction returns the action identifier of the selected menu item
func (m MenuModel) GetSelectedAction() string {
	return m.selectedItem
}

// IsSelected returns whether a menu item has been selected
func (m MenuModel) IsSelected() bool {
	return m.selected
}
