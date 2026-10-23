package lib

import "testing"

func TestAnimUpdateAllFinishesAnimationAndRunsCallback(t *testing.T) {
	AnimCancelAll()
	value := 0.0
	done := false
	AnimTo([]AnimGoal{{Field: &value, Goal: 10}}, 1, "linear", func() { done = true })

	AnimUpdateAll(0.5)
	if value != 5 {
		t.Fatalf("value after half the duration = %v, want 5", value)
	}
	if AnimActiveCount() != 1 {
		t.Fatalf("active count = %d, want 1", AnimActiveCount())
	}

	AnimUpdateAll(0.5)
	if value != 10 || !done {
		t.Fatalf("value = %v, done = %v, want 10 and true", value, done)
	}
	if AnimActiveCount() != 0 {
		t.Fatalf("active count = %d, want 0", AnimActiveCount())
	}
}

func TestAnimUpdateAllKeepsAnimationsAddedByCallbacks(t *testing.T) {
	AnimCancelAll()
	first, second := 0.0, 0.0
	AnimTo([]AnimGoal{{Field: &first, Goal: 1}}, 1, "linear", func() {
		AnimTo([]AnimGoal{{Field: &second, Goal: 1}}, 1, "linear", nil)
	})

	AnimUpdateAll(1)
	if AnimActiveCount() != 1 {
		t.Fatalf("active count after callback = %d, want 1", AnimActiveCount())
	}
	if second != 0 {
		t.Fatalf("new animation advanced in the frame it was added: %v", second)
	}

	AnimUpdateAll(1)
	if second != 1 || AnimActiveCount() != 0 {
		t.Fatalf("second = %v, active = %d, want 1 and 0", second, AnimActiveCount())
	}
}

func TestAnimUpdateAllKeepsUnfinishedAnimationsInOrder(t *testing.T) {
	AnimCancelAll()
	short, long := 0.0, 0.0
	AnimTo([]AnimGoal{{Field: &short, Goal: 1}}, 1, "linear", nil)
	AnimTo([]AnimGoal{{Field: &long, Goal: 1}}, 4, "linear", nil)

	AnimUpdateAll(1)
	if AnimActiveCount() != 1 {
		t.Fatalf("active count = %d, want 1", AnimActiveCount())
	}
	AnimUpdateAll(3)
	if long != 1 || AnimActiveCount() != 0 {
		t.Fatalf("long = %v, active = %d, want 1 and 0", long, AnimActiveCount())
	}
}
