package game

import (
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
	"github.com/thesimpledev/radianthorizon/lib"
)

// ResolveResult describes what happened when a card was resolved.
type ResolveResult struct {
	Action     string
	Card       *entities.Card
	Slot       int
	Damage     int
	Healed     int
	Wasted     bool
	UsedWeapon bool
	OldWeapon  *entities.Card
	OldSlain   []*entities.Card
}

// TurnManager orchestrates the turn cycle.
type TurnManager struct {
	Deck       *entities.Deck
	Room       *Room
	WeaponSlot *entities.WeaponSlot

	HP              int
	AvoidedLastRoom bool
	LastCard        *entities.Card

	CarryOver     *entities.Card
	CarryOverSlot int // -1 if none

	GameOver bool
	Victory  bool
	Score    int

	Phase string // "ready", "dealing", "resolving", "room_complete"
}

// NewTurnManager creates a new game instance.
func NewTurnManager() *TurnManager {
	return &TurnManager{
		Deck:          entities.NewDeck(),
		Room:          NewRoom(),
		WeaponSlot:    entities.NewWeaponSlot(),
		HP:            core.StartHP,
		CarryOverSlot: -1,
		Phase:         "ready",
	}
}

// DealRoom deals cards into the room.
// Returns the slot cards array and the carry-over slot index (-1 if none).
func (tm *TurnManager) DealRoom() ([4]*entities.Card, int) {
	tm.Phase = "dealing"
	var slotCards [4]*entities.Card
	coSlot := -1

	// Carry-over stays in its original slot
	if tm.CarryOver != nil {
		coSlot = tm.CarryOverSlot
		slotCards[coSlot] = tm.CarryOver
		tm.CarryOver = nil
		tm.CarryOverSlot = -1
	}

	// Fill empty slots from the dungeon deck
	for i := 0; i < 4; i++ {
		if slotCards[i] == nil && !tm.Deck.IsEmpty() {
			slotCards[i] = tm.Deck.DrawCard()
		}
	}

	tm.Room.SetCards(slotCards)
	tm.Phase = "resolving"

	lib.EventEmit("room_dealt", map[string]any{
		"slot_cards":      slotCards,
		"carry_over_slot": coSlot,
	})
	return slotCards, coSlot
}

// CanAvoid returns true if the player can skip the current room.
func (tm *TurnManager) CanAvoid() bool {
	return CanAvoidRoom(tm.AvoidedLastRoom) &&
		tm.Room.CardCount() == 4 &&
		tm.Phase == "resolving" &&
		tm.Room.ResolvedCount == 0
}

// AvoidRoom sends all 4 room cards to the bottom of the deck.
func (tm *TurnManager) AvoidRoom() bool {
	if !tm.CanAvoid() {
		return false
	}

	cards := tm.Room.GetAllCards()
	tm.Room.Clear()
	tm.Deck.PlaceAtBottom(cards)
	tm.AvoidedLastRoom = true

	lib.EventEmit("room_avoided", map[string]any{
		"cards": cards,
	})
	return true
}

// ResolveCard resolves a card at the given slot index.
func (tm *TurnManager) ResolveCard(slotIndex int, useWeapon bool) *ResolveResult {
	if tm.Phase != "resolving" {
		return nil
	}
	if !tm.Room.HasCardAt(slotIndex) {
		return nil
	}

	card := tm.Room.Slots[slotIndex]
	var result *ResolveResult

	if card.IsPotion {
		result = tm.resolvePotion(slotIndex, card)
	} else if card.IsWeapon {
		result = tm.resolveWeapon(slotIndex, card)
	} else if card.IsMonster {
		result = tm.resolveMonster(slotIndex, card, useWeapon)
	}

	if tm.Room.IsComplete() {
		tm.finishRoom()
	}

	if tm.HP <= 0 {
		tm.endGame()
	}

	return result
}

