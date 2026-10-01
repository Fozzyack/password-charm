// Package ui provides shared full-screen rendering for the terminal interface.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Run displays a model centered in the terminal and returns its final state.
func Run(model tea.Model) (tea.Model, error) {
	finalModel, err := tea.NewProgram(centeredModel{model: model}, tea.WithAltScreen()).Run()
	if finalModel == nil {
		return nil, err
	}
	return finalModel.(centeredModel).model, err
}

type centeredModel struct {
	model  tea.Model
	width  int
	height int
}

func (m centeredModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m centeredModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m centeredModel) View() string {
	// Wait for the initial terminal size to avoid briefly drawing at the top left.
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.model.View())
}
