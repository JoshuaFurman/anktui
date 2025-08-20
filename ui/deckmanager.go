package ui

import (
	"anktui/models"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DeckManagerState represents the current state of the deck manager
type DeckManagerState int

const (
	DeckManagerMenu DeckManagerState = iota
	CreatingDeck
	EditingDeck
	DeletingDeck
	ManagingCards
)

// DeckManagerModel represents the deck management screen
type DeckManagerModel struct {
	decks        []*models.Deck
	selectedDeck int
	state        DeckManagerState
	editingDeck  *models.Deck
	isNewDeck    bool

	// Form fields
	nameInput        string
	descriptionInput string
	currentField     int

	// Confirmation
	confirmingDelete bool

	width  int
	height int
}

// NewDeckManagerModel creates a new deck manager model
func NewDeckManagerModel(decks []*models.Deck, editDeck *models.Deck) *DeckManagerModel {
	m := &DeckManagerModel{
		decks:        decks,
		selectedDeck: 0,
		state:        DeckManagerMenu,
	}

	// If a deck is passed for editing, go straight to edit mode
	if editDeck != nil {
		m.state = EditingDeck
		m.editingDeck = editDeck
		m.nameInput = editDeck.Name
		m.descriptionInput = editDeck.Description
		m.isNewDeck = false
	}

	return m
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case DeckManagerMenu:
			return m.updateMenu(msg)
		case CreatingDeck, EditingDeck:
			return m.updateForm(msg)
		case DeletingDeck:
			return m.updateDelete(msg)
		}
	case DecksLoadedMsg:
		// Handle successful deck operations
		m.decks = msg.Decks

		// If we were creating or editing a deck, return to menu
		if m.state == CreatingDeck || m.state == EditingDeck {
			m.state = DeckManagerMenu
			m.nameInput = ""
			m.descriptionInput = ""
			m.currentField = 0
			m.editingDeck = nil
			m.isNewDeck = false
		}

		// If we were deleting, also return to menu
		if m.state == DeletingDeck {
			m.state = DeckManagerMenu
			m.confirmingDelete = false
			// Adjust selected deck if it was deleted
			if m.selectedDeck >= len(m.decks) && len(m.decks) > 0 {
				m.selectedDeck = len(m.decks) - 1
			} else if len(m.decks) == 0 {
				m.selectedDeck = 0
			}
		}
	}

	return m, nil
}

// updateMenu handles the main deck manager menu
func (m *DeckManagerModel) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedDeck > 0 {
			m.selectedDeck--
		}
	case "down", "j":
		if len(m.decks) > 0 && m.selectedDeck < len(m.decks)-1 {
			m.selectedDeck++
		}
	case "n":
		// Create new deck
		m.state = CreatingDeck
		m.editingDeck = &models.Deck{}
		m.nameInput = ""
		m.descriptionInput = ""
		m.currentField = 0
		m.isNewDeck = true
	case "e", "enter":
		if len(m.decks) > 0 {
			// Edit selected deck
			selectedDeck := m.decks[m.selectedDeck]
			m.state = EditingDeck
			m.editingDeck = selectedDeck
			m.nameInput = selectedDeck.Name
			m.descriptionInput = selectedDeck.Description
			m.currentField = 0
			m.isNewDeck = false
		}
	case "d":
		if len(m.decks) > 0 {
			// Delete selected deck
			m.state = DeletingDeck
			m.confirmingDelete = false
		}
	case "c":
		if len(m.decks) > 0 {
			// Manage cards in selected deck
			selectedDeck := m.decks[m.selectedDeck]
			return m, func() tea.Msg {
				return NavigateMsg{
					Screen: CardEditorScreen,
					Data:   selectedDeck,
				}
			}
		}
	case "esc":
		return m, func() tea.Msg {
			return NavigateMsg{Screen: MenuScreen}
		}
	}

	return m, nil
}

