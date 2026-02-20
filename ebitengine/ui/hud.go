package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
	"github.com/thesimpledev/radianthorizon/game"
)

const (
	infoBarH = 60
	infoBarY = core.WindowH - infoBarH
)

// HUD draws the heads-up display: HP bar, deck count, weapon info, skip indicator, card info.
type HUD struct {
	Font      *text.GoTextFace
	SmallFont *text.GoTextFace
}

// NewHUD creates the HUD with fonts.
func NewHUD(font, smallFont *text.GoTextFace) *HUD {
	return &HUD{Font: font, SmallFont: smallFont}
}

// Draw renders the full HUD.
func (h *HUD) Draw(screen *ebiten.Image, hp, maxHP, deckRemaining int, ws *entities.WeaponSlot, canAvoid bool, phase string, hoveredCard *entities.Card, room *game.Room) {
	h.drawHPBar(screen, hp, maxHP)
	h.drawDeckCount(screen, deckRemaining)
	h.drawWeaponInfo(screen, ws)
	h.drawSkipIndicator(screen, canAvoid)
	h.drawCardInfo(screen, hoveredCard, ws, room, hp)
}

func (h *HUD) drawHPBar(screen *ebiten.Image, hp, maxHP int) {
	x, y := float32(20), float32(20)
	w, h2 := float32(200), float32(24)
	fill := (float32(hp) / float32(maxHP)) * w

	// Background
	vector.DrawFilledRect(screen, x, y, w, h2, core.ColorHPBG, true)

	// Fill
	var fillColor color.NRGBA
	if hp <= 5 {
		fillColor = core.ColorRed
	} else {
		fillColor = core.ColorHPBar
	}
	vector.DrawFilledRect(screen, x, y, fill, h2, fillColor, true)

	// Border
	vector.StrokeRect(screen, x, y, w, h2, 1, core.ColorWhite, true)

	// Text
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x)+8, float64(y)+4)
	op.ColorScale.ScaleWithColor(core.ColorWhite)
	text.Draw(screen, fmt.Sprintf("HP: %d / %d", hp, maxHP), h.Font, op)
}

func (h *HUD) drawDeckCount(screen *ebiten.Image, remaining int) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(core.DungeonX, core.DungeonY-22)
	op.ColorScale.ScaleWithColor(core.ColorWhite)
	text.Draw(screen, fmt.Sprintf("Dungeon: %d", remaining), h.Font, op)
}

func (h *HUD) drawWeaponInfo(screen *ebiten.Image, ws *entities.WeaponSlot) {
	x := float64(core.WeaponX)
	y := float64(core.WeaponY) - 22

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(core.ColorWhite)

	if ws.HasWeapon() {
		weapon := ws.Weapon
		text.Draw(screen, fmt.Sprintf("Weapon: %s of %s", weapon.Rank, weapon.Suit), h.Font, op)

		if len(ws.Slain) > 0 {
			last := ws.Slain[len(ws.Slain)-1]
			op2 := &text.DrawOptions{}
			slainY := y + core.CardDrawH + float64(len(ws.Slain))*core.SlainOffsetY + 10
			op2.GeoM.Translate(x, slainY)
			op2.ColorScale.ScaleWithColor(core.ColorWhite)
			text.Draw(screen, fmt.Sprintf("Last slain: %s (can fight \u2264%d)", last.Rank, last.Value), h.Font, op2)
		}
	} else {
		text.Draw(screen, "No weapon equipped", h.Font, op)
	}
}

func (h *HUD) drawSkipIndicator(screen *ebiten.Image, canAvoid bool) {
	x := float64(core.WindowW - 220)
	y := float64(20)

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)

	if canAvoid {
		op.ColorScale.ScaleWithColor(core.ColorGold)
		text.Draw(screen, "[ Press S to skip room ]", h.Font, op)
	} else {
		op.ColorScale.ScaleWithColor(color.NRGBA{153, 153, 153, 153})
		text.Draw(screen, "Cannot skip this room", h.Font, op)
	}
}

