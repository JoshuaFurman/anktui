package ui

import (
	"anktui/models"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CardEditorState represents the current state of the card editor
type CardEditorState int

const (
	CardListView CardEditorState = iota
	CardForm
	CardDeleteConfirm
)

// CardEditorModel represents the card editor screen
type CardEditorModel struct {
	deck         *models.Deck
	state        CardEditorState
	selectedCard int
	editingCard  *models.Card
	isNewCard    bool

	// Form fields
	frontInput   string
	backInput    string
	currentField int

	// Confirmation
	confirmingDelete bool

	width  int
	height int
}

// NewCardEditorModel creates a new card editor model
func NewCardEditorModel(deck *models.Deck) *CardEditorModel {
	return &CardEditorModel{
		deck:         deck,
		state:        CardListView,
		selectedCard: 0,
	}
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case CardListView:
			return m.updateList(msg)
		case CardForm:
			return m.updateForm(msg)
		case CardDeleteConfirm:
			return m.updateDelete(msg)
		}
	case DecksLoadedMsg:
		// Update deck data with fresh information
		for _, deck := range msg.Decks {
			if deck.ID == m.deck.ID {
				m.deck = deck
				break
			}
		}

		// If we were creating or editing a card, return to list view
		if m.state == CardForm {
			m.state = CardListView
			m.frontInput = ""
			m.backInput = ""
			m.currentField = 0
			m.editingCard = nil
			m.isNewCard = false
		}

		// If we were deleting, also return to list view
		if m.state == CardDeleteConfirm {
			m.state = CardListView
			m.confirmingDelete = false
			// Adjust selected card if it was deleted
			if m.selectedCard >= len(m.deck.Cards) && len(m.deck.Cards) > 0 {
				m.selectedCard = len(m.deck.Cards) - 1
			} else if len(m.deck.Cards) == 0 {
				m.selectedCard = 0
			}
		}
	}

	return m, nil
}

// updateList handles the card list view
func (m *CardEditorModel) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedCard > 0 {
			m.selectedCard--
		}
	case "down", "j":
		if len(m.deck.Cards) > 0 && m.selectedCard < len(m.deck.Cards)-1 {
			m.selectedCard++
		}
	case "n":
		// Create new card
		m.state = CardForm
		m.editingCard = &models.Card{}
		m.frontInput = ""
		m.backInput = ""
		m.currentField = 0
		m.isNewCard = true
	case "e", "enter":
		if len(m.deck.Cards) > 0 {
			// Edit selected card
			selectedCard := &m.deck.Cards[m.selectedCard]
			m.state = CardForm
			m.editingCard = selectedCard
			m.frontInput = selectedCard.Front
			m.backInput = selectedCard.Back
			m.currentField = 0
			m.isNewCard = false
		}
	case "d":
		if len(m.deck.Cards) > 0 {
			// Delete selected card
			m.state = CardDeleteConfirm
		}
	case "esc":
		return m, func() tea.Msg {
			return NavigateMsg{Screen: DeckManagerScreen}
		}
	}

	return m, nil
}

// updateForm handles card creation/editing form
func (m *CardEditorModel) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		// Save card
		if m.frontInput == "" || m.backInput == "" {
			return m, nil // Don't save without both sides
		}

		m.editingCard.Front = m.frontInput
		m.editingCard.Back = m.backInput

		if m.isNewCard {
			// Create new card
			return m, func() tea.Msg {
				return CreateCardMsg{Deck: m.deck, Card: m.editingCard}
			}
		} else {
			// Update existing card
			return m, func() tea.Msg {
				return UpdateCardMsg{Deck: m.deck, Card: m.editingCard}
			}
		}
	case "esc":
		// Cancel editing
		m.state = CardListView
	default:
		// Handle text input
		if m.currentField == 0 {
			// Front field
			if msg.String() == "backspace" {
				if len(m.frontInput) > 0 {
					m.frontInput = m.frontInput[:len(m.frontInput)-1]
				}
			} else if len(msg.Runes) > 0 {
				m.frontInput += string(msg.Runes)
			}
		} else {
			// Back field
			if msg.String() == "backspace" {
				if len(m.backInput) > 0 {
					m.backInput = m.backInput[:len(m.backInput)-1]
				}
			} else if len(msg.Runes) > 0 {
				m.backInput += string(msg.Runes)
			}
		}
	}

	return m, nil
}

// updateDelete handles card deletion confirmation
func (m *CardEditorModel) updateDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm deletion
		selectedCard := &m.deck.Cards[m.selectedCard]
		return m, func() tea.Msg {
			return DeleteCardMsg{Deck: m.deck, Card: selectedCard}
		}
	case "n", "N", "esc":
		// Cancel deletion
		m.state = CardListView
	}

	return m, nil
}

// View implements tea.Model
func (m *CardEditorModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	switch m.state {
	case CardListView:
		return m.viewList()
	case CardForm:
		return m.viewForm()
	case CardDeleteConfirm:
		return m.viewDelete()
	default:
		return "Unknown state"
	}
}

