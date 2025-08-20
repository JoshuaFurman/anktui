package ui

import (
	"anktui/algorithms"
	"anktui/models"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StudyState represents the current state of the study session
type StudyState int

const (
	ShowingQuestion StudyState = iota
	ShowingAnswer
	SessionComplete
)

// StudyModel represents the study session screen
type StudyModel struct {
	session        *models.StudySession
	deck           *models.Deck
	state          StudyState
	selectedRating int
	width          int
	height         int
}

// NewStudyModel creates a new study model
func NewStudyModel(deck *models.Deck, maxCards int) *StudyModel {
	session := models.NewStudySession(deck, maxCards)
	return &StudyModel{
		session:        session,
		deck:           deck,
		state:          ShowingQuestion,
		selectedRating: 2, // Default to "Good"
	}
}

// SetSize sets the terminal size
func (m *StudyModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model
func (m *StudyModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *StudyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case ShowingQuestion:
			switch msg.String() {
			case "space", "enter", "f":
				// Flip card to show answer
				m.session.ShowAnswer()
				m.state = ShowingAnswer
				return m, nil
			case "esc":
				// Return to deck list
				return m, func() tea.Msg {
					return NavigateMsg{Screen: DeckListScreen}
				}
			}

		case ShowingAnswer:
			switch msg.String() {
			case "1":
				return m.rateCardAndContinue(models.Again)
			case "2":
				return m.rateCardAndContinue(models.Hard)
			case "3":
				return m.rateCardAndContinue(models.Good)
			case "4":
				return m.rateCardAndContinue(models.Easy)
			case "left", "h":
				if m.selectedRating > 0 {
					m.selectedRating--
				}
			case "right", "l":
				if m.selectedRating < 3 {
					m.selectedRating++
				}
			case "enter", "space":
				return m.rateCardAndContinue(models.Rating(m.selectedRating))
			case "esc":
				// Return to deck list
				return m, func() tea.Msg {
					return NavigateMsg{Screen: DeckListScreen}
				}
			}

		case SessionComplete:
			switch msg.String() {
			case "enter", "space", "esc":
				// Return to deck list
				return m, func() tea.Msg {
					return NavigateMsg{Screen: DeckListScreen}
				}
			case "r":
				// Restart session
				m.session = models.NewStudySession(m.deck, 20)
				m.state = ShowingQuestion
				return m, nil
			}
		}
	}

	return m, nil
}

// rateCardAndContinue rates the current card and moves to the next one
func (m *StudyModel) rateCardAndContinue(rating models.Rating) (tea.Model, tea.Cmd) {
	currentCard := m.session.GetCurrentCard()
	if currentCard == nil {
		return m, nil
	}

	// Update the card with spaced repetition algorithm
	algorithms.UpdateCardReview(currentCard, rating)

	// Update the card in the deck
	deckCard := m.deck.GetCard(currentCard.ID)
	if deckCard != nil {
		*deckCard = *currentCard
		m.deck.MarkModified()
	}

	// Move to next card
	if m.session.NextCard() {
		m.state = ShowingQuestion
		m.selectedRating = 2 // Reset to "Good"
	} else {
		m.state = SessionComplete
	}

	// Save the deck after each card (simple approach for now)
	return m, func() tea.Msg {
		return SaveDeckMsg{m.deck}
	}
}

// View implements tea.Model
func (m *StudyModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	switch m.state {
	case ShowingQuestion:
		return m.viewQuestion()
	case ShowingAnswer:
		return m.viewAnswer()
	case SessionComplete:
		return m.viewSessionComplete()
	default:
		return "Unknown state"
	}
}

// viewQuestion renders the question side of the card
func (m *StudyModel) viewQuestion() string {
	currentCard := m.session.GetCurrentCard()
	if currentCard == nil {
		return "No card to display"
	}

	// Progress indicator
	current, total := m.session.GetProgress()
	progressText := fmt.Sprintf("Card %d of %d", current, total)
	progress := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		Render(progressText)

	// Deck name
	deckName := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(1).
		Render(m.session.DeckName)

	// Card content (question)
	cardContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		PaddingLeft(4).
		PaddingRight(4).
		PaddingTop(2).
		PaddingBottom(2).
		Width(60).
		Height(8).
		Align(lipgloss.Center).
		Foreground(textColor).
		Render(m.wrapText(currentCard.Front, 50))

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(2).
		Render("Press Space or F to flip card • Esc to exit")

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		progress,
		deckName,
		cardContent,
		instructions,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// viewAnswer renders the answer side of the card with rating options
func (m *StudyModel) viewAnswer() string {
	currentCard := m.session.GetCurrentCard()
	if currentCard == nil {
		return "No card to display"
	}

	// Progress indicator
	current, total := m.session.GetProgress()
	progressText := fmt.Sprintf("Card %d of %d", current, total)
	progress := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		Render(progressText)

	// Deck name
	deckName := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(1).
		Render(m.session.DeckName)

	// Card content (question and answer)
	questionText := lipgloss.NewStyle().
		Foreground(mutedColor).
		Bold(true).
		Render("Q: " + currentCard.Front)

	answerText := lipgloss.NewStyle().
		Foreground(textColor).
		Bold(true).
		PaddingTop(1).
		Render("A: " + currentCard.Back)

	cardContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		PaddingLeft(4).
		PaddingRight(4).
		PaddingTop(2).
		PaddingBottom(2).
		Width(60).
		Height(8).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Left, questionText, answerText))

	// Rating buttons
	ratingOptions := []string{"1 Again", "2 Hard", "3 Good", "4 Easy"}
	ratingColors := []lipgloss.Color{errorColor, accentColor, secondaryColor, primaryColor}

	var ratings []string
	for i, option := range ratingOptions {
		style := lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2).
			PaddingTop(1).
			PaddingBottom(1).
			Margin(0, 1)

		if i == m.selectedRating {
			style = style.
				Background(ratingColors[i]).
				Foreground(backgroundColor).
				Bold(true)
		} else {
			style = style.
				Foreground(ratingColors[i]).
				Border(lipgloss.NormalBorder()).
				BorderForeground(ratingColors[i])
		}

		ratings = append(ratings, style.Render(option))
	}

	ratingRow := lipgloss.JoinHorizontal(lipgloss.Center, ratings...)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		PaddingTop(2).
		Render("Use 1-4 keys or ←/→ arrows + Enter to rate • Esc to exit")

	// Combine elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		progress,
		deckName,
		cardContent,
		ratingRow,
		instructions,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// viewSessionComplete renders the session completion screen
func (m *StudyModel) viewSessionComplete() string {
	// Session stats
	completedCards := m.session.CardsStudied
	if !m.session.IsFinished() {
		completedCards++
	}

	title := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render("🎉 Session Complete!")

	statsText := fmt.Sprintf("You studied %d cards from %s", completedCards, m.session.DeckName)
	stats := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center).
		PaddingBottom(2).
		Render(statsText)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Align(lipgloss.Center).
		Render("Press Enter to return to deck list • R to restart session")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		stats,
		instructions,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

// wrapText wraps text to the specified width
func (m *StudyModel) wrapText(text string, width int) string {
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

// SaveDeckMsg is a message to save a deck
type SaveDeckMsg struct {
	Deck *models.Deck
}
