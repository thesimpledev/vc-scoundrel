package states

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
	"github.com/thesimpledev/radianthorizon/game"
	"github.com/thesimpledev/radianthorizon/lib"
	"github.com/thesimpledev/radianthorizon/ui"
)

type fakeInput struct {
	x, y     int
	pressed  bool
	released bool
	keys     map[ebiten.Key]bool
}

func (f *fakeInput) CursorPosition() (int, int)         { return f.x, f.y }
func (f *fakeInput) MouseJustPressed() bool             { return f.pressed }
func (f *fakeInput) MouseJustReleased() bool            { return f.released }
func (f *fakeInput) KeyJustPressed(key ebiten.Key) bool { return f.keys[key] }

type gameplayFixture struct {
	gp     *GamePlay
	states *testStates
	input  *fakeInput
	clock  time.Time
	sounds []string
}

func newGameplayFixture(t *testing.T) *gameplayFixture {
	t.Helper()
	lib.AnimCancelAll()
	lib.EventClear()
	t.Cleanup(func() {
		lib.AnimCancelAll()
		lib.EventClear()
	})

	f := &gameplayFixture{
		states: newTestStates(),
		input:  &fakeInput{keys: map[ebiten.Key]bool{}},
		clock:  time.Unix(1000, 0),
	}
	fonts := testFonts(t)
	gp := NewGamePlay(f.states.sm, fonts)
	gp.input = f.input
	gp.now = func() time.Time { return f.clock }
	gp.playSFX = func(name string) { f.sounds = append(f.sounds, name) }
	gp.turn = game.NewTurnManager()
	gp.hud = ui.NewHUD(fonts.Default, fonts.Default)
	gp.combatSlot = -1
	gp.lastUpdate = f.clock
	f.gp = gp
	return f
}

func slotX(slot int) float64 {
	return core.RoomStartX + float64(slot)*core.RoomSpacing
}

// setRoom puts cards face up and clickable in the room, both in the rules and on screen.
func (f *gameplayFixture) setRoom(cards ...*entities.Card) {
	var slots [4]*entities.Card
	for i, c := range cards {
		if c == nil {
			continue
		}
		c.X, c.Y = slotX(i), core.RoomY
		c.SetFaceUp(true)
		c.Clickable = true
		slots[i] = c
	}
	f.gp.turn.Room.SetCards(slots)
	f.gp.turn.Phase = "resolving"
	f.gp.roomCards = slots
}

func (f *gameplayFixture) pointAtSlot(slot int) {
	f.input.x, f.input.y = int(slotX(slot))+5, int(core.RoomY)+5
}

func (f *gameplayFixture) pointAtButton(b *ui.Button) {
	f.input.x, f.input.y = int(b.X)+5, int(b.Y)+5
}

func (f *gameplayFixture) update(t *testing.T) {
	t.Helper()
	if err := f.gp.Update(); err != nil {
		t.Fatal(err)
	}
	f.input.pressed, f.input.released = false, false
	f.input.keys = map[ebiten.Key]bool{}
}

func runAnimations() {
	for range 20 {
		lib.AnimUpdateAll(1)
	}
}

func monster(rank string) *entities.Card { return entities.NewCard(core.SuitClubs, rank) }
func potion(rank string) *entities.Card  { return entities.NewCard(core.SuitHearts, rank) }
func weapon(rank string) *entities.Card  { return entities.NewCard(core.SuitDiamonds, rank) }

func TestUpdateAdvancesAnimationsByElapsedTime(t *testing.T) {
	f := newGameplayFixture(t)
	value := 0.0
	lib.AnimTo([]lib.AnimGoal{{Field: &value, Goal: 1}}, 1, "linear", nil)

	f.clock = f.clock.Add(500 * time.Millisecond)
	f.update(t)

	if value != 0.5 {
		t.Fatalf("value = %v, want 0.5", value)
	}
}

func TestUpdateWaitsForGameOverTimer(t *testing.T) {
	f := newGameplayFixture(t)
	f.gp.gameOverPending, f.gp.gameOverTimer = true, 1

	f.clock = f.clock.Add(400 * time.Millisecond)
	f.update(t)

	if len(f.states.gameOver.entered) != 0 || !f.gp.gameOverPending {
		t.Fatal("switched to game over before the timer ran out")
	}
}

func TestUpdateSwitchesToGameOverWhenTimerRunsOut(t *testing.T) {
	f := newGameplayFixture(t)
	data := map[string]any{"score": 5}
	f.gp.gameOverPending, f.gp.gameOverTimer, f.gp.gameOverData = true, 1, data
	f.input.keys[ebiten.KeyQ] = true

	f.clock = f.clock.Add(time.Second)
	f.update(t)

	if len(f.states.gameOver.entered) != 1 || f.states.gameOver.entered[0].(map[string]any)["score"] != 5 {
		t.Fatal("did not switch to game over with its data")
	}
	if f.gp.gameOverPending {
		t.Fatal("game over still pending")
	}
	if len(f.states.menu.entered) != 0 {
		t.Fatal("input was handled in the frame that switched to game over")
	}
}

