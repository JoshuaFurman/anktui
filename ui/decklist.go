package ui

import (
	"anktui/models"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DeckListModel represents the deck selection screen
type DeckListModel struct {
	decks    []*models.Deck
	selected int
	width    int
	height   int
}

// NewDeckListModel creates a new deck list model
func NewDeckListModel(decks []*models.Deck) *DeckListModel {
	return &DeckListModel{
		decks:    decks,
		selected: 0,
	}
}

// SetSize sets the terminal size
func (m *DeckListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model
func (m *DeckListModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *DeckListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if len(m.decks) > 0 && m.selected < len(m.decks)-1 {
				m.selected++
			}
		case "enter", " ":
			if len(m.decks) > 0 {
				// Start studying the selected deck
				selectedDeck := m.decks[m.selected]
				return m, func() tea.Msg {
					return NavigateMsg{
						Screen: StudyScreen,
						Data:   selectedDeck,
					}
				}
			}
		case "n":
			// Create new deck
			return m, func() tea.Msg {
				return NavigateMsg{Screen: DeckManagerScreen}
			}
		case "e":
			if len(m.decks) > 0 {
				// Edit selected deck
				selectedDeck := m.decks[m.selected]
				return m, func() tea.Msg {
					return NavigateMsg{
						Screen: DeckManagerScreen,
						Data:   selectedDeck,
					}
				}
			}
		case "d":
			if len(m.decks) > 0 {
				// Delete selected deck (TODO: implement confirmation)
				// For now, just return to menu
				return m, func() tea.Msg {
					return NavigateMsg{Screen: MenuScreen}
				}
			}
		case "esc":
			return m, func() tea.Msg {
				return NavigateMsg{Screen: MenuScreen}
			}
		}
	}

	return m, nil
}

// View implements tea.Model
func (m *DeckListModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Title
	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render("Select a Deck to Study")

	// Create deck list
	var deckItems []string

	if len(m.decks) == 0 {
		noDeckMsg := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Align(lipgloss.Center).
			Render("No decks found. Press 'n' to create a new deck.")
		deckItems = []string{noDeckMsg}
	} else {
		for i, deck := range m.decks {
			total, new, review := deck.GetCardStats()

			// Create deck info
			deckName := deck.Name
			if deck.Description != "" {
				deckName = fmt.Sprintf("%s - %s", deck.Name, deck.Description)
			}

			stats := fmt.Sprintf("Total: %d • New: %d • Review: %d", total, new, review)

			// Style the item
			itemStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				PaddingBottom(1).
				Margin(0, 2, 1, 2).
				Width(60)

			if i == m.selected {
				itemStyle = itemStyle.
					BorderForeground(primaryColor).
					Background(lipgloss.Color("#374151"))
			}

			deckContent := lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(textColor).Render(deckName),
				lipgloss.NewStyle().Foreground(mutedColor).Render(stats),
			)

			deckItems = append(deckItems, itemStyle.Render(deckContent))
		}
	}

	// Join deck items
	deckList := lipgloss.JoinVertical(lipgloss.Center, deckItems...)

	// Help text
	var helpText string
	if len(m.decks) > 0 {
		helpText = "↑/↓ or j/k: navigate • Enter: study • e: edit • d: delete • n: new deck • Esc: back"
	} else {
		helpText = "n: create new deck • Esc: back to menu"
	}

	help := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(2).
		Render(helpText)

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		deckList,
		help,
	)

	// Center everything in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
