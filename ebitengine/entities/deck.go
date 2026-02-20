package entities

import (
	"math/rand/v2"

	"github.com/thesimpledev/radianthorizon/core"
)

// Deck is the 44-card Scoundrel dungeon deck.
type Deck struct {
	Cards []*Card
}

// NewDeck builds and shuffles a 44-card Scoundrel deck.
func NewDeck() *Deck {
	d := &Deck{}
	d.build()
	d.Shuffle()
	return d
}

// build creates the 44-card deck according to Scoundrel rules.
func (d *Deck) build() {
	d.Cards = nil

	// Clubs and Spades: all 13 ranks (monsters)
	for _, suit := range []string{core.SuitClubs, core.SuitSpades} {
		for _, rank := range core.Ranks {
			d.Cards = append(d.Cards, NewCard(suit, rank))
		}
	}

	// Diamonds and Hearts: only 2-10 (no face cards or aces)
	redRanks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10"}
	for _, suit := range []string{core.SuitDiamonds, core.SuitHearts} {
		for _, rank := range redRanks {
			d.Cards = append(d.Cards, NewCard(suit, rank))
		}
	}
}

// Shuffle randomizes the deck using Fisher-Yates.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// DrawCard removes and returns the top card. Returns nil if empty.
func (d *Deck) DrawCard() *Card {
	if len(d.Cards) == 0 {
		return nil
	}
	n := len(d.Cards)
	card := d.Cards[n-1]
	d.Cards = d.Cards[:n-1]
	return card
}

// PlaceAtBottom puts cards at the bottom of the deck (for room avoidance).
func (d *Deck) PlaceAtBottom(cards []*Card) {
	d.Cards = append(cards, d.Cards...)
}

// Remaining returns how many cards are left.
func (d *Deck) Remaining() int {
	return len(d.Cards)
}

// IsEmpty returns true if the deck has no cards.
func (d *Deck) IsEmpty() bool {
	return len(d.Cards) == 0
}

// CountRemainingMonsterValue sums the value of all remaining monster cards.
func (d *Deck) CountRemainingMonsterValue() int {
	total := 0
	for _, card := range d.Cards {
		if card.IsMonster {
			total += card.Value
		}
	}
	return total
}
