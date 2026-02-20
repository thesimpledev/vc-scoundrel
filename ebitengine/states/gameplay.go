package states

import (
	"fmt"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/entities"
	game "github.com/thesimpledev/radianthorizon/game"
	"github.com/thesimpledev/radianthorizon/lib"
	"github.com/thesimpledev/radianthorizon/ui"
)

// GamePlay is the main gameplay state.
type GamePlay struct {
	sm    *core.StateManager
	fonts *Fonts

	turn *game.TurnManager

	// Visual card tracking
	roomCards   [4]*entities.Card
	weaponCard  *entities.Card
	slainCards  []*entities.Card
	discardTop  *entities.Card
	flyingCards []*entities.Card

	// Combat choice UI
	combatButtons []*ui.Button
	choosingCombat bool
	combatSlot     int

	// Animation lock
	animating bool

	// Hovered card for info bar
	hoveredCard *entities.Card

	// Game over transition
	gameOverPending bool
	gameOverTimer   float64
	gameOverData    map[string]any

	// HUD
	hud *ui.HUD
}

// NewGamePlay creates the gameplay state.
func NewGamePlay(sm *core.StateManager, fonts *Fonts) *GamePlay {
	return &GamePlay{sm: sm, fonts: fonts}
}

func (gp *GamePlay) Enter(data any) {
	lib.AnimCancelAll()
	lib.EventClear()

	gp.turn = game.NewTurnManager()

	gp.roomCards = [4]*entities.Card{}
	gp.weaponCard = nil
	gp.slainCards = nil
	gp.discardTop = nil
	gp.flyingCards = nil
	gp.combatButtons = nil
	gp.choosingCombat = false
	gp.combatSlot = -1
	gp.animating = false
	gp.hoveredCard = nil
	gp.gameOverPending = false
	gp.gameOverTimer = 0
	gp.gameOverData = nil

	gp.hud = ui.NewHUD(gp.fonts.Default, gp.fonts.Default)

	gp.setupEvents()
	gp.doDeal()
}

func (gp *GamePlay) Exit() {
	lib.AnimCancelAll()
	lib.EventClear()
}

// === Event Subscriptions ===

func (gp *GamePlay) setupEvents() {
	lib.EventOn("monster_fought", func(data any) {
		result := data.(*game.ResolveResult)
		if result.Damage > 0 {
			core.PlaySFX("audio/casino/card-shove-1")
		} else {
			core.PlaySFX("audio/casino/card-place-2")
		}
	})

	lib.EventOn("potion_used", func(data any) {
		result := data.(*game.ResolveResult)
		if result.Wasted {
			core.PlaySFX("audio/casino/card-shove-2")
		} else {
			core.PlaySFX("audio/casino/card-place-1")
		}
	})

	lib.EventOn("weapon_equipped", func(data any) {
		core.PlaySFX("audio/casino/card-place-3")
	})

	lib.EventOn("room_avoided", func(data any) {
		core.PlaySFX("audio/casino/card-shuffle")
	})

	lib.EventOn("game_over", func(data any) {
		d := data.(map[string]any)
		gp.gameOverPending = true
		gp.gameOverTimer = 1.0
		gp.gameOverData = d
	})
}

// === Deal Animation ===

func (gp *GamePlay) doDeal() {
	gp.animating = true
	slotCards, carrySlot := gp.turn.DealRoom()

	// Find new cards that need animation (everything except carry-over)
	var newCardSlots []int
	for i := 0; i < 4; i++ {
		if slotCards[i] != nil && i != carrySlot {
			newCardSlots = append(newCardSlots, i)
		}
	}

	pendingAnims := len(newCardSlots)

	if pendingAnims == 0 {
		gp.roomCards = slotCards
		gp.animating = false
		gp.enableRoomInteraction()
		return
	}

	// Place all cards into visual slots
	gp.roomCards = slotCards

	// Animate new cards from dungeon pile
	for staggerI, slot := range newCardSlots {
		card := slotCards[slot]
		card.X = core.DungeonX
		card.Y = core.DungeonY
		card.SetFaceUp(false)
		card.Visible = true

		targetX := core.RoomStartX + float64(slot)*core.RoomSpacing
		delay := float64(staggerI) * core.DealStagger

		// Capture slot and card in closure
		capturedCard := card
		gp.delayedCall(delay, func() {
			core.PlaySFX(fmt.Sprintf("audio/casino/card-slide-%d", rand.IntN(8)+1))
			capturedCard.MoveTo(targetX, core.RoomY, core.MoveDuration, "ease_out_quad", func() {
				capturedCard.Flip(func() {
					pendingAnims--
					if pendingAnims <= 0 {
						gp.animating = false
						gp.enableRoomInteraction()
					}
				})
			})
		})
	}
}

