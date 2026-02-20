package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/ui"
)

// GameOver is the win/loss screen state.
type GameOver struct {
	sm    *core.StateManager
	fonts *Fonts

	victory      bool
	score        int
	hp           int
	playAgainBtn *ui.Button
	menuBtn      *ui.Button
}

// NewGameOver creates the game over state.
func NewGameOver(sm *core.StateManager, fonts *Fonts) *GameOver {
	return &GameOver{sm: sm, fonts: fonts}
}

func (go_ *GameOver) Enter(data any) {
	if d, ok := data.(map[string]any); ok {
		if v, ok := d["victory"].(bool); ok {
			go_.victory = v
		}
		if s, ok := d["score"].(int); ok {
			go_.score = s
		}
		if h, ok := d["hp"].(int); ok {
			go_.hp = h
		}
	}

	btnW, btnH := 200.0, 50.0
	btnX := (core.WindowW - btnW) / 2
	btnY := core.WindowH * 0.65

	go_.playAgainBtn = ui.NewButton("Play Again", btnX, btnY, btnW, btnH, go_.fonts.Default, func() {
		go_.sm.Switch("game_play", nil)
	})

	menuY := btnY + btnH + 16
	go_.menuBtn = ui.NewButton("Main Menu", btnX, menuY, btnW, btnH, go_.fonts.Default, func() {
		go_.sm.Switch("menu", nil)
	})
}

func (go_ *GameOver) Exit() {}

func (go_ *GameOver) Update() error {
	mx, my := ebiten.CursorPosition()
	go_.playAgainBtn.HandleMouseMove(mx, my)
	go_.menuBtn.HandleMouseMove(mx, my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		go_.playAgainBtn.HandleMousePress(mx, my)
		go_.menuBtn.HandleMousePress(mx, my)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		go_.playAgainBtn.HandleMouseRelease(mx, my)
		go_.menuBtn.HandleMouseRelease(mx, my)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		go_.sm.Switch("game_play", nil)
	}

	return nil
}

func (go_ *GameOver) Draw(screen *ebiten.Image) {
	// Background
	vector.DrawFilledRect(screen, 0, 0, core.WindowW, core.WindowH, core.ColorBG, false)

	// Title
	var title string
	var titleColor = core.ColorRed
	if go_.victory {
		titleColor = core.ColorGold
		title = "VICTORY!"
	} else {
		title = "DEFEATED"
	}
	titleOp := &text.DrawOptions{}
	tw, _ := text.Measure(title, go_.fonts.Title, 0)
	titleOp.GeoM.Translate((core.WindowW-tw)/2, core.WindowH*0.25)
	titleOp.ColorScale.ScaleWithColor(titleColor)
	text.Draw(screen, title, go_.fonts.Title, titleOp)

	// Score
	scoreText := fmt.Sprintf("Score: %d", go_.score)
	scoreOp := &text.DrawOptions{}
	sw, _ := text.Measure(scoreText, go_.fonts.Score, 0)
	scoreOp.GeoM.Translate((core.WindowW-sw)/2, core.WindowH*0.38)
	scoreOp.ColorScale.ScaleWithColor(core.ColorWhite)
	text.Draw(screen, scoreText, go_.fonts.Score, scoreOp)

	// HP
	hpText := fmt.Sprintf("Final HP: %d", go_.hp)
	hpOp := &text.DrawOptions{}
	hw, _ := text.Measure(hpText, go_.fonts.Score, 0)
	hpOp.GeoM.Translate((core.WindowW-hw)/2, core.WindowH*0.38+36)
	hpOp.ColorScale.ScaleWithColor(core.ColorWhite)
	text.Draw(screen, hpText, go_.fonts.Score, hpOp)

	// Buttons
	go_.playAgainBtn.Draw(screen)
	go_.menuBtn.Draw(screen)
}
