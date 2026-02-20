package core

import "github.com/hajimehoshi/ebiten/v2"

// State is the interface every game screen implements.
type State interface {
	Enter(data any)
	Exit()
	Update() error
	Draw(screen *ebiten.Image)
}

// StateManager routes Update/Draw to the current state and handles transitions.
type StateManager struct {
	registry map[string]State
	current  State
}

// NewStateManager creates an empty state manager.
func NewStateManager() *StateManager {
	return &StateManager{
		registry: make(map[string]State),
	}
}

// Register adds a state by name.
func (sm *StateManager) Register(name string, state State) {
	sm.registry[name] = state
}

// Switch transitions to a new state, calling Exit on the old and Enter on the new.
func (sm *StateManager) Switch(name string, data any) {
	s, ok := sm.registry[name]
	if !ok {
		panic("state_manager: no state registered with name: " + name)
	}

	if sm.current != nil {
		sm.current.Exit()
	}
	sm.current = s
	sm.current.Enter(data)
}

// Update forwards to the current state.
func (sm *StateManager) Update() error {
	if sm.current != nil {
		return sm.current.Update()
	}
	return nil
}

// Draw forwards to the current state.
func (sm *StateManager) Draw(screen *ebiten.Image) {
	if sm.current != nil {
		sm.current.Draw(screen)
	}
}
