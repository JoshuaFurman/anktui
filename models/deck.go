package models

import (
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cards       []Card    `json:"cards"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// NewDeck creates a new deck with the given name and description
func NewDeck(name, description string) *Deck {
	now := time.Now()
	return &Deck{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Cards:       make([]Card, 0),
		Created:     now,
		Modified:    now,
	}
}

// AddCard adds a new card to the deck
func (d *Deck) AddCard(card *Card) {
	d.Cards = append(d.Cards, *card)
	d.MarkModified()
}

// RemoveCard removes a card by ID from the deck
func (d *Deck) RemoveCard(cardID string) bool {
	for i, card := range d.Cards {
		if card.ID == cardID {
			d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)
			d.MarkModified()
			return true
		}
	}
	return false
}

// GetCard returns a card by ID
func (d *Deck) GetCard(cardID string) *Card {
	for i, card := range d.Cards {
		if card.ID == cardID {
			return &d.Cards[i]
		}
	}
	return nil
}

// GetReviewCards returns all cards that are due for review
func (d *Deck) GetReviewCards() []Card {
	var reviewCards []Card
	for _, card := range d.Cards {
		if card.IsReviewDue() {
			reviewCards = append(reviewCards, card)
		}
	}
	return reviewCards
}

// GetNewCards returns all cards that have never been reviewed
func (d *Deck) GetNewCards() []Card {
	var newCards []Card
	for _, card := range d.Cards {
		if card.Repetition == 0 {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

// GetCardStats returns statistics about the deck
func (d *Deck) GetCardStats() (total, new, review int) {
	total = len(d.Cards)
	for _, card := range d.Cards {
		if card.Repetition == 0 {
			new++
		} else if card.IsReviewDue() {
			review++
		}
	}
	return total, new, review
}

// MarkModified updates the modified timestamp
func (d *Deck) MarkModified() {
	d.Modified = time.Now()
}

// UpdateInfo updates the deck name and description
func (d *Deck) UpdateInfo(name, description string) {
	d.Name = name
	d.Description = description
	d.MarkModified()
}
