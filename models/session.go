package models

import (
	"time"
)

// StudySession represents an active study session for a deck
type StudySession struct {
	DeckID        string    `json:"deck_id"`
	DeckName      string    `json:"deck_name"`
	Cards         []Card    `json:"cards"`
	CurrentIndex  int       `json:"current_index"`
	ShowingAnswer bool      `json:"showing_answer"`
	SessionStart  time.Time `json:"session_start"`
	CardsStudied  int       `json:"cards_studied"`
}

// Rating represents how well the user knew a card
type Rating int

const (
	Again Rating = iota // 0 - Didn't know it, show again soon
	Hard                // 1 - Knew it with difficulty
	Good                // 2 - Knew it well
	Easy                // 3 - Knew it very easily
)

// String returns a human-readable representation of the rating
func (r Rating) String() string {
	switch r {
	case Again:
		return "Again"
	case Hard:
		return "Hard"
	case Good:
		return "Good"
	case Easy:
		return "Easy"
	default:
		return "Unknown"
	}
}

// NewStudySession creates a new study session for the given deck
func NewStudySession(deck *Deck, maxCards int) *StudySession {
	reviewCards := deck.GetReviewCards()
	newCards := deck.GetNewCards()

	// Combine review and new cards, prioritizing review cards
	var sessionCards []Card
	sessionCards = append(sessionCards, reviewCards...)

	// Add new cards up to the limit
	remainingSlots := maxCards - len(reviewCards)
	if remainingSlots > 0 && len(newCards) > 0 {
		newCardsToAdd := remainingSlots
		if newCardsToAdd > len(newCards) {
			newCardsToAdd = len(newCards)
		}
		sessionCards = append(sessionCards, newCards[:newCardsToAdd]...)
	}

	// Limit total cards to maxCards
	if len(sessionCards) > maxCards {
		sessionCards = sessionCards[:maxCards]
	}

	return &StudySession{
		DeckID:        deck.ID,
		DeckName:      deck.Name,
		Cards:         sessionCards,
		CurrentIndex:  0,
		ShowingAnswer: false,
		SessionStart:  time.Now(),
		CardsStudied:  0,
	}
}

// GetCurrentCard returns the current card being studied
func (s *StudySession) GetCurrentCard() *Card {
	if s.CurrentIndex >= len(s.Cards) || s.CurrentIndex < 0 {
		return nil
	}
	return &s.Cards[s.CurrentIndex]
}

// ShowAnswer toggles to showing the answer
func (s *StudySession) ShowAnswer() {
	s.ShowingAnswer = true
}

// NextCard moves to the next card
func (s *StudySession) NextCard() bool {
	if s.CurrentIndex < len(s.Cards)-1 {
		s.CurrentIndex++
		s.ShowingAnswer = false
		s.CardsStudied++
		return true
	}
	return false
}

// IsFinished returns true if all cards have been studied
func (s *StudySession) IsFinished() bool {
	return s.CurrentIndex >= len(s.Cards) || len(s.Cards) == 0
}

// GetProgress returns current progress as (current, total)
func (s *StudySession) GetProgress() (int, int) {
	return s.CurrentIndex + 1, len(s.Cards)
}

// GetRemainingCards returns the number of cards left to study
func (s *StudySession) GetRemainingCards() int {
	remaining := len(s.Cards) - s.CurrentIndex - 1
	if remaining < 0 {
		return 0
	}
	return remaining
}
