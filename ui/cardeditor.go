package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// CardEditorModel represents the card editor screen
type CardEditorModel struct {
	width  int
	height int
}

// NewCardEditorModel creates a new card editor model
func NewCardEditorModel() *CardEditorModel {
	return &CardEditorModel{}
}

// SetSize sets the terminal size
func (m *CardEditorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model
func (m *CardEditorModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *CardEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement card editor logic
	return m, nil
}

// View implements tea.Model
func (m *CardEditorModel) View() string {
	return "Card editor screen - coming soon!"
}