func TestUpdateHoverMarksCardUnderCursor(t *testing.T) {
	f := newGameplayFixture(t)
	first, second := monster("5"), monster("6")
	f.setRoom(first, second)

	f.pointAtSlot(1)
	f.update(t)

	if f.gp.hoveredCard != second || !second.Hovered || first.Hovered {
		t.Fatal("hover not set on the card under the cursor")
	}

	f.input.x, f.input.y = 0, 0
	f.update(t)
	if f.gp.hoveredCard != nil || second.Hovered {
		t.Fatal("hover not cleared when the cursor left the card")
	}
}

func TestUpdateClearsHoverWhileAnimating(t *testing.T) {
	f := newGameplayFixture(t)
	card := monster("5")
	f.setRoom(card)
	f.gp.hoveredCard = card
	f.gp.animating = true

	f.pointAtSlot(0)
	f.update(t)

	if f.gp.hoveredCard != nil {
		t.Fatal("hovered card kept while animating")
	}
}

func TestUpdateClickOnPotionResolvesIt(t *testing.T) {
	f := newGameplayFixture(t)
	f.gp.turn.HP = 10
	heal := potion("4")
	f.setRoom(heal, monster("5"), monster("6"), monster("7"))

	f.pointAtSlot(0)
	f.input.pressed = true
	f.update(t)

	if f.gp.roomCards[0] != nil || f.gp.turn.HP != 14 || !f.gp.animating {
		t.Fatalf("potion not resolved: HP %d, animating %v", f.gp.turn.HP, f.gp.animating)
	}
	runAnimations()
	if f.gp.discardTop != heal || f.gp.animating {
		t.Fatal("potion did not land on the discard pile")
	}
	if !f.gp.roomCards[1].Clickable {
		t.Fatal("room not unlocked after the potion landed")
	}
}

func TestUpdateIgnoresClicksWhileAnimating(t *testing.T) {
	f := newGameplayFixture(t)
	f.setRoom(potion("4"))
	f.gp.animating = true

	f.pointAtSlot(0)
	f.input.pressed = true
	f.update(t)

	if f.gp.roomCards[0] == nil {
		t.Fatal("card resolved while animating")
	}
}

func (f *gameplayFixture) openCombatChoice(t *testing.T) *entities.Card {
	t.Helper()
	f.gp.turn.WeaponSlot.Equip(weapon("5"))
	enemy := monster("3")
	f.setRoom(enemy, monster("9"), monster("10"), potion("2"))
	f.pointAtSlot(0)
	f.input.pressed = true
	f.update(t)
	return enemy
}

func TestUpdateClickOnMonsterWithWeaponShowsCombatChoice(t *testing.T) {
	f := newGameplayFixture(t)
	f.openCombatChoice(t)

	if !f.gp.choosingCombat || f.gp.combatSlot != 0 || len(f.gp.combatButtons) != 2 {
		t.Fatal("combat choice not shown")
	}
	if f.gp.combatButtons[0].Text != "Barehanded (-3 HP)" || f.gp.combatButtons[1].Text != "Use Weapon (-0 HP)" {
		t.Fatalf("button labels %q, %q", f.gp.combatButtons[0].Text, f.gp.combatButtons[1].Text)
	}
	if f.gp.roomCards[1].Clickable {
		t.Fatal("room still clickable during combat choice")
	}
}

func TestUpdateHoversCombatButtonsWhileChoosing(t *testing.T) {
	f := newGameplayFixture(t)
	f.openCombatChoice(t)
	weaponBtn := f.gp.combatButtons[1]

	f.pointAtButton(weaponBtn)
	f.update(t)

	if !weaponBtn.Hovered || f.gp.combatButtons[0].Hovered {
		t.Fatal("combat button hover not updated")
	}
}

func TestUpdateReleaseOnWeaponButtonFightsWithWeapon(t *testing.T) {
	f := newGameplayFixture(t)
	enemy := f.openCombatChoice(t)
	weaponBtn := f.gp.combatButtons[1]

	f.pointAtButton(weaponBtn)
	f.input.pressed = true
	f.update(t)
	f.input.released = true
	f.update(t)

	if f.gp.choosingCombat || f.gp.combatButtons != nil || f.gp.combatSlot != -1 {
		t.Fatal("combat choice not closed")
	}
	slain := f.gp.turn.WeaponSlot.Slain
	if len(slain) != 1 || slain[0] != enemy || f.gp.turn.HP != core.StartHP {
		t.Fatalf("monster not fought with the weapon: slain %v, HP %d", slain, f.gp.turn.HP)
	}
}

