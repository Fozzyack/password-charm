// Package list provides a scrollable list view for displaying password entries.
// It uses Bubble Tea for TUI functionality and maintains consistent styling with the rest of the application.
package list

import (
	"fmt"
	"strings"
	"time"

	"github.com/Fozzyack/password-manager/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PasswordEntry represents a password entry in the list view
type PasswordEntry struct {
	Filename  string    // The actual filename (without .gpg)
	SiteName  string    // Display name for the site
	Username  string    // Username for the entry
	Email     string    // Email for the entry
	CreatedAt time.Time // When the entry was created
}

// ListModel represents the state of the password list
type ListModel struct {
	entries       []PasswordEntry
	cursor        int
	selected      bool
	selectedEntry PasswordEntry
	options       *types.Options
}

// List styling
var (
	listTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center)

	listContainerStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Margin(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Width(80).
		Align(lipgloss.Left)

	listItemStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 0, 1, 0)

	selectedItemStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Margin(0, 0, 1, 0).
		Background(lipgloss.Color("#7D56F4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	listHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		PaddingLeft(4).
		Italic(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	emptyListStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Italic(true).
		Align(lipgloss.Center).
		Padding(4, 2)
)

// NewPasswordList creates a new password list with the given entries
func NewPasswordList(entries []PasswordEntry, options *types.Options) ListModel {
	// Clear screen for clean list display
	fmt.Print("\033[2J\033[H")

	return ListModel{
		entries: entries,
		cursor:  0,
		options: options,
	}
}

// Init implements the tea.Model interface
func (m ListModel) Init() tea.Cmd {
	return nil
}

// Update handles user input and list navigation
func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			// Return to main menu
			return m, tea.Quit

		case "enter", " ":
			// Select the current entry
			if len(m.entries) > 0 && m.cursor < len(m.entries) {
				m.selected = true
				m.selectedEntry = m.entries[m.cursor]
				return m, tea.Quit
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case "home":
			m.cursor = 0

		case "end":
			if len(m.entries) > 0 {
				m.cursor = len(m.entries) - 1
			}
		}
	}

	return m, nil
}

// View renders the password list interface
func (m ListModel) View() string {
	var content strings.Builder

	// Title
	title := listTitleStyle.Render("🔐 Password List")
	content.WriteString(title + "\n\n")

	// Check if list is empty
	if len(m.entries) == 0 {
		emptyMsg := emptyListStyle.Render("No passwords found.\nUse the 'Add New Password' option to create your first entry.")
		content.WriteString(listContainerStyle.Render(emptyMsg))
		content.WriteString(listHelpStyle.Render("Press Esc to return to main menu"))
		return content.String()
	}

	// List content
	listContent := ""
	
	// Add header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 2).
		Margin(0, 0, 1, 0)
	
	listContent += headerStyle.Render(fmt.Sprintf("%-25s %-20s %-15s %s", 
		"Site/Service", "Username", "Email", "Created"))
	listContent += "\n"
	
	// Add separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262"))
	listContent += separatorStyle.Render(strings.Repeat("─", 70)) + "\n\n"

	// List entries
	for i, entry := range m.entries {
		// Format the entry data
		siteName := entry.SiteName
		if len(siteName) > 24 {
			siteName = siteName[:21] + "..."
		}
		
		username := entry.Username
		if len(username) > 19 {
			username = username[:16] + "..."
		}
		
		email := entry.Email
		if len(email) > 14 {
			email = email[:11] + "..."
		}
		
		createdAt := entry.CreatedAt.Format("Jan 02, 2006")
		
		entryText := fmt.Sprintf("%-25s %-20s %-15s %s", 
			siteName, username, email, createdAt)

		// Apply styling based on cursor position
		if i == m.cursor {
			listContent += selectedItemStyle.Render("► " + entryText) + "\n"
		} else {
			listContent += listItemStyle.Render("  " + entryText) + "\n"
		}
	}

	content.WriteString(listContainerStyle.Render(listContent))

	// Help text
	help := listHelpStyle.Render("↑↓/j/k: Navigate • Enter/Space: View Details • Esc/q: Back to Menu")
	content.WriteString(help)

	return content.String()
}

// IsSelected returns whether an entry was selected
func (m ListModel) IsSelected() bool {
	return m.selected
}

// GetSelectedEntry returns the selected password entry
func (m ListModel) GetSelectedEntry() PasswordEntry {
	return m.selectedEntry
}

// GetCursor returns the current cursor position
func (m ListModel) GetCursor() int {
	return m.cursor
}
