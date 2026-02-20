package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thesimpledev/radianthorizon/core"
)

// Button is a clickable UI button with hover and press states.
type Button struct {
	Text    string
	X, Y    float64
	W, H    float64
	OnClick func()
	Font    *text.GoTextFace

	Hovered bool
	Pressed bool
	Visible bool
}

// NewButton creates a button.
func NewButton(label string, x, y, w, h float64, font *text.GoTextFace, onClick func()) *Button {
	return &Button{
		Text:    label,
		X:       x,
		Y:       y,
		W:       w,
		H:       h,
		OnClick: onClick,
		Font:    font,
		Visible: true,
	}
}

// Draw renders the button.
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.Visible {
		return
	}

	// Background
	var bgColor color.NRGBA
	if b.Hovered {
		bgColor = core.ColorButtonHov
	} else {
		bgColor = core.ColorButton
	}
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), bgColor, true)

	// Border
	vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, core.ColorWhite, true)

	// Text centered
	op := &text.DrawOptions{}
	tw, th := text.Measure(b.Text, b.Font, 0)
	op.GeoM.Translate(b.X+(b.W-tw)/2, b.Y+(b.H-th)/2)
	op.ColorScale.ScaleWithColor(core.ColorButtonTxt)
	text.Draw(screen, b.Text, b.Font, op)
}

// ContainsPoint checks if a point is inside the button.
func (b *Button) ContainsPoint(mx, my int) bool {
	fx, fy := float64(mx), float64(my)
	return fx >= b.X && fx <= b.X+b.W &&
		fy >= b.Y && fy <= b.Y+b.H
}

// HandleMouseMove updates hover state.
func (b *Button) HandleMouseMove(mx, my int) {
	if !b.Visible {
		return
	}
	b.Hovered = b.ContainsPoint(mx, my)
}

// HandleMousePress records a press.
func (b *Button) HandleMousePress(mx, my int) {
	if !b.Visible {
		return
	}
	if b.ContainsPoint(mx, my) {
		b.Pressed = true
	}
}

// HandleMouseRelease triggers the click if appropriate.
func (b *Button) HandleMouseRelease(mx, my int) {
	if !b.Visible {
		return
	}
	if b.Pressed && b.ContainsPoint(mx, my) {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
	b.Pressed = false
}