func TestUpdateEscapeClosesCombatChoice(t *testing.T) {
	f := newGameplayFixture(t)
	f.openCombatChoice(t)

	f.input.keys[ebiten.KeyEscape] = true
	f.update(t)

	if f.gp.choosingCombat || f.gp.combatButtons != nil || f.gp.combatSlot != -1 {
		t.Fatal("combat choice not closed")
	}
	if !f.gp.roomCards[1].Clickable {
		t.Fatal("room not clickable after closing the combat choice")
	}
}

func TestUpdateSkipKeyAvoidsRoom(t *testing.T) {
	f := newGameplayFixture(t)
	f.setRoom(monster("2"), monster("3"), monster("4"), monster("5"))
	deckBefore := f.gp.turn.Deck.Remaining()

	f.input.keys[ebiten.KeyS] = true
	f.update(t)

	if !f.gp.animating {
		t.Fatal("skip did not start")
	}
	runAnimations()
	if !f.gp.turn.AvoidedLastRoom || f.gp.turn.Deck.Remaining() != deckBefore {
		t.Fatalf("room not avoided and redealt: deck %d, want %d", f.gp.turn.Deck.Remaining(), deckBefore)
	}
}

func TestUpdateSkipKeyIgnoredWhenRoomCannotBeAvoided(t *testing.T) {
	f := newGameplayFixture(t)
	f.setRoom(monster("2"), monster("3"), monster("4"), monster("5"))
	f.gp.turn.AvoidedLastRoom = true

	f.input.keys[ebiten.KeyS] = true
	f.update(t)

	if f.gp.animating {
		t.Fatal("skipped a room that cannot be avoided")
	}
}

func TestUpdateQuitKeyReturnsToMenu(t *testing.T) {
	f := newGameplayFixture(t)
	value := 0.0
	lib.AnimTo([]lib.AnimGoal{{Field: &value, Goal: 1}}, 10, "linear", nil)

	f.input.keys[ebiten.KeyQ] = true
	f.update(t)

	if len(f.states.menu.entered) != 1 {
		t.Fatal("did not return to the menu")
	}
	if lib.AnimActiveCount() != 0 {
		t.Fatal("animations not cancelled on quit")
	}
}

func TestUpdateIgnoresKeysWhileAnimating(t *testing.T) {
	f := newGameplayFixture(t)
	f.gp.animating = true

	f.input.keys[ebiten.KeyQ] = true
	f.update(t)

	if len(f.states.menu.entered) != 0 {
		t.Fatal("quit handled while animating")
	}
}

func TestDoDealAnimatesNewCardsIntoRoom(t *testing.T) {
	f := newGameplayFixture(t)

	f.gp.doDeal()

	if !f.gp.animating {
		t.Fatal("not animating while dealing")
	}
	for slot, card := range f.gp.roomCards {
		if card == nil {
			t.Fatalf("slot %d not dealt", slot)
		}
	}
	runAnimations()
	if f.gp.animating {
		t.Fatal("still animating after the deal finished")
	}
	for slot, card := range f.gp.roomCards {
		if card.X != slotX(slot) || card.Y != core.RoomY || !card.FaceUp || !card.Clickable {
			t.Fatalf("slot %d card not in place: %+v", slot, card)
		}
	}
	if len(f.sounds) != 4 {
		t.Fatalf("played %d deal sounds, want 4", len(f.sounds))
	}
}

func TestDoDealKeepsCarryOverInPlace(t *testing.T) {
	f := newGameplayFixture(t)
	carried := monster("7")
	carried.X, carried.Y = slotX(2), core.RoomY
	carried.SetFaceUp(true)
	f.gp.turn.CarryOver, f.gp.turn.CarryOverSlot = carried, 2

	f.gp.doDeal()

	if carried.X != slotX(2) || !carried.FaceUp {
		t.Fatal("carry-over card was moved back to the dungeon")
	}
	runAnimations()
	if f.gp.roomCards[2] != carried || !carried.FaceUp || len(f.sounds) != 3 {
		t.Fatalf("carry-over flipped or wrong sound count %d", len(f.sounds))
	}
}

func TestDoDealWithNothingToDealUnlocksImmediately(t *testing.T) {
	f := newGameplayFixture(t)
	f.gp.turn.Deck.Cards = nil

	f.gp.doDeal()

	if f.gp.animating {
		t.Fatal("animating with nothing to deal")
	}
}

func TestDrawRendersWithoutPanicking(t *testing.T) {
	f := newGameplayFixture(t)
	f.gp.turn.Deck.Cards = nil
	hidden := monster("5")
	f.setRoom(hidden)
	hidden.Visible = false
	f.gp.hoveredCard = hidden

	f.gp.Draw(ebiten.NewImage(core.WindowW, core.WindowH))
}