func (h *HUD) drawCardInfo(screen *ebiten.Image, card *entities.Card, ws *entities.WeaponSlot, room *game.Room, hp int) {
	// Background bar
	vector.DrawFilledRect(screen, 0, float32(infoBarY), core.WindowW, infoBarH, color.NRGBA{0, 0, 0, 153}, true)
	vector.StrokeRect(screen, 0, float32(infoBarY), core.WindowW, infoBarH, 1, color.NRGBA{102, 102, 102, 128}, true)

	if card == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(20, float64(infoBarY)+8)
		op.ColorScale.ScaleWithColor(color.NRGBA{128, 128, 128, 128})
		text.Draw(screen, "Hover over a card to see what it does", h.Font, op)
		return
	}

	name := card.Rank + " of " + card.Suit
	var line1, line2 string

	if card.IsMonster {
		line1, line2 = h.monsterInfo(card, ws)
	} else if card.IsWeapon {
		line1, line2 = h.weaponInfoText(card, ws)
	} else if card.IsPotion {
		line1, line2 = h.potionInfo(card, room, hp)
	}

	// Card name with type color
	nameOp := &text.DrawOptions{}
	nameOp.GeoM.Translate(20, float64(infoBarY)+8)
	if card.IsMonster {
		nameOp.ColorScale.ScaleWithColor(core.ColorRed)
	} else if card.IsWeapon {
		nameOp.ColorScale.ScaleWithColor(color.NRGBA{153, 179, 255, 255})
	} else if card.IsPotion {
		nameOp.ColorScale.ScaleWithColor(core.ColorHPBar)
	}
	text.Draw(screen, name, h.Font, nameOp)

	// Description line 1
	nameW, _ := text.Measure(name, h.Font, 0)
	if line1 != "" {
		op := &text.DrawOptions{}
		op.GeoM.Translate(20+nameW+20, float64(infoBarY)+8)
		op.ColorScale.ScaleWithColor(core.ColorWhite)
		text.Draw(screen, line1, h.Font, op)
	}

	// Description line 2
	if line2 != "" {
		op := &text.DrawOptions{}
		op.GeoM.Translate(20, float64(infoBarY)+30)
		op.ColorScale.ScaleWithColor(color.NRGBA{204, 204, 204, 204})
		text.Draw(screen, line2, h.Font, op)
	}
}

func (h *HUD) monsterInfo(card *entities.Card, ws *entities.WeaponSlot) (string, string) {
	line1 := fmt.Sprintf("Monster  |  Deals %d damage", card.Value)
	parts := fmt.Sprintf("Barehanded: -%d HP", card.Value)

	if ws.HasWeapon() {
		if ws.CanUseAgainst(card.Value) {
			weapDmg := card.Value - ws.WeaponValue()
			if weapDmg < 0 {
				weapDmg = 0
			}
			if weapDmg == 0 {
				parts += "    Weapon: no damage!"
			} else {
				parts += fmt.Sprintf("    Weapon: -%d HP", weapDmg)
			}
		} else {
			last := ws.Slain[len(ws.Slain)-1]
			parts += fmt.Sprintf("    Weapon: blocked (monster %d > last slain %d)", card.Value, last.Value)
		}
	}

	return line1, parts
}

func (h *HUD) weaponInfoText(card *entities.Card, ws *entities.WeaponSlot) (string, string) {
	line1 := fmt.Sprintf("Weapon  |  Absorbs up to %d damage per fight", card.Value)
	var line2 string

	if ws.HasWeapon() {
		old := ws.Weapon
		if card.Value > old.Value {
			line2 = fmt.Sprintf("Replaces current weapon (%s of %s) - upgrade!", old.Rank, old.Suit)
		} else if card.Value < old.Value {
			line2 = fmt.Sprintf("Replaces current weapon (%s of %s) - downgrade", old.Rank, old.Suit)
		} else {
			line2 = "Replaces current weapon (same strength, resets slain stack)"
		}
	} else {
		line2 = "You have no weapon - equipping this lets you fight monsters with reduced damage"
	}

	return line1, line2
}

func (h *HUD) potionInfo(card *entities.Card, room *game.Room, hp int) (string, string) {
	line1 := fmt.Sprintf("Health Potion  |  Heals %d HP (max %d)", card.Value, core.MaxHP)
	var line2 string

	if room.PotionUsed {
		line2 = "Already used a potion this room - this one will be wasted!"
	} else {
		actualHeal := card.Value
		if core.MaxHP-hp < actualHeal {
			actualHeal = core.MaxHP - hp
		}
		if actualHeal <= 0 {
			line2 = "HP is already full - healing would be wasted"
		} else if actualHeal < card.Value {
			line2 = fmt.Sprintf("Would heal %d HP (capped at max %d)", actualHeal, core.MaxHP)
		} else {
			line2 = fmt.Sprintf("Would heal %d HP  (%d -> %d)", actualHeal, hp, hp+actualHeal)
		}
	}

	return line1, line2
}
