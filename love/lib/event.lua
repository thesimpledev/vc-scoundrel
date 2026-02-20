-- lib/event.lua
-- Pub/sub event bus. Decouples game systems so they don't need direct references.
--
-- When a monster is slain, the combat system emits "monster_slain". The HUD,
-- sound system, and animation system each subscribe independently — no one
-- system needs to know about the others.
--
-- Usage:
--   local Event = require("lib.event")
--   Event.on("monster_slain", function(data) print(data.damage) end)
--   Event.emit("monster_slain", { damage = 5, card = monster_card })

local Event = {}

-- Private table of event_name → list of listener functions
local listeners = {}

-- Subscribe to an event. `fn(data)` is called when the event fires.
-- Returns the listener function (useful for unsubscribing later).
function Event.on(event_name, fn)
    if not listeners[event_name] then
        listeners[event_name] = {}
    end
    table.insert(listeners[event_name], fn)
    return fn
end

-- Emit an event. All subscribers are called in registration order.
-- `data` is an optional table passed to each listener.
function Event.emit(event_name, data)
    if not listeners[event_name] then return end
    for _, fn in ipairs(listeners[event_name]) do
        fn(data)
    end
end

-- Unsubscribe a specific listener from an event.
function Event.off(event_name, fn)
    if not listeners[event_name] then return end
    for i, listener in ipairs(listeners[event_name]) do
        if listener == fn then
            table.remove(listeners[event_name], i)
            return
        end
    end
end

-- Remove ALL listeners (e.g., when switching game states).
function Event.clear()
    listeners = {}
end

return Event
