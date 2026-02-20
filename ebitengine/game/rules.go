package game

import "github.com/thesimpledev/radianthorizon/entities"

// Pure validation functions — no side effects.

// CanUsePotion returns true if the player hasn't used a potion this room yet.
func CanUsePotion(room *Room) bool {
	return !room.PotionUsed
}

// CanUseWeapon returns true if the weapon can fight this monster.
func CanUseWeapon(ws *entities.WeaponSlot, card *entities.Card) bool {
	if !ws.HasWeapon() {
		return false
	}
	return ws.CanUseAgainst(card.Value)
}

// CanAvoidRoom returns true if the player didn't avoid the last room.
func CanAvoidRoom(avoidedLastRoom bool) bool {
	return !avoidedLastRoom
}

// IsGameOver returns true if HP <= 0 or dungeon + room are both empty.
func IsGameOver(hp int, deck *entities.Deck, room *Room) bool {
	if hp <= 0 {
		return true
	}
	if deck.IsEmpty() && room.CardCount() == 0 {
		return true
	}
	return false
}

// IsVictory returns true if the player survived the entire dungeon.
func IsVictory(hp int, deck *entities.Deck, room *Room) bool {
	return hp > 0 && deck.IsEmpty() && room.CardCount() == 0
}

// CardType returns "monster", "weapon", "potion", or "unknown".
func CardType(card *entities.Card) string {
	if card.IsMonster {
		return "monster"
	}
	if card.IsWeapon {
		return "weapon"
	}
	if card.IsPotion {
		return "potion"
	}
	return "unknown"
}
