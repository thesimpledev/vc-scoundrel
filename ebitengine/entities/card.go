package entities

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/thesimpledev/radianthorizon/core"
	"github.com/thesimpledev/radianthorizon/lib"
)

// Card co-locates game data (suit, rank, value) with visual state
// (position, scale, flip animation).
type Card struct {
	// Game data
	Suit      string
	Rank      string
	Value     int
	IsMonster bool
	IsWeapon  bool
	IsPotion  bool

	// Asset names
	faceAsset string
	backAsset string

	// Visual state
	X       float64
	Y       float64
	ScaleX  float64
	ScaleY  float64
	FaceUp  bool
	Visible bool

	// Flip animation state
	Flipping bool

	// Hover/interaction state
	Hovered   bool
	Clickable bool
}

// NewCard creates a card with the given suit and rank.
func NewCard(suit, rank string) *Card {
	return &Card{
		Suit:      suit,
		Rank:      rank,
		Value:     core.RankValues[rank],
		IsMonster: core.MonsterSuits[suit],
		IsWeapon:  suit == core.WeaponSuit,
		IsPotion:  suit == core.PotionSuit,
		faceAsset: "cards/card" + suit + rank,
		backAsset: "cards/" + core.CardBack,
		ScaleX:    core.CardScale,
		ScaleY:    core.CardScale,
		Visible:   true,
	}
}

// GetImage returns the current image (face or back).
func (c *Card) GetImage() *ebiten.Image {
	if c.FaceUp {
		return core.LoadImage(c.faceAsset)
	}
	return core.LoadImage(c.backAsset)
}

// MoveTo animates the card sliding to a new position.
func (c *Card) MoveTo(x, y, duration float64, easing string, onComplete func()) {
	if duration == 0 {
		duration = core.MoveDuration
	}
	if easing == "" {
		easing = "ease_out_quad"
	}
	lib.AnimTo([]lib.AnimGoal{
		{Field: &c.X, Goal: x},
		{Field: &c.Y, Goal: y},
	}, duration, easing, onComplete)
}

// Flip animates a card flip (face-down to face-up or vice versa).
// Two-phase: shrink ScaleX to 0, swap face, grow ScaleX back.
func (c *Card) Flip(onComplete func()) {
	if c.Flipping {
		return
	}
	c.Flipping = true
	half := core.FlipDuration / 2

	// Phase 1: shrink to zero width
	lib.AnimTo([]lib.AnimGoal{
		{Field: &c.ScaleX, Goal: 0},
	}, half, "ease_in_quad", func() {
		// At midpoint, swap face
		c.FaceUp = !c.FaceUp

		// Phase 2: grow back
		lib.AnimTo([]lib.AnimGoal{
			{Field: &c.ScaleX, Goal: core.CardScale},
		}, half, "ease_out_quad", func() {
			c.Flipping = false
			if onComplete != nil {
				onComplete()
			}
		})
	})
}

// SetFaceUp instantly sets the face-up state without animation.
func (c *Card) SetFaceUp(faceUp bool) {
	c.FaceUp = faceUp
	c.ScaleX = core.CardScale
	c.Flipping = false
}

// Draw renders the card to the screen.
func (c *Card) Draw(screen *ebiten.Image) {
	if !c.Visible {
		return
	}

	img := c.GetImage()
	op := &ebiten.DrawImageOptions{}

	// Origin at center-x, top-y of the source image for flip centering
	ox := float64(core.CardW) / 2
	op.GeoM.Translate(-ox, 0)
	op.GeoM.Scale(c.ScaleX, c.ScaleY)
	op.GeoM.Translate(c.X+core.CardDrawW/2, c.Y)

	screen.DrawImage(img, op)

	// Hover highlight
	if c.Hovered && c.Clickable {
		vector.DrawFilledRect(screen,
			float32(c.X), float32(c.Y),
			float32(core.CardDrawW), float32(core.CardDrawH),
			core.ColorHighlight, true)
	}
}

// ContainsPoint checks if a screen point is inside this card's bounds.
func (c *Card) ContainsPoint(mx, my int) bool {
	fx, fy := float64(mx), float64(my)
	return fx >= c.X && fx <= c.X+core.CardDrawW &&
		fy >= c.Y && fy <= c.Y+core.CardDrawH
}

// String returns a debug string.
func (c *Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}
