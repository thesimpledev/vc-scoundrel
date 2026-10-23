package entities

import "testing"

func cardCounts(cards []*Card) map[string]int {
	counts := map[string]int{}
	for _, c := range cards {
		counts[c.String()]++
	}
	return counts
}

func TestNewDeckHasFortyFourUniqueCards(t *testing.T) {
	d := NewDeck()
	if d.Remaining() != 44 {
		t.Fatalf("deck size = %d, want 44", d.Remaining())
	}
	for name, n := range cardCounts(d.Cards) {
		if n != 1 {
			t.Fatalf("%s appears %d times", name, n)
		}
	}
}

func TestShuffleKeepsTheSameCards(t *testing.T) {
	d := NewDeck()
	before := cardCounts(d.Cards)

	d.Shuffle()

	after := cardCounts(d.Cards)
	if len(after) != len(before) {
		t.Fatalf("shuffle changed the card set: %d before, %d after", len(before), len(after))
	}
	for name := range before {
		if after[name] != 1 {
			t.Fatalf("%s missing after shuffle", name)
		}
	}
}
