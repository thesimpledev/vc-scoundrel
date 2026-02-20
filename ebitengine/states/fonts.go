package states

import "github.com/hajimehoshi/ebiten/v2/text/v2"

// Fonts holds the shared font faces used across all states.
type Fonts struct {
	Title    *text.GoTextFace
	Subtitle *text.GoTextFace
	Score    *text.GoTextFace
	Default  *text.GoTextFace
}
