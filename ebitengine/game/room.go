package game

import "github.com/thesimpledev/radianthorizon/entities"

// Room tracks the 4 card slots and per-turn rules.
type Room struct {
	Slots         [4]*entities.Card // nil = empty
	ResolvedCount int
	PotionUsed    bool
	InitialCount  int
}

// NewRoom creates an empty room.
func NewRoom() *Room {
	return &Room{}
}

// SetCards sets up the room with cards at specific slot indices.
// slotCards is keyed by index 0-3.
func (r *Room) SetCards(slotCards [4]*entities.Card) {
	r.Slots = slotCards
	r.ResolvedCount = 0
	r.PotionUsed = false
	count := 0
	for i := 0; i < 4; i++ {
		if r.Slots[i] != nil {
			count++
		}
	}
	r.InitialCount = count
}

// ResolveCard marks a card as taken. Returns the card.
func (r *Room) ResolveCard(index int) *entities.Card {
	card := r.Slots[index]
	r.Slots[index] = nil
	r.ResolvedCount++
	return card
}

// IsComplete checks if this room is finished.
func (r *Room) IsComplete() bool {
	if r.InitialCount >= 4 {
		return r.ResolvedCount >= 3
	} else if r.InitialCount > 1 {
		return r.ResolvedCount >= r.InitialCount-1
	}
	return r.ResolvedCount >= 1
}

// GetRemainingCard returns the carry-over card and its slot index.
// Returns nil, -1 if empty.
func (r *Room) GetRemainingCard() (*entities.Card, int) {
	for i := 0; i < 4; i++ {
		if r.Slots[i] != nil {
			return r.Slots[i], i
		}
	}
	return nil, -1
}

// GetAllCards returns all non-nil cards currently in the room.
func (r *Room) GetAllCards() []*entities.Card {
	var cards []*entities.Card
	for i := 0; i < 4; i++ {
		if r.Slots[i] != nil {
			cards = append(cards, r.Slots[i])
		}
	}
	return cards
}

// CardCount returns how many cards are in the room.
func (r *Room) CardCount() int {
	count := 0
	for i := 0; i < 4; i++ {
		if r.Slots[i] != nil {
			count++
		}
	}
	return count
}

// HasCardAt checks if a slot is occupied.
func (r *Room) HasCardAt(index int) bool {
	return index >= 0 && index < 4 && r.Slots[index] != nil
}

// Clear empties the room.
func (r *Room) Clear() {
	r.Slots = [4]*entities.Card{}
	r.ResolvedCount = 0
	r.PotionUsed = false
	r.InitialCount = 0
}