func (gp *GamePlay) delayedCall(delay float64, fn func()) {
	if delay <= 0 {
		fn()
		return
	}
	dummy := 0.0
	lib.AnimTo([]lib.AnimGoal{
		{Field: &dummy, Goal: 1},
	}, delay, "linear", fn)
}

// === Card Interaction ===

func (gp *GamePlay) enableRoomInteraction() {
	for i := 0; i < 4; i++ {
		if gp.roomCards[i] != nil {
			gp.roomCards[i].Clickable = true
		}
	}
}

func (gp *GamePlay) disableRoomInteraction() {
	for i := 0; i < 4; i++ {
		if gp.roomCards[i] != nil {
			gp.roomCards[i].Clickable = false
			gp.roomCards[i].Hovered = false
		}
	}
}

func (gp *GamePlay) onCardClicked(slotIndex int) {
	if gp.animating || gp.choosingCombat {
		return
	}

	card := gp.roomCards[slotIndex]
	if card == nil {
		return
	}

	actions := gp.turn.GetActions(slotIndex)
	if len(actions) == 0 {
		return
	}

	// Monster with weapon option: show combat choice
	if card.IsMonster && len(actions) > 1 {
		gp.showCombatChoice(slotIndex, actions)
		return
	}

	useWeapon := actions[0] == "fight_with_weapon"
	gp.doResolve(slotIndex, useWeapon)
}

// === Combat Choice UI ===

func (gp *GamePlay) showCombatChoice(slotIndex int, actions []string) {
	gp.choosingCombat = true
	gp.combatSlot = slotIndex
	gp.disableRoomInteraction()

	card := gp.roomCards[slotIndex]
	btnW, btnH := 180.0, 40.0
	btnX := card.X + core.CardDrawW + 12
	btnY := card.Y

	gp.combatButtons = nil

	// Barehanded button
	bareDmg := card.Value
	gp.combatButtons = append(gp.combatButtons, ui.NewButton(
		fmt.Sprintf("Barehanded (-%d HP)", bareDmg),
		btnX, btnY, btnW, btnH, gp.fonts.Default,
		func() { gp.chooseCombat(false) },
	))

	// Weapon button if available
	for _, action := range actions {
		if action == "fight_with_weapon" {
			weapDmg := card.Value - gp.turn.WeaponSlot.WeaponValue()
			if weapDmg < 0 {
				weapDmg = 0
			}
			gp.combatButtons = append(gp.combatButtons, ui.NewButton(
				fmt.Sprintf("Use Weapon (-%d HP)", weapDmg),
				btnX, btnY+btnH+8, btnW, btnH, gp.fonts.Default,
				func() { gp.chooseCombat(true) },
			))
		}
	}
}

func (gp *GamePlay) chooseCombat(useWeapon bool) {
	gp.choosingCombat = false
	gp.combatButtons = nil
	gp.doResolve(gp.combatSlot, useWeapon)
	gp.combatSlot = -1
}

func (gp *GamePlay) hideCombatChoice() {
	gp.choosingCombat = false
	gp.combatButtons = nil
	gp.combatSlot = -1
	gp.enableRoomInteraction()
}

// === Resolve a Card ===

