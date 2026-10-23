package game

import (
	"testing"

	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
	"github.com/thesimpledev/radianthorizon/lib"
)

func TestDealRoomFillsFourSlotsFromDeck(t *testing.T) {
	lib.EventClear()
	tm := NewTurnManager()
	before := tm.Deck.Remaining()

	slots, carry := tm.DealRoom()

	if carry != -1 {
		t.Fatalf("carry slot = %d, want -1", carry)
	}
	for i, c := range slots {
		if c == nil {
			t.Fatalf("slot %d is empty", i)
		}
	}
	if tm.Deck.Remaining() != before-4 {
		t.Fatalf("deck remaining = %d, want %d", tm.Deck.Remaining(), before-4)
	}
	if tm.Phase != "resolving" {
		t.Fatalf("phase = %q, want resolving", tm.Phase)
	}
}

func TestDealRoomKeepsCarryOverInItsSlot(t *testing.T) {
	lib.EventClear()
	tm := NewTurnManager()
	carried := entities.NewCard(core.SuitHearts, "3")
	tm.CarryOver = carried
	tm.CarryOverSlot = 2

	slots, carry := tm.DealRoom()

	if carry != 2 || slots[2] != carried {
		t.Fatalf("carry slot = %d, want 2 holding the carried card", carry)
	}
	if tm.CarryOver != nil || tm.CarryOverSlot != -1 {
		t.Fatal("carry-over not cleared after dealing")
	}
}

func TestDealRoomLeavesSlotsEmptyWhenDeckRunsOut(t *testing.T) {
	lib.EventClear()
	tm := NewTurnManager()
	tm.Deck.Cards = tm.Deck.Cards[:2]

	slots, _ := tm.DealRoom()

	if slots[0] == nil || slots[1] == nil || slots[2] != nil || slots[3] != nil {
		t.Fatalf("slots = %v, want first two filled", slots)
	}
}

func TestResolveWeaponReturnsOldWeaponAndSlain(t *testing.T) {
	lib.EventClear()
	tm := NewTurnManager()
	oldWeapon := entities.NewCard(core.SuitDiamonds, "4")
	slain := entities.NewCard(core.SuitClubs, "3")
	tm.WeaponSlot.Equip(oldWeapon)
	tm.WeaponSlot.RecordKill(slain)
	newWeapon := entities.NewCard(core.SuitDiamonds, "7")
	tm.Room.SetCards([4]*entities.Card{newWeapon})
	tm.Phase = "resolving"

	result := tm.ResolveCard(0, false)

	if result.Action != "weapon_equipped" || result.Card != newWeapon || result.Slot != 0 {
		t.Fatalf("unexpected result %+v", result)
	}
	if result.OldWeapon != oldWeapon || len(result.OldSlain) != 1 || result.OldSlain[0] != slain {
		t.Fatalf("old weapon and slain not returned: %+v", result)
	}
}
