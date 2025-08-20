package algorithms

import (
	"anktui/models"
	"math"
	"time"
)

// UpdateCardReview updates a card based on the user's rating using the SM-2 algorithm
// This is based on the SuperMemo 2 algorithm for spaced repetition
func UpdateCardReview(card *models.Card, rating models.Rating) {
	now := time.Now()
	card.LastReview = now
	card.MarkModified()

	switch rating {
	case models.Again:
		// Reset the card - user didn't know it
		card.Repetition = 0
		card.Interval = 1
		card.EaseFactor = math.Max(1.3, card.EaseFactor-0.2) // Decrease ease factor but don't go below 1.3

	case models.Hard:
		// Increase repetition count
		card.Repetition++

		// Calculate new interval
		if card.Repetition == 1 {
			card.Interval = 1
		} else if card.Repetition == 2 {
			card.Interval = 6
		} else {
			// For repetition >= 3: Interval(n) = Interval(n-1) * EaseFactor
			card.Interval = int(float64(card.Interval) * card.EaseFactor)
		}

		// Decrease ease factor for hard cards
		card.EaseFactor = math.Max(1.3, card.EaseFactor-0.15)

	case models.Good:
		// This is the standard path - user knew the card normally
		card.Repetition++

		// Calculate new interval using standard SM-2 formula
		if card.Repetition == 1 {
			card.Interval = 1
		} else if card.Repetition == 2 {
			card.Interval = 6
		} else {
			// For repetition >= 3: Interval(n) = Interval(n-1) * EaseFactor
			card.Interval = int(float64(card.Interval) * card.EaseFactor)
		}

		// Ease factor remains unchanged for "Good" rating

	case models.Easy:
		// User found it very easy
		card.Repetition++

		// Calculate new interval with bonus for easy cards
		if card.Repetition == 1 {
			card.Interval = 4 // Start with longer interval for easy cards
		} else if card.Repetition == 2 {
			card.Interval = 8
		} else {
			// For repetition >= 3: Interval(n) = Interval(n-1) * EaseFactor * 1.3 (bonus for easy)
			card.Interval = int(float64(card.Interval) * card.EaseFactor * 1.3)
		}

		// Increase ease factor for easy cards
		card.EaseFactor = card.EaseFactor + 0.1
	}

	// Set the next review date
	card.NextReview = now.AddDate(0, 0, card.Interval)
}

// GetDueCards returns cards that are due for review from a slice of cards
func GetDueCards(cards []models.Card) []models.Card {
	var dueCards []models.Card
	for _, card := range cards {
		if card.IsReviewDue() {
			dueCards = append(dueCards, card)
		}
	}
	return dueCards
}

// GetNewCards returns cards that have never been studied (repetition == 0)
func GetNewCards(cards []models.Card) []models.Card {
	var newCards []models.Card
	for _, card := range cards {
		if card.Repetition == 0 {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

// CalculateRetentionStats calculates retention statistics for a set of cards
func CalculateRetentionStats(cards []models.Card) (total, mature, young, new int) {
	total = len(cards)

	for _, card := range cards {
		if card.Repetition == 0 {
			new++
		} else if card.Interval >= 21 { // Cards with interval >= 21 days are considered mature
			mature++
		} else {
			young++
		}
	}

	return total, mature, young, new
}
