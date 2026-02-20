package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/ui"
)

// Menu is the title screen state.
type Menu struct {
	sm          *core.StateManager
	fonts       *Fonts
	newGameBtn  *ui.Button
}

// NewMenu creates the menu state.
func NewMenu(sm *core.StateManager, fonts *Fonts) *Menu {
	return &Menu{sm: sm, fonts: fonts}
}

func (m *Menu) Enter(data any) {
	btnW, btnH := 200.0, 50.0
	btnX := (core.WindowW - btnW) / 2
	btnY := core.WindowH * 0.6

	m.newGameBtn = ui.NewButton("New Game", btnX, btnY, btnW, btnH, m.fonts.Default, func() {
		m.sm.Switch("game_play", nil)
	})
}

func (m *Menu) Exit() {}

func (m *Menu) Update() error {
	mx, my := ebiten.CursorPosition()
	m.newGameBtn.HandleMouseMove(mx, my)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		m.newGameBtn.HandleMousePress(mx, my)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		m.newGameBtn.HandleMouseRelease(mx, my)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		m.sm.Switch("game_play", nil)
	}

	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	// Background
	vector.DrawFilledRect(screen, 0, 0, core.WindowW, core.WindowH, core.ColorBG, false)

	// Title
	title := "SCOUNDREL"
	titleOp := &text.DrawOptions{}
	tw, _ := text.Measure(title, m.fonts.Title, 0)
	titleOp.GeoM.Translate((core.WindowW-tw)/2, core.WindowH*0.25)
	titleOp.ColorScale.ScaleWithColor(core.ColorWhite)
	text.Draw(screen, title, m.fonts.Title, titleOp)

	// Subtitle
	subtitle := "A Scoundrel Card Game"
	subOp := &text.DrawOptions{}
	sw, _ := text.Measure(subtitle, m.fonts.Subtitle, 0)
	subOp.GeoM.Translate((core.WindowW-sw)/2, core.WindowH*0.25+60)
	subOp.ColorScale.ScaleWithColor(core.ColorDim)
	text.Draw(screen, subtitle, m.fonts.Subtitle, subOp)

	// Button
	m.newGameBtn.Draw(screen)
}