func (gp *GamePlay) doResolve(slotIndex int, useWeapon bool) {
	gp.animating = true
	gp.disableRoomInteraction()

	card := gp.roomCards[slotIndex]
	result := gp.turn.ResolveCard(slotIndex, useWeapon)
	if result == nil {
		gp.animating = false
		gp.enableRoomInteraction()
		return
	}

	card.Clickable = false
	card.Hovered = false

	switch result.Action {
	case "potion":
		gp.animateToDiscard(card, slotIndex, func() {
			gp.afterResolve()
		})
	case "weapon_equipped":
		gp.animateEquipWeapon(card, slotIndex, result, func() {
			gp.afterResolve()
		})
	case "monster_fought":
		if result.UsedWeapon {
			gp.animateToSlain(card, slotIndex, func() {
				gp.afterResolve()
			})
		} else {
			gp.animateToDiscard(card, slotIndex, func() {
				gp.afterResolve()
			})
		}
	}
}

func (gp *GamePlay) afterResolve() {
	if gp.turn.GameOver {
		gp.animating = false
		return
	}

	if gp.turn.Room.IsComplete() || gp.turn.Phase == "room_complete" {
		gp.handleRoomComplete()
	} else {
		gp.animating = false
		gp.enableRoomInteraction()
	}
}

func (gp *GamePlay) handleRoomComplete() {
	hasCarry := gp.turn.CarryOver != nil

	if !hasCarry {
		gp.roomCards = [4]*entities.Card{}
	}

	if !gp.turn.Deck.IsEmpty() || hasCarry {
		if !gp.turn.GameOver {
			gp.delayedCall(0.3, func() {
				gp.doDeal()
			})
		} else {
			gp.animating = false
		}
	} else {
		gp.animating = false
		if !gp.turn.GameOver {
			gp.turn.EndGame()
		}
	}
}

// === Card Animations ===

func (gp *GamePlay) startFlying(card *entities.Card) {
	gp.flyingCards = append(gp.flyingCards, card)
}

func (gp *GamePlay) stopFlying(card *entities.Card) {
	for i, c := range gp.flyingCards {
		if c == card {
			gp.flyingCards = append(gp.flyingCards[:i], gp.flyingCards[i+1:]...)
			return
		}
	}
}

func (gp *GamePlay) animateToDiscard(card *entities.Card, slotIndex int, onComplete func()) {
	gp.roomCards[slotIndex] = nil
	gp.startFlying(card)
	card.MoveTo(core.DiscardX, core.DiscardY, core.MoveDuration, "ease_in_out_quad", func() {
		gp.stopFlying(card)
		gp.discardTop = card
		card.SetFaceUp(false)
		if onComplete != nil {
			onComplete()
		}
	})
}

func (gp *GamePlay) animateEquipWeapon(card *entities.Card, slotIndex int, result *game.ResolveResult, onComplete func()) {
	gp.roomCards[slotIndex] = nil
	gp.startFlying(card)

	// Discard old weapon + slain stack
	if result.OldWeapon != nil {
		oldCards := []*entities.Card{result.OldWeapon}
		oldCards = append(oldCards, result.OldSlain...)

		discardPending := len(oldCards)
		for _, oldCard := range oldCards {
			oc := oldCard // capture
			gp.startFlying(oc)
			oc.MoveTo(core.DiscardX, core.DiscardY, core.MoveDuration, "ease_in_out_quad", func() {
				gp.stopFlying(oc)
				oc.SetFaceUp(false)
				discardPending--
				if discardPending <= 0 {
					gp.discardTop = oc
				}
			})
		}
		gp.weaponCard = nil
		gp.slainCards = nil
	}

	// Move new weapon to weapon slot
	card.MoveTo(core.WeaponX, core.WeaponY, core.MoveDuration, "ease_out_back", func() {
		gp.stopFlying(card)
		gp.weaponCard = card
		if onComplete != nil {
			onComplete()
		}
	})
}

func (gp *GamePlay) animateToSlain(card *entities.Card, slotIndex int, onComplete func()) {
	gp.roomCards[slotIndex] = nil
	slainIndex := len(gp.slainCards)
	gp.slainCards = append(gp.slainCards, card)

	targetY := core.SlainY + float64(slainIndex)*core.SlainOffsetY
	card.MoveTo(core.SlainX, targetY, core.MoveDuration, "ease_out_quad", func() {
		if onComplete != nil {
			onComplete()
		}
	})
}

// === Update ===

