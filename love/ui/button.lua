-- ui/button.lua
-- A clickable button with hover and press visual states.

local Class = require("lib.class")
local C     = require("core.constants")

local Button = Class:extend()

function Button:init(text, x, y, w, h, on_click)
    self.text = text
    self.x = x
    self.y = y
    self.w = w
    self.h = h
    self.on_click = on_click -- callback function

    self.hovered = false
    self.pressed = false
    self.visible = true
end

function Button:draw()
    if not self.visible then return end

    -- Background color changes on hover
    if self.hovered then
        love.graphics.setColor(C.COLOR_BUTTON_HOV)
    else
        love.graphics.setColor(C.COLOR_BUTTON)
    end

    love.graphics.rectangle("fill", self.x, self.y, self.w, self.h, 6, 6)

    -- Border
    love.graphics.setColor(C.COLOR_WHITE)
    love.graphics.rectangle("line", self.x, self.y, self.w, self.h, 6, 6)

    -- Text (centered in the button)
    love.graphics.setColor(C.COLOR_BUTTON_TXT)
    local font = love.graphics.getFont()
    local text_w = font:getWidth(self.text)
    local text_h = font:getHeight()
    love.graphics.print(
        self.text,
        self.x + (self.w - text_w) / 2,
        self.y + (self.h - text_h) / 2
    )
end

function Button:contains_point(mx, my)
    return mx >= self.x and mx <= self.x + self.w
       and my >= self.y and my <= self.y + self.h
end

function Button:mousemoved(mx, my)
    if not self.visible then return end
    self.hovered = self:contains_point(mx, my)
end

function Button:mousepressed(mx, my, button)
    if not self.visible then return end
    if button == 1 and self:contains_point(mx, my) then
        self.pressed = true
    end
end

function Button:mousereleased(mx, my, button)
    if not self.visible then return end
    if button == 1 and self.pressed and self:contains_point(mx, my) then
        if self.on_click then
            self.on_click()
        end
    end
    self.pressed = false
end

return Button
