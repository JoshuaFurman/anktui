package storage

import (
	"anktui/config"
	"anktui/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// JSONStorage implements the Storage interface using JSON files
type JSONStorage struct {
	dataDir string
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(cfg *config.Config) (*JSONStorage, error) {
	dataDir, err := cfg.GetExpandedDataDir()
	if err != nil {
		return nil, err
	}

	// Ensure data directory exists
	if err := cfg.EnsureDataDir(); err != nil {
		return nil, err
	}

	return &JSONStorage{dataDir: dataDir}, nil
}

// getDeckFilePath returns the file path for a deck
func (s *JSONStorage) getDeckFilePath(deckID string) string {
	return filepath.Join(s.dataDir, fmt.Sprintf("%s.json", deckID))
}

// SaveDeck saves a deck to a JSON file
func (s *JSONStorage) SaveDeck(deck *models.Deck) error {
	filePath := s.getDeckFilePath(deck.ID)

	// Marshal deck to JSON with proper indentation
	data, err := json.MarshalIndent(deck, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal deck: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write deck file: %w", err)
	}

	return nil
}

// LoadDeck loads a deck by ID from a JSON file
func (s *JSONStorage) LoadDeck(id string) (*models.Deck, error) {
	filePath := s.getDeckFilePath(id)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("deck with ID %s not found", id)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read deck file: %w", err)
	}

	// Unmarshal JSON
	var deck models.Deck
	if err := json.Unmarshal(data, &deck); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deck: %w", err)
	}

	return &deck, nil
}

// LoadAllDecks loads all decks from JSON files
func (s *JSONStorage) LoadAllDecks() ([]*models.Deck, error) {
	// Read all JSON files in data directory
	files, err := filepath.Glob(filepath.Join(s.dataDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list deck files: %w", err)
	}

	var decks []*models.Deck
	for _, file := range files {
		// Extract deck ID from filename
		filename := filepath.Base(file)
		deckID := strings.TrimSuffix(filename, ".json")

		// Load the deck
		deck, err := s.LoadDeck(deckID)
		if err != nil {
			// Log error but continue with other decks
			continue
		}

		decks = append(decks, deck)
	}

	return decks, nil
}

// DeleteDeck deletes a deck by ID
func (s *JSONStorage) DeleteDeck(id string) error {
	filePath := s.getDeckFilePath(id)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("deck with ID %s not found", id)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete deck file: %w", err)
	}

	return nil
}

// DeckExists checks if a deck exists
func (s *JSONStorage) DeckExists(id string) bool {
	filePath := s.getDeckFilePath(id)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// ListDeckIDs returns a list of all deck IDs
func (s *JSONStorage) ListDeckIDs() ([]string, error) {
	// Read all JSON files in data directory
	files, err := filepath.Glob(filepath.Join(s.dataDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list deck files: %w", err)
	}

	var deckIDs []string
	for _, file := range files {
		// Extract deck ID from filename
		filename := filepath.Base(file)
		deckID := strings.TrimSuffix(filename, ".json")
		deckIDs = append(deckIDs, deckID)
	}

	return deckIDs, nil
}