// viewList renders the card list view
func (m *CardEditorModel) viewList() string {
	// Title
	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(1).
		Render(fmt.Sprintf("Cards in: %s", m.deck.Name))

	// Stats
	stats := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render(fmt.Sprintf("Total: %d cards", len(m.deck.Cards)))

	// Create card list
	var cardItems []string

	if len(m.deck.Cards) == 0 {
		noCardMsg := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Align(lipgloss.Center).
			Render("No cards found. Press 'n' to create a new card.")
		cardItems = []string{noCardMsg}
	} else {
		for i, card := range m.deck.Cards {
			// Truncate long text
			front := card.Front
			if len(front) > 40 {
				front = front[:37] + "..."
			}
			back := card.Back
			if len(back) > 40 {
				back = back[:37] + "..."
			}

			// Style the item
			itemStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor).
				PaddingLeft(2).
				PaddingRight(2).
				PaddingTop(1).
				PaddingBottom(1).
				Margin(0, 2, 1, 2).
				Width(70)

			if i == m.selectedCard {
				itemStyle = itemStyle.
					BorderForeground(secondaryColor).
					Background(lipgloss.Color("#374151"))
			}

			cardContent := lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(textColor).Render("Q: "+front),
				lipgloss.NewStyle().Foreground(mutedColor).Render("A: "+back),
			)

			cardItems = append(cardItems, itemStyle.Render(cardContent))
		}
	}

	// Join card items (limit height for scrolling)
	maxItems := 8 // Show max 8 cards at once
	startIdx := 0
	endIdx := len(cardItems)

	if len(cardItems) > maxItems {
		// Calculate scroll position
		if m.selectedCard >= maxItems/2 {
			startIdx = m.selectedCard - maxItems/2
			endIdx = startIdx + maxItems
			if endIdx > len(cardItems) {
				endIdx = len(cardItems)
				startIdx = endIdx - maxItems
			}
		} else {
			endIdx = maxItems
		}
	}

	visibleItems := cardItems[startIdx:endIdx]
	cardList := lipgloss.JoinVertical(lipgloss.Center, visibleItems...)

	// Help text
	var helpText string
	if len(m.deck.Cards) > 0 {
		helpText = "↑/↓: navigate • Enter/e: edit • d: delete • n: new card • Esc: back"
	} else {
		helpText = "n: create new card • Esc: back to deck manager"
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
		stats,
		cardList,
		help,
	)

	// Center everything in the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// viewForm renders the card creation/editing form
func (m *CardEditorModel) viewForm() string {
	var titleText string
	if m.isNewCard {
		titleText = "Create New Card"
	} else {
		titleText = "Edit Card"
	}

	title := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render(titleText)

	// Front field
	frontFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		PaddingLeft(2).
		PaddingRight(2).
		Width(60).
		Height(3)

	if m.currentField == 0 {
		frontFieldStyle = frontFieldStyle.BorderForeground(primaryColor)
	}

	frontLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		PaddingBottom(1).
		Render("Front (Question):")

	frontValue := m.wrapText(m.frontInput, 56)
	if m.currentField == 0 {
		frontValue += "█"
	}

	frontField := lipgloss.JoinVertical(
		lipgloss.Left,
		frontLabel,
		frontFieldStyle.Render(frontValue),
	)

	// Back field
	backFieldStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		PaddingLeft(2).
		PaddingRight(2).
		Width(60).
		Height(3)

	if m.currentField == 1 {
		backFieldStyle = backFieldStyle.BorderForeground(primaryColor)
	}

	backLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		PaddingBottom(1).
		PaddingTop(2).
		Render("Back (Answer):")

	backValue := m.wrapText(m.backInput, 56)
	if m.currentField == 1 {
		backValue += "█"
	}

	backField := lipgloss.JoinVertical(
		lipgloss.Left,
		backLabel,
		backFieldStyle.Render(backValue),
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
		frontField,
		backField,
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
func (m *CardEditorModel) viewDelete() string {
	selectedCard := &m.deck.Cards[m.selectedCard]

	title := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render("⚠️  Delete Card")

	// Show card preview
	front := selectedCard.Front
	if len(front) > 50 {
		front = front[:47] + "..."
	}

	cardPreview := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		PaddingLeft(2).
		PaddingRight(2).
		PaddingTop(1).
		PaddingBottom(1).
		Width(60).
		Align(lipgloss.Center).
		Render(front)

	warning := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center).
		PaddingTop(2).
		PaddingBottom(3).
		Render("Are you sure you want to delete this card?")

	help := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		Render("Y: confirm deletion • N/Esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		cardPreview,
		warning,
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// wrapText wraps text to the specified width
func (m *CardEditorModel) wrapText(text string, width int) string {
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

// Message types for card operations
type CreateCardMsg struct {
	Deck *models.Deck
	Card *models.Card
}

type UpdateCardMsg struct {
	Deck *models.Deck
	Card *models.Card
}

type DeleteCardMsg struct {
	Deck *models.Deck
	Card *models.Card
}
