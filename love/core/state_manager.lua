-- core/state_manager.lua
-- Game state machine: switch between screens (menu, gameplay, game over).
--
-- Each state is a plain table with optional LOVE2D callback methods:
--   enter(data), exit(), update(dt), draw(),
--   keypressed(key), mousepressed(x, y, button), mousereleased(x, y, button),
--   mousemoved(x, y, dx, dy)
--
-- States register by name string to avoid circular require() issues.
-- The state_manager routes LOVE2D callbacks to the current state.
--
-- Usage:
--   local SM = require("core.state_manager")
--   SM.register("menu", require("states.menu"))
--   SM.switch("menu")

local StateManager = {}

-- Registered states: name → state table
local registry = {}

-- The currently active state
local current = nil

-- Register a state by name. Call this once per state at startup.
function StateManager.register(name, state)
    registry[name] = state
end

-- Switch to a new state. Calls exit() on old state, enter(data) on new state.
-- `data` is an optional table passed to the new state's enter() method.
function StateManager.switch(name, data)
    assert(registry[name], "No state registered with name: " .. tostring(name))

    if current and current.exit then
        current:exit()
    end

    current = registry[name]

    if current.enter then
        current:enter(data)
    end
end

-- Forward LOVE2D callbacks to the active state.
-- Each method checks that the state exists and implements the callback.

function StateManager.update(dt)
    if current and current.update then
        current:update(dt)
    end
end

function StateManager.draw()
    if current and current.draw then
        current:draw()
    end
end

function StateManager.keypressed(key)
    if current and current.keypressed then
        current:keypressed(key)
    end
end

function StateManager.mousepressed(x, y, button)
    if current and current.mousepressed then
        current:mousepressed(x, y, button)
    end
end

function StateManager.mousereleased(x, y, button)
    if current and current.mousereleased then
        current:mousereleased(x, y, button)
    end
end

function StateManager.mousemoved(x, y, dx, dy)
    if current and current.mousemoved then
        current:mousemoved(x, y, dx, dy)
    end
end

return StateManager
