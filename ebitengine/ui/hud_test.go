package ui

import (
	"bytes"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
	"github.com/thesimpledev/radianthorizon/game"
	"golang.org/x/image/font/gofont/goregular"
)

func testFont(t *testing.T) *text.GoTextFace {
	t.Helper()
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatal(err)
	}
	return &text.GoTextFace{Source: src, Size: 16}
}

func weaponSlot(weapon string, slain ...string) *entities.WeaponSlot {
	ws := entities.NewWeaponSlot()
	if weapon != "" {
		ws.Equip(entities.NewCard(core.SuitDiamonds, weapon))
	}
	for _, rank := range slain {
		ws.RecordKill(entities.NewCard(core.SuitClubs, rank))
	}
	return ws
}

func TestMonsterInfoText(t *testing.T) {
	h := &HUD{}
	monster := entities.NewCard(core.SuitClubs, "8")
	tests := []struct {
		name string
		ws   *entities.WeaponSlot
		want string
	}{
		{"no weapon", weaponSlot(""), "Barehanded: -8 HP"},
		{"weapon reduces damage", weaponSlot("5"), "Barehanded: -8 HP    Weapon: -3 HP"},
		{"weapon blocks all damage", weaponSlot("10"), "Barehanded: -8 HP    Weapon: no damage!"},
		{"weapon blocked by slain stack", weaponSlot("5", "4"), "Barehanded: -8 HP    Weapon: blocked (monster 8 > last slain 4)"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line1, line2 := h.monsterInfo(monster, tc.ws)
			if line1 != "Monster  |  Deals 8 damage" {
				t.Fatalf("line1 = %q", line1)
			}
			if line2 != tc.want {
				t.Fatalf("line2 = %q, want %q", line2, tc.want)
			}
		})
	}
}

func TestPotionInfoText(t *testing.T) {
	h := &HUD{}
	potion := entities.NewCard(core.SuitHearts, "5")
	tests := []struct {
		name string
		used bool
		hp   int
		want string
	}{
		{"already used", true, 10, "Already used a potion this room - this one will be wasted!"},
		{"full heal", false, 10, "Would heal 5 HP  (10 -> 15)"},
		{"capped heal", false, core.MaxHP - 2, "Would heal 2 HP (capped at max 20)"},
		{"hp full", false, core.MaxHP, "HP is already full - healing would be wasted"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			room := game.NewRoom()
			room.PotionUsed = tc.used
			_, line2 := h.potionInfo(potion, room, tc.hp)
			if line2 != tc.want {
				t.Fatalf("line2 = %q, want %q", line2, tc.want)
			}
		})
	}
}

func TestHUDDrawHandlesEveryCardType(t *testing.T) {
	font := testFont(t)
	h := NewHUD(font, font)
	screen := ebiten.NewImage(core.WindowW, core.WindowH)
	cards := []*entities.Card{
		nil,
		entities.NewCard(core.SuitClubs, "8"),
		entities.NewCard(core.SuitDiamonds, "4"),
		entities.NewCard(core.SuitHearts, "3"),
	}
	for _, card := range cards {
		h.Draw(screen, HUDState{
			HP:            12,
			MaxHP:         core.MaxHP,
			DeckRemaining: 30,
			Weapon:        weaponSlot("5", "9"),
			CanAvoid:      card == nil,
			HoveredCard:   card,
			Room:          game.NewRoom(),
		})
	}
}
