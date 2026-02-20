package game

// Pure damage calculation. No graphics — testable in isolation.

// BarehandedDamage returns damage from fighting a monster barehanded.
func BarehandedDamage(monsterValue int) int {
	return monsterValue
}

// WeaponDamage returns damage from fighting a monster with a weapon.
func WeaponDamage(monsterValue, weaponValue int) int {
	dmg := monsterValue - weaponValue
	if dmg < 0 {
		return 0
	}
	return dmg
}

// ApplyDamage reduces HP by damage, clamped to 0.
func ApplyDamage(currentHP, damage int) int {
	hp := currentHP - damage
	if hp < 0 {
		return 0
	}
	return hp
}

// ApplyHealing increases HP by a potion value, clamped to max.
func ApplyHealing(currentHP, potionValue, maxHP int) int {
	hp := currentHP + potionValue
	if hp > maxHP {
		return maxHP
	}
	return hp
}
