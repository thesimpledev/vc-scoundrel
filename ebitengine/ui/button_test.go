package ui

import "testing"

func TestNewButtonUsesBoundsAndIsVisible(t *testing.T) {
	b := NewButton("Go", Rect{X: 10, Y: 20, W: 100, H: 40}, nil, nil)

	if b.Text != "Go" || b.X != 10 || b.Y != 20 || b.W != 100 || b.H != 40 {
		t.Fatalf("unexpected button %+v", b)
	}
	if !b.Visible {
		t.Fatal("new button is not visible")
	}
}

func TestButtonClickRequiresPressAndReleaseInside(t *testing.T) {
	clicks := 0
	b := NewButton("Go", Rect{X: 10, Y: 20, W: 100, H: 40}, nil, func() { clicks++ })

	b.HandleMouseRelease(50, 30)
	if clicks != 0 {
		t.Fatal("release without a press clicked")
	}

	b.HandleMousePress(50, 30)
	b.HandleMouseRelease(500, 300)
	if clicks != 0 {
		t.Fatal("release outside the button clicked")
	}

	b.HandleMousePress(50, 30)
	b.HandleMouseRelease(60, 40)
	if clicks != 1 {
		t.Fatalf("clicks = %d, want 1", clicks)
	}
}

func TestButtonHover(t *testing.T) {
	b := NewButton("Go", Rect{X: 10, Y: 20, W: 100, H: 40}, nil, nil)
	b.HandleMouseMove(10, 20)
	if !b.Hovered {
		t.Fatal("not hovered at the top-left corner")
	}
	b.HandleMouseMove(9, 20)
	if b.Hovered {
		t.Fatal("hovered outside the button")
	}
}

func TestHiddenButtonIgnoresInput(t *testing.T) {
	clicks := 0
	b := NewButton("Go", Rect{X: 0, Y: 0, W: 10, H: 10}, nil, func() { clicks++ })
	b.Visible = false

	b.HandleMouseMove(5, 5)
	b.HandleMousePress(5, 5)
	b.HandleMouseRelease(5, 5)

	if b.Hovered || clicks != 0 {
		t.Fatal("hidden button reacted to input")
	}
}
