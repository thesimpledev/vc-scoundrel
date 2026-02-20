package entities

// WeaponSlot tracks the equipped weapon and the stack of monsters slain by it.
type WeaponSlot struct {
	Weapon *Card   // the equipped Diamond card, or nil
	Slain  []*Card // monsters killed by this weapon
}

// NewWeaponSlot creates an empty weapon slot.
func NewWeaponSlot() *WeaponSlot {
	return &WeaponSlot{}
}

// Equip replaces the weapon. Returns the old weapon and old slain stack for discarding.
func (ws *WeaponSlot) Equip(weaponCard *Card) (oldWeapon *Card, oldSlain []*Card) {
	oldWeapon = ws.Weapon
	oldSlain = ws.Slain
	ws.Weapon = weaponCard
	ws.Slain = nil
	return
}

// CanUseAgainst checks if the weapon can fight a monster of the given value.
// First kill: always allowed. Subsequent: monster value must be <= last slain.
func (ws *WeaponSlot) CanUseAgainst(monsterValue int) bool {
	if ws.Weapon == nil {
		return false
	}
	if len(ws.Slain) == 0 {
		return true
	}
	lastSlainValue := ws.Slain[len(ws.Slain)-1].Value
	return monsterValue <= lastSlainValue
}

// RecordKill adds a monster to the slain stack.
func (ws *WeaponSlot) RecordKill(monsterCard *Card) {
	ws.Slain = append(ws.Slain, monsterCard)
}

// CalculateDamage returns damage taken when using this weapon against a monster.
func (ws *WeaponSlot) CalculateDamage(monsterValue int) int {
	if ws.Weapon == nil {
		return monsterValue
	}
	dmg := monsterValue - ws.Weapon.Value
	if dmg < 0 {
		return 0
	}
	return dmg
}

// HasWeapon returns true if a weapon is equipped.
func (ws *WeaponSlot) HasWeapon() bool {
	return ws.Weapon != nil
}

// WeaponValue returns the weapon's value (0 if no weapon).
func (ws *WeaponSlot) WeaponValue() int {
	if ws.Weapon != nil {
		return ws.Weapon.Value
	}
	return 0
}

// Clear empties the weapon slot.
func (ws *WeaponSlot) Clear() {
	ws.Weapon = nil
	ws.Slain = nil
}
