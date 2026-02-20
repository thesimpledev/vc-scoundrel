-- states/game_over.lua
-- Win/loss screen showing the final score with a "Play Again" button.

local Button       = require("ui.button")
local C            = require("core.constants")
local StateManager = require("core.state_manager")

local GameOver = {}

local title_font = nil
local score_font = nil
local default_font = nil

function GameOver:enter(data)
    if not title_font then
        title_font = love.graphics.newFont(48)
        score_font = love.graphics.newFont(24)
        default_font = love.graphics.newFont(16)
    end

    self.victory = data and data.victory or false
    self.score = data and data.score or 0
    self.hp = data and data.hp or 0

    -- Play Again button
    local btn_w, btn_h = 200, 50
    local btn_x = (C.WINDOW_W - btn_w) / 2
    local btn_y = C.WINDOW_H * 0.65

    self.play_again_btn = Button:new("Play Again", btn_x, btn_y, btn_w, btn_h, function()
        StateManager.switch("game_play")
    end)

    -- Main Menu button
    local menu_y = btn_y + btn_h + 16
    self.menu_btn = Button:new("Main Menu", btn_x, menu_y, btn_w, btn_h, function()
        StateManager.switch("menu")
    end)
end

function GameOver:update(dt)
end

function GameOver:draw()
    -- Background
    love.graphics.setColor(C.COLOR_BG)
    love.graphics.rectangle("fill", 0, 0, C.WINDOW_W, C.WINDOW_H)

    -- Title
    love.graphics.setFont(title_font)
    local title
    if self.victory then
        love.graphics.setColor(C.COLOR_GOLD)
        title = "VICTORY!"
    else
        love.graphics.setColor(C.COLOR_RED)
        title = "DEFEATED"
    end
    local title_w = title_font:getWidth(title)
    love.graphics.print(title, (C.WINDOW_W - title_w) / 2, C.WINDOW_H * 0.25)

    -- Score
    love.graphics.setFont(score_font)
    love.graphics.setColor(C.COLOR_WHITE)
    local score_text = "Score: " .. self.score
    local score_w = score_font:getWidth(score_text)
    love.graphics.print(score_text, (C.WINDOW_W - score_w) / 2, C.WINDOW_H * 0.38)

    -- HP
    local hp_text = "Final HP: " .. self.hp
    local hp_w = score_font:getWidth(hp_text)
    love.graphics.print(hp_text, (C.WINDOW_W - hp_w) / 2, C.WINDOW_H * 0.38 + 36)

    -- Buttons (use default font)
    love.graphics.setFont(default_font)
    self.play_again_btn:draw()
    self.menu_btn:draw()
end

function GameOver:mousepressed(x, y, button)
    self.play_again_btn:mousepressed(x, y, button)
    self.menu_btn:mousepressed(x, y, button)
end

function GameOver:mousereleased(x, y, button)
    self.play_again_btn:mousereleased(x, y, button)
    self.menu_btn:mousereleased(x, y, button)
end

function GameOver:mousemoved(x, y, dx, dy)
    self.play_again_btn:mousemoved(x, y)
    self.menu_btn:mousemoved(x, y)
end

function GameOver:keypressed(key)
    if key == "return" or key == "space" then
        StateManager.switch("game_play")
    end
end

return GameOver
