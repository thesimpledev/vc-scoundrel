-- main.lua
-- Entry point: wires LOVE2D callbacks to the state manager.
-- This file should stay thin — all logic lives in states and game modules.

local StateManager = require("core.state_manager")

function love.load()
    -- Use nearest-neighbor scaling for crisp pixel art
    love.graphics.setDefaultFilter("nearest", "nearest")

    -- Set up a clean, readable font
    local font = love.graphics.newFont(16)
    love.graphics.setFont(font)

    -- Seed the RNG (LOVE2D seeds automatically since 0.9.0, but explicit is clear)
    math.randomseed(os.time())

    -- Register all game states by name (avoids circular require issues)
    StateManager.register("menu", require("states.menu"))
    StateManager.register("game_play", require("states.game_play"))
    StateManager.register("game_over", require("states.game_over"))

    -- Start at the menu
    StateManager.switch("menu")
end

function love.update(dt)
    StateManager.update(dt)
end

function love.draw()
    StateManager.draw()
end

function love.keypressed(key)
    -- Let each state handle keys first (game_play uses escape for combat cancel)
    StateManager.keypressed(key)
end

function love.mousepressed(x, y, button)
    StateManager.mousepressed(x, y, button)
end

function love.mousereleased(x, y, button)
    StateManager.mousereleased(x, y, button)
end

function love.mousemoved(x, y, dx, dy)
    StateManager.mousemoved(x, y, dx, dy)
end
