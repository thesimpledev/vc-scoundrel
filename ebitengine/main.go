package main

import (
	"bytes"
	"embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/states"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed assets/*
var assetFS embed.FS

// Game implements ebiten.Game and delegates to the current State.
type Game struct {
	sm *core.StateManager
}

func (g *Game) Update() error {
	return g.sm.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sm.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return core.WindowW, core.WindowH
}

func main() {
	core.InitAssets(assetFS)

	// Create font source from Go's regular font
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}

	fonts := &states.Fonts{
		Title:    &text.GoTextFace{Source: fontSource, Size: 48},
		Subtitle: &text.GoTextFace{Source: fontSource, Size: 20},
		Score:    &text.GoTextFace{Source: fontSource, Size: 24},
		Default:  &text.GoTextFace{Source: fontSource, Size: 16},
	}

	sm := core.NewStateManager()
	sm.Register("menu", states.NewMenu(sm, fonts))
	sm.Register("game_play", states.NewGamePlay(sm, fonts))
	sm.Register("game_over", states.NewGameOver(sm, fonts))

	sm.Switch("menu", nil)

	ebiten.SetWindowSize(core.WindowW, core.WindowH)
	ebiten.SetWindowTitle("Scoundrel")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&Game{sm: sm}); err != nil {
		log.Fatal(err)
	}
}
