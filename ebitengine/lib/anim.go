package lib

import "math"

// EasingFunc takes t in [0,1] and returns eased t in [0,1].
type EasingFunc func(t float64) float64

// Easing functions
var Easings = map[string]EasingFunc{
	"linear":           easeLinear,
	"ease_in_quad":     easeInQuad,
	"ease_out_quad":    easeOutQuad,
	"ease_in_out_quad": easeInOutQuad,
	"ease_out_back":    easeOutBack,
	"ease_out_elastic": easeOutElastic,
}

func easeLinear(t float64) float64 { return t }

func easeInQuad(t float64) float64 { return t * t }

func easeOutQuad(t float64) float64 { return t * (2 - t) }

func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

func easeOutBack(t float64) float64 {
	s := 1.70158
	t -= 1
	return t*t*((s+1)*t+s) + 1
}

func easeOutElastic(t float64) float64 {
	if t == 0 || t == 1 {
		return t
	}
	return math.Pow(2, -10*t)*math.Sin((t-0.075)*(2*math.Pi)/0.3) + 1
}

// AnimTarget is any value the animation system can interpolate.
// We use a pointer-to-float64 approach: each goal maps a *float64 to a target value.
type AnimGoal struct {
	Field *float64
	Start float64
	Goal  float64
}

// Anim animates numeric fields toward goal values over a duration.
type Anim struct {
	Goals      []AnimGoal
	Duration   float64
	EasingFn   EasingFunc
	Elapsed    float64
	Done       bool
	OnComplete func()
}

// NewAnim creates an animation. goals maps field pointers to target values.
func NewAnim(goals []AnimGoal, duration float64, easing string, onComplete func()) *Anim {
	fn, ok := Easings[easing]
	if !ok {
		fn = easeOutQuad
	}
	// Snapshot starting values
	for i := range goals {
		goals[i].Start = *goals[i].Field
	}
	return &Anim{
		Goals:      goals,
		Duration:   duration,
		EasingFn:   fn,
		OnComplete: onComplete,
	}
}

// Update advances the animation by dt seconds. Returns true when finished.
func (a *Anim) Update(dt float64) bool {
	if a.Done {
		return true
	}

	a.Elapsed += dt
	t := a.Elapsed / a.Duration
	if t > 1 {
		t = 1
	}
	eased := a.EasingFn(t)

	for _, g := range a.Goals {
		*g.Field = Lerp(g.Start, g.Goal, eased)
	}

	if t >= 1 {
		a.Done = true
		if a.OnComplete != nil {
			a.OnComplete()
		}
		return true
	}
	return false
}

// Global animation list (fire-and-forget style)
var activeAnims []*Anim

// AnimTo creates an animation and adds it to the global list.
func AnimTo(goals []AnimGoal, duration float64, easing string, onComplete func()) *Anim {
	a := NewAnim(goals, duration, easing, onComplete)
	activeAnims = append(activeAnims, a)
	return a
}

// AnimUpdateAll updates all active global animations. Call once per frame.
func AnimUpdateAll(dt float64) {
	// Snapshot count: callbacks may append new anims beyond this index.
	count := len(activeAnims)
	n := 0
	for i := 0; i < count; i++ {
		if !activeAnims[i].Update(dt) {
			activeAnims[n] = activeAnims[i]
			n++
		}
	}
	// Preserve animations added by callbacks (they sit at indices count+).
	added := make([]*Anim, len(activeAnims)-count)
	copy(added, activeAnims[count:])
	// Clear old references for GC
	for i := n; i < len(activeAnims); i++ {
		activeAnims[i] = nil
	}
	activeAnims = append(activeAnims[:n], added...)
}

// AnimCancelAll cancels all active animations.
func AnimCancelAll() {
	activeAnims = activeAnims[:0]
}

// AnimActiveCount returns how many animations are currently running.
func AnimActiveCount() int {
	return len(activeAnims)
}

// Lerp linearly interpolates between a and b by t.
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Clamp restricts v to [min, max].
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
