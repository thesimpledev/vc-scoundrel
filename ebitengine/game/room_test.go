package game

import (
	"testing"

	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
)

func roomWith(slots ...int) (*Room, [4]*entities.Card) {
	var cards [4]*entities.Card
	for _, s := range slots {
		cards[s] = entities.NewCard(core.SuitClubs, "5")
	}
	r := NewRoom()
	r.SetCards(cards)
	return r, cards
}

func TestRoomSetCardsCountsFilledSlots(t *testing.T) {
	r, _ := roomWith(0, 2, 3)
	if r.InitialCount != 3 {
		t.Fatalf("InitialCount = %d, want 3", r.InitialCount)
	}
	if r.CardCount() != 3 {
		t.Fatalf("CardCount = %d, want 3", r.CardCount())
	}
}

func TestRoomGetRemainingCardReturnsFirstFilledSlot(t *testing.T) {
	r, cards := roomWith(1, 3)
	card, slot := r.GetRemainingCard()
	if card != cards[1] || slot != 1 {
		t.Fatalf("got slot %d, want slot 1", slot)
	}
}

func TestRoomGetRemainingCardOnEmptyRoom(t *testing.T) {
	r := NewRoom()
	card, slot := r.GetRemainingCard()
	if card != nil || slot != -1 {
		t.Fatalf("got %v, %d, want nil, -1", card, slot)
	}
}

func TestRoomGetAllCardsSkipsEmptySlots(t *testing.T) {
	r, cards := roomWith(0, 3)
	got := r.GetAllCards()
	if len(got) != 2 || got[0] != cards[0] || got[1] != cards[3] {
		t.Fatalf("GetAllCards returned %v", got)
	}
}

func TestRoomCompletesAfterThreeOfFour(t *testing.T) {
	r, _ := roomWith(0, 1, 2, 3)
	r.ResolveCard(0)
	r.ResolveCard(1)
	if r.IsComplete() {
		t.Fatal("room complete after two cards")
	}
	r.ResolveCard(2)
	if !r.IsComplete() {
		t.Fatal("room not complete after three cards")
	}
}
