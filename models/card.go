package models

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID       string    `json:"id"`
	Front    string    `json:"front"`
	Back     string    `json:"back"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`

	// Spaced repetition data
	Interval   int       `json:"interval"`    // Days until next review
	Repetition int       `json:"repetition"`  // Number of successful reviews
	EaseFactor float64   `json:"ease_factor"` // Difficulty multiplier (default 2.5)
	NextReview time.Time `json:"next_review"`
	LastReview time.Time `json:"last_review,omitempty"`
}

// NewCard creates a new flashcard with default values
func NewCard(front, back string) *Card {
	now := time.Now()
	return &Card{
		ID:         uuid.New().String(),
		Front:      front,
		Back:       back,
		Created:    now,
		Modified:   now,
		Interval:   1,
		Repetition: 0,
		EaseFactor: 2.5,
		NextReview: now, // Available for immediate review
	}
}

// IsReviewDue checks if the card is ready for review
func (c *Card) IsReviewDue() bool {
	return time.Now().After(c.NextReview) || time.Now().Equal(c.NextReview)
}

// MarkModified updates the modified timestamp
func (c *Card) MarkModified() {
	c.Modified = time.Now()
}

// UpdateContent updates the front and back content of the card
func (c *Card) UpdateContent(front, back string) {
	c.Front = front
	c.Back = back
	c.MarkModified()
}