func (tm *TurnManager) resolvePotion(slotIndex int, card *entities.Card) *ResolveResult {
	healed := 0
	wasted := false

	if CanUsePotion(tm.Room) {
		oldHP := tm.HP
		tm.HP = ApplyHealing(tm.HP, card.Value, core.MaxHP)
		healed = tm.HP - oldHP
		tm.Room.PotionUsed = true
	} else {
		wasted = true
	}

	tm.Room.ResolveCard(slotIndex)
	tm.LastCard = card

	result := &ResolveResult{
		Action: "potion",
		Card:   card,
		Slot:   slotIndex,
		Healed: healed,
		Wasted: wasted,
	}
	lib.EventEmit("potion_used", result)
	return result
}

func (tm *TurnManager) resolveWeapon(slotIndex int, card *entities.Card) *ResolveResult {
	oldWeapon, oldSlain := tm.WeaponSlot.Equip(card)
	tm.Room.ResolveCard(slotIndex)
	tm.LastCard = card

	result := &ResolveResult{
		Action:   "weapon_equipped",
		Card:     card,
		Slot:     slotIndex,
		OldWeapon: oldWeapon,
		OldSlain:  oldSlain,
	}
	lib.EventEmit("weapon_equipped", result)
	return result
}

func (tm *TurnManager) resolveMonster(slotIndex int, card *entities.Card, useWeapon bool) *ResolveResult {
	var damage int
	usedWeapon := false

	if useWeapon && CanUseWeapon(tm.WeaponSlot, card) {
		damage = WeaponDamage(card.Value, tm.WeaponSlot.WeaponValue())
		tm.WeaponSlot.RecordKill(card)
		usedWeapon = true
	} else {
		damage = BarehandedDamage(card.Value)
	}

	tm.HP = ApplyDamage(tm.HP, damage)
	tm.Room.ResolveCard(slotIndex)
	tm.LastCard = card

	result := &ResolveResult{
		Action:     "monster_fought",
		Card:       card,
		Slot:       slotIndex,
		Damage:     damage,
		UsedWeapon: usedWeapon,
	}
	lib.EventEmit("monster_fought", result)
	return result
}

func (tm *TurnManager) finishRoom() {
	carryOver, coSlot := tm.Room.GetRemainingCard()
	tm.CarryOver = carryOver
	tm.CarryOverSlot = coSlot
	tm.Room.Clear()
	tm.AvoidedLastRoom = false

	tm.Phase = "room_complete"
	lib.EventEmit("room_complete", map[string]any{
		"carry_over":      tm.CarryOver,
		"carry_over_slot": tm.CarryOverSlot,
	})

	if tm.Deck.IsEmpty() && tm.CarryOver == nil {
		tm.endGame()
	}
}

// EndGame ends the game and calculates the score.
func (tm *TurnManager) EndGame() {
	tm.endGame()
}

func (tm *TurnManager) endGame() {
	if tm.GameOver {
		return
	}
	tm.GameOver = true
	tm.Victory = tm.HP > 0
	tm.Score = ScoreCalculate(tm.HP, tm.Deck, tm.LastCard)

	lib.EventEmit("game_over", map[string]any{
		"victory": tm.Victory,
		"score":   tm.Score,
		"hp":      tm.HP,
	})
}

// GetActions returns available actions for a card at the given slot.
func (tm *TurnManager) GetActions(slotIndex int) []string {
	if tm.Phase != "resolving" {
		return nil
	}
	if !tm.Room.HasCardAt(slotIndex) {
		return nil
	}

	card := tm.Room.Slots[slotIndex]
	var actions []string

	if card.IsPotion {
		actions = append(actions, "use_potion")
	} else if card.IsWeapon {
		actions = append(actions, "equip_weapon")
	} else if card.IsMonster {
		actions = append(actions, "fight_barehanded")
		if CanUseWeapon(tm.WeaponSlot, card) {
			actions = append(actions, "fight_with_weapon")
		}
	}
	return actions
}
