-- states/menu.lua
-- Title screen with a "New Game" button.

local Button       = require("ui.button")
local C            = require("core.constants")
local StateManager = require("core.state_manager")

local Menu = {}

-- Fonts are created once and reused (creating fonts is expensive)
local title_font = nil
local subtitle_font = nil
local default_font = nil

function Menu:enter()
    if not title_font then
        title_font = love.graphics.newFont(48)
        subtitle_font = love.graphics.newFont(20)
        default_font = love.graphics.newFont(16)
    end

    -- Center the button horizontally, place it in the lower third of the screen
    local btn_w, btn_h = 200, 50
    local btn_x = (C.WINDOW_W - btn_w) / 2
    local btn_y = C.WINDOW_H * 0.6

    self.new_game_btn = Button:new("New Game", btn_x, btn_y, btn_w, btn_h, function()
        StateManager.switch("game_play")
    end)
end

function Menu:update(dt)
    -- Nothing to update on the menu
end

function Menu:draw()
    -- Background
    love.graphics.setColor(C.COLOR_BG)
    love.graphics.rectangle("fill", 0, 0, C.WINDOW_W, C.WINDOW_H)

    -- Title
    love.graphics.setColor(C.COLOR_WHITE)
    love.graphics.setFont(title_font)
    local title = "SCOUNDREL"
    local title_w = title_font:getWidth(title)
    love.graphics.print(title, (C.WINDOW_W - title_w) / 2, C.WINDOW_H * 0.25)

    -- Subtitle
    love.graphics.setColor(0.8, 0.8, 0.8, 0.7)
    love.graphics.setFont(subtitle_font)
    local subtitle = "A Scoundrel Card Game"
    local sub_w = subtitle_font:getWidth(subtitle)
    love.graphics.print(subtitle, (C.WINDOW_W - sub_w) / 2, C.WINDOW_H * 0.25 + 60)

    -- Reset to default font for button
    love.graphics.setFont(default_font)

    -- Button
    self.new_game_btn:draw()
end

function Menu:mousepressed(x, y, button)
    self.new_game_btn:mousepressed(x, y, button)
end

function Menu:mousereleased(x, y, button)
    self.new_game_btn:mousereleased(x, y, button)
end

function Menu:mousemoved(x, y, dx, dy)
    self.new_game_btn:mousemoved(x, y)
end

function Menu:keypressed(key)
    if key == "return" or key == "space" then
        StateManager.switch("game_play")
    end
end

return Menu