// updateForm handles deck creation/editing form
func (m *DeckManagerModel) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		if m.currentField < 1 {
			m.currentField++
		}
	case "shift+tab", "up":
		if m.currentField > 0 {
			m.currentField--
		}
	case "enter":
		// Save deck
		if m.nameInput == "" {
			return m, nil // Don't save without name
		}

		m.editingDeck.Name = m.nameInput
		m.editingDeck.Description = m.descriptionInput

		if m.isNewDeck {
			// Create new deck
			return m, func() tea.Msg {
				return CreateDeckMsg{Deck: m.editingDeck}
			}
		} else {
			// Update existing deck
			return m, func() tea.Msg {
				return UpdateDeckMsg{Deck: m.editingDeck}
			}
		}
	case "esc":
		// Cancel editing
		m.state = DeckManagerMenu
	default:
		// Handle text input
		if m.currentField == 0 {
			// Name field
			if msg.String() == "backspace" {
				if len(m.nameInput) > 0 {
					m.nameInput = m.nameInput[:len(m.nameInput)-1]
				}
			} else if len(msg.Runes) > 0 {
				m.nameInput += string(msg.Runes)
			}
		} else {
			// Description field
			if msg.String() == "backspace" {
				if len(m.descriptionInput) > 0 {
					m.descriptionInput = m.descriptionInput[:len(m.descriptionInput)-1]
				}
			} else if len(msg.Runes) > 0 {
				m.descriptionInput += string(msg.Runes)
			}
		}
	}

	return m, nil
}

// updateDelete handles deck deletion confirmation
func (m *DeckManagerModel) updateDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm deletion
		selectedDeck := m.decks[m.selectedDeck]
		return m, func() tea.Msg {
			return DeleteDeckMsg{Deck: selectedDeck}
		}
	case "n", "N", "esc":
		// Cancel deletion
		m.state = DeckManagerMenu
	}

	return m, nil
}

// View implements tea.Model
func (m *DeckManagerModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	switch m.state {
	case DeckManagerMenu:
		return m.viewMenu()
	case CreatingDeck:
		return m.viewForm("Create New Deck")
	case EditingDeck:
		return m.viewForm("Edit Deck")
	case DeletingDeck:
		return m.viewDelete()
	default:
		return "Unknown state"
	}
}

// viewMenu renders the main deck manager menu
func (m *DeckManagerModel) viewMenu() string {
	// Title
	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render("Deck Manager")

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

			stats := fmt.Sprintf("Cards: %d (New: %d, Review: %d)", total, new, review)

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

			if i == m.selectedDeck {
				itemStyle = itemStyle.
					BorderForeground(primaryColor)
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
		helpText = "↑/↓: navigate • Enter/e: edit • c: manage cards • d: delete • n: new deck • Esc: back"
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

// viewForm renders the deck creation/editing form
func (m *DeckManagerModel) viewForm(titleText string) string {
	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render(titleText)

	// Name field
	nameFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		PaddingLeft(2).
		PaddingRight(2).
		Width(50)

	if m.currentField == 0 {
		nameFieldStyle = nameFieldStyle.BorderForeground(primaryColor)
	}

	nameLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		PaddingBottom(1).
		Render("Deck Name:")

	nameValue := m.nameInput
	if m.currentField == 0 {
		nameValue += "█"
	}

	nameField := lipgloss.JoinVertical(
		lipgloss.Left,
		nameLabel,
		nameFieldStyle.Render(nameValue),
	)

	// Description field
	descFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		PaddingLeft(2).
		PaddingRight(2).
		Width(50).
		Height(3)

	if m.currentField == 1 {
		descFieldStyle = descFieldStyle.BorderForeground(primaryColor)
	}

	descLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		PaddingBottom(1).
		PaddingTop(2).
		Render("Description (optional):")

	descValue := m.wrapText(m.descriptionInput, 46)
	if m.currentField == 1 {
		descValue += "█"
	}

	descField := lipgloss.JoinVertical(
		lipgloss.Left,
		descLabel,
		descFieldStyle.Render(descValue),
	)

	// Help text
	help := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(3).
		Render("Tab/↑↓: switch fields • Enter: save • Esc: cancel")

	// Combine all elements
	form := lipgloss.JoinVertical(
		lipgloss.Left,
		nameField,
		descField,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		form,
		help,
	)

	// Center everything in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// viewDelete renders the deletion confirmation
func (m *DeckManagerModel) viewDelete() string {
	selectedDeck := m.decks[m.selectedDeck]

	title := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render("⚠️  Delete Deck")

	warning := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render(fmt.Sprintf("Are you sure you want to delete '%s'?", selectedDeck.Name))

	cardCount := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingBottom(3).
		Render(fmt.Sprintf("This will permanently delete %d cards.", len(selectedDeck.Cards)))

	help := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		Render("Y: confirm deletion • N/Esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		warning,
		cardCount,
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// wrapText wraps text to the specified width
func (m *DeckManagerModel) wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// Message types for deck operations
type CreateDeckMsg struct {
	Deck *models.Deck
}

type UpdateDeckMsg struct {
	Deck *models.Deck
}

type DeleteDeckMsg struct {
	Deck *models.Deck
}
