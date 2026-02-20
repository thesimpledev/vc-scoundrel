package lib

// EventBus is a pub/sub event system. It uses a global instance for simplicity,
// matching the Lua version.

type listenerEntry struct {
	fn func(any)
	id int
}

var (
	listeners = map[string][]listenerEntry{}
	nextID    int
)

// EventOn subscribes to an event. Returns a listener ID for unsubscribing.
func EventOn(eventName string, fn func(any)) int {
	nextID++
	id := nextID
	listeners[eventName] = append(listeners[eventName], listenerEntry{fn: fn, id: id})
	return id
}

// EventEmit broadcasts an event to all subscribers.
func EventEmit(eventName string, data any) {
	for _, entry := range listeners[eventName] {
		entry.fn(data)
	}
}

// EventOff removes a specific listener by its ID.
func EventOff(eventName string, id int) {
	entries := listeners[eventName]
	for i, entry := range entries {
		if entry.id == id {
			listeners[eventName] = append(entries[:i], entries[i+1:]...)
			return
		}
	}
}

// EventClear removes all listeners.
func EventClear() {
	listeners = map[string][]listenerEntry{}
}
