package ui

import (
	"anktui/models"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DeckListState represents the current state of the deck list screen
type DeckListState int

const (
	SelectingDeck DeckListState = iota
	SelectingMode
)

// DeckListModel represents the deck selection screen
type DeckListModel struct {
	decks        []*models.Deck
	selected     int
	selectedMode int
	state        DeckListState
	width        int
	height       int
}

// NewDeckListModel creates a new deck list model
func NewDeckListModel(decks []*models.Deck) *DeckListModel {
	return &DeckListModel{
		decks:        decks,
		selected:     0,
		selectedMode: 0, // Default to Review Mode
		state:        SelectingDeck,
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
		switch m.state {
		case SelectingDeck:
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
					// Move to mode selection
					m.state = SelectingMode
					m.selectedMode = 0 // Default to Review Mode
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

		case SelectingMode:
			switch msg.String() {
			case "up", "k":
				if m.selectedMode > 0 {
					m.selectedMode--
				}
			case "down", "j":
				if m.selectedMode < 1 {
					m.selectedMode++
				}
			case "enter", " ":
				// Start studying with selected mode
				selectedDeck := m.decks[m.selected]
				mode := models.ReviewMode
				if m.selectedMode == 1 {
					mode = models.PracticeMode
				}
				return m, func() tea.Msg {
					return NavigateMsg{
						Screen: StudyScreen,
						Data: &StudyRequest{
							Deck: selectedDeck,
							Mode: mode,
						},
					}
				}
			case "esc":
				// Go back to deck selection
				m.state = SelectingDeck
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

	switch m.state {
	case SelectingDeck:
		return m.viewDeckSelection()
	case SelectingMode:
		return m.viewModeSelection()
	default:
		return "Unknown state"
	}
}

// viewDeckSelection renders the deck selection screen
func (m *DeckListModel) viewDeckSelection() string {
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
		helpText = "↑/↓ or j/k: navigate • Enter: select deck • e: edit • d: delete • n: new deck • Esc: back"
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

// viewModeSelection renders the study mode selection screen
func (m *DeckListModel) viewModeSelection() string {
	selectedDeck := m.decks[m.selected]

	// Title
	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(1).
		Render(fmt.Sprintf("Study Mode for: %s", selectedDeck.Name))

	// Mode options
	modes := []struct {
		name        string
		description string
	}{
		{"📚 Review Mode", "Only cards due for review + new cards"},
		{"🔄 Practice Mode", "All cards for practice (ignores schedule)"},
	}

	var modeItems []string
	for i, mode := range modes {
		// Style the item
		itemStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			PaddingLeft(3).
			PaddingRight(3).
			PaddingTop(1).
			PaddingBottom(1).
			Margin(0, 2, 1, 2).
			Width(50)

		if i == m.selectedMode {
			itemStyle = itemStyle.
				BorderForeground(secondaryColor).
				Background(lipgloss.Color("#374151"))
		}

		modeContent := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Foreground(textColor).Render(mode.name),
			lipgloss.NewStyle().Foreground(mutedColor).Render(mode.description),
		)

		modeItems = append(modeItems, itemStyle.Render(modeContent))
	}

	// Join mode items
	modeList := lipgloss.JoinVertical(lipgloss.Center, modeItems...)

	// Help text
	help := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(2).
		Render("↑/↓ or j/k: navigate • Enter: start study • Esc: back to deck list")

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		modeList,
		help,
	)

	// Center everything in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
