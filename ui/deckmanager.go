package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// DeckManagerModel represents the deck management screen
type DeckManagerModel struct {
	width  int
	height int
}

// NewDeckManagerModel creates a new deck manager model
func NewDeckManagerModel() *DeckManagerModel {
	return &DeckManagerModel{}
}

// SetSize sets the terminal size
func (m *DeckManagerModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model
func (m *DeckManagerModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *DeckManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement deck management logic
	return m, nil
}

// View implements tea.Model
func (m *DeckManagerModel) View() string {
	return "Deck manager screen - coming soon!"
}
