package storage

import (
	"anktui/models"
)

// Storage interface defines the methods for persisting and retrieving deck data
type Storage interface {
	// SaveDeck saves a deck to storage
	SaveDeck(deck *models.Deck) error

	// LoadDeck loads a deck by ID from storage
	LoadDeck(id string) (*models.Deck, error)

	// LoadAllDecks loads all decks from storage
	LoadAllDecks() ([]*models.Deck, error)

	// DeleteDeck deletes a deck by ID from storage
	DeleteDeck(id string) error

	// DeckExists checks if a deck exists in storage
	DeckExists(id string) bool

	// ListDeckIDs returns a list of all deck IDs in storage
	ListDeckIDs() ([]string, error)
}
