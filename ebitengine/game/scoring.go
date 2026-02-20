package game

import "github.com/thesimpledev/radianthorizon/entities"

// Calculate computes the final score.
func ScoreCalculate(hp int, deck *entities.Deck, lastCard *entities.Card) int {
	if hp <= 0 {
		return -deck.CountRemainingMonsterValue()
	}

	score := hp
	if lastCard != nil && lastCard.IsPotion {
		score += lastCard.Value
	}
	return score
}
