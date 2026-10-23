package states

import (
	"bytes"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/thesimpledev/radianthorizon/core"
	"golang.org/x/image/font/gofont/goregular"
)

type recordingState struct {
	entered []any
}

func (r *recordingState) Enter(data any)            { r.entered = append(r.entered, data) }
func (r *recordingState) Exit()                     {}
func (r *recordingState) Update() error             { return nil }
func (r *recordingState) Draw(screen *ebiten.Image) {}

type testStates struct {
	sm       *core.StateManager
	menu     *recordingState
	gamePlay *recordingState
	gameOver *recordingState
}

func newTestStates() *testStates {
	ts := &testStates{
		sm:       core.NewStateManager(),
		menu:     &recordingState{},
		gamePlay: &recordingState{},
		gameOver: &recordingState{},
	}
	ts.sm.Register("menu", ts.menu)
	ts.sm.Register("game_play", ts.gamePlay)
	ts.sm.Register("game_over", ts.gameOver)
	return ts
}

func testFonts(t *testing.T) *Fonts {
	t.Helper()
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		t.Fatal(err)
	}
	face := &text.GoTextFace{Source: src, Size: 16}
	return &Fonts{Title: face, Subtitle: face, Score: face, Default: face}
}

func TestMenuEnterPlacesNewGameButton(t *testing.T) {
	ts := newTestStates()
	m := NewMenu(ts.sm, testFonts(t))

	m.Enter(nil)

	b := m.newGameBtn
	if b.Text != "New Game" || b.X != (core.WindowW-200.0)/2 || b.Y != core.WindowH*0.6 || b.W != 200 || b.H != 50 {
		t.Fatalf("unexpected button %+v", b)
	}
	b.OnClick()
	if len(ts.gamePlay.entered) != 1 {
		t.Fatal("New Game did not start the game")
	}
}

func TestGameOverEnterPlacesButtons(t *testing.T) {
	ts := newTestStates()
	g := NewGameOver(ts.sm, testFonts(t))

	g.Enter(map[string]any{"victory": true, "score": 7, "hp": 3})

	btnX, btnY := (core.WindowW-200.0)/2, core.WindowH*0.65
	play, menu := g.playAgainBtn, g.menuBtn
	if play.Text != "Play Again" || play.X != btnX || play.Y != btnY || play.W != 200 || play.H != 50 {
		t.Fatalf("unexpected play again button %+v", play)
	}
	if menu.Text != "Main Menu" || menu.X != btnX || menu.Y != btnY+66 || menu.W != 200 || menu.H != 50 {
		t.Fatalf("unexpected menu button %+v", menu)
	}
	play.OnClick()
	menu.OnClick()
	if len(ts.gamePlay.entered) != 1 || len(ts.menu.entered) != 1 {
		t.Fatal("buttons did not switch states")
	}
}