func (gp *GamePlay) Update() error {
	// Get dt from TPS
	dt := 1.0 / float64(ebiten.TPS())
	lib.AnimUpdateAll(dt)

	// Game over delay timer
	if gp.gameOverPending {
		gp.gameOverTimer -= dt
		if gp.gameOverTimer <= 0 {
			gp.gameOverPending = false
			gp.sm.Switch("game_over", gp.gameOverData)
			return nil
		}
	}

	// Input handling
	mx, my := ebiten.CursorPosition()

	// Mouse move
	if gp.choosingCombat {
		for _, btn := range gp.combatButtons {
			btn.HandleMouseMove(mx, my)
		}
	} else if !gp.animating {
		gp.hoveredCard = nil
		for i := 0; i < 4; i++ {
			card := gp.roomCards[i]
			if card != nil && card.Clickable {
				over := card.ContainsPoint(mx, my)
				card.Hovered = over
				if over {
					gp.hoveredCard = card
				}
			}
		}
	} else {
		gp.hoveredCard = nil
	}

	// Mouse press
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if gp.choosingCombat {
			for _, btn := range gp.combatButtons {
				btn.HandleMousePress(mx, my)
			}
		} else if !gp.animating {
			for i := 0; i < 4; i++ {
				card := gp.roomCards[i]
				if card != nil && card.Clickable && card.ContainsPoint(mx, my) {
					gp.onCardClicked(i)
					break
				}
			}
		}
	}

	// Mouse release
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if gp.choosingCombat {
			for _, btn := range gp.combatButtons {
				btn.HandleMouseRelease(mx, my)
			}
		}
	}

	// Keyboard
	if !gp.animating {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) && gp.turn.CanAvoid() {
			gp.doAvoidRoom()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && gp.choosingCombat {
			gp.hideCombatChoice()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			lib.AnimCancelAll()
			lib.EventClear()
			gp.sm.Switch("menu", nil)
		}
	}

	return nil
}

// === Draw ===

func (gp *GamePlay) Draw(screen *ebiten.Image) {
	// Background
	vector.DrawFilledRect(screen, 0, 0, core.WindowW, core.WindowH, core.ColorBG, false)

	// Dungeon pile
	if gp.turn.Deck.Remaining() > 0 {
		backImg := core.LoadImage("cards/" + core.CardBack)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(core.CardScale, core.CardScale)
		op.GeoM.Translate(core.DungeonX, core.DungeonY)
		screen.DrawImage(backImg, op)
	}

	// Discard pile
	if gp.discardTop != nil {
		gp.discardTop.Draw(screen)
	}

	// Weapon card
	if gp.weaponCard != nil {
		gp.weaponCard.Draw(screen)
	}

	// Slain monster stack
	for _, card := range gp.slainCards {
		card.Draw(screen)
	}

	// Room cards
	for i := 0; i < 4; i++ {
		if gp.roomCards[i] != nil {
			gp.roomCards[i].Draw(screen)
		}
	}

	// Flying cards (on top)
	for _, card := range gp.flyingCards {
		card.Draw(screen)
	}

	// Combat choice buttons
	for _, btn := range gp.combatButtons {
		btn.Draw(screen)
	}

	// HUD
	gp.hud.Draw(screen,
		gp.turn.HP,
		core.MaxHP,
		gp.turn.Deck.Remaining(),
		gp.turn.WeaponSlot,
		gp.turn.CanAvoid(),
		gp.turn.Phase,
		gp.hoveredCard,
		gp.turn.Room,
	)
}

// === Avoid Room ===

func (gp *GamePlay) doAvoidRoom() {
	gp.animating = true
	gp.disableRoomInteraction()

	var cardsToMove []*entities.Card
	for i := 0; i < 4; i++ {
		if gp.roomCards[i] != nil {
			cardsToMove = append(cardsToMove, gp.roomCards[i])
		}
	}

	pending := len(cardsToMove)
	for _, card := range cardsToMove {
		c := card // capture
		c.Flip(func() {
			c.MoveTo(core.DungeonX, core.DungeonY, core.MoveDuration, "ease_in_out_quad", func() {
				c.Visible = false
				pending--
				if pending <= 0 {
					gp.roomCards = [4]*entities.Card{}
					gp.turn.AvoidRoom()
					gp.doDeal()
				}
			})
		})
	}
}
