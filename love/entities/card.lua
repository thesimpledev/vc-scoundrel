-- entities/card.lua
-- The Card entity: co-locates game data (suit, rank, value) with visual state
-- (position, scale, flip animation). In LOVE2D you *are* the engine, so
-- splitting data/view creates sync overhead with no benefit.
--
-- Cards can be face-up or face-down, animate position changes, and handle
-- mouse hit-testing for click interaction.

local Class  = require("lib.class")
local Anim   = require("lib.anim")
local Assets = require("core.assets")
local C      = require("core.constants")

local Card = Class:extend()

-- Build the asset filename for a card face. E.g., suit="Hearts", rank="A" → "cards/cardHeartsA"
local function face_asset_name(suit, rank)
    return "cards/card" .. suit .. rank
end

function Card:init(suit, rank)
    -- Game data
    self.suit  = suit
    self.rank  = rank
    self.value = C.RANK_VALUES[rank]

    -- Derived type for quick checks
    self.is_monster = C.MONSTER_SUITS[suit] or false
    self.is_weapon  = (suit == C.WEAPON_SUIT)
    self.is_potion  = (suit == C.POTION_SUIT)

    -- Asset names (resolved to actual images lazily via Assets.image)
    self.face_asset = face_asset_name(suit, rank)
    self.back_asset = "cards/" .. C.CARD_BACK

    -- Visual state
    self.x = 0
    self.y = 0
    self.scale_x = C.CARD_SCALE -- horizontal scale (used for flip animation)
    self.scale_y = C.CARD_SCALE
    self.face_up = false
    self.visible = true

    -- Flip animation state: when flipping, we scale_x to 0 (showing edge),
    -- swap the image, then scale_x back to full. This two-phase approach
    -- creates a convincing 2D card flip.
    self.flipping = false

    -- Hover/interaction state
    self.hovered = false
    self.clickable = false
end

-- Returns the image to draw (face or back depending on flip state).
function Card:get_image()
    if self.face_up then
        return Assets.image(self.face_asset)
    else
        return Assets.image(self.back_asset)
    end
end

-- Animate the card sliding to a new position.
function Card:move_to(x, y, duration, easing, on_complete)
    duration = duration or C.MOVE_DURATION
    Anim.to(self, { x = x, y = y }, duration, easing or "ease_out_quad", on_complete)
end

-- Animate a card flip (face-down → face-up or vice versa).
-- Two-phase animation: shrink scale_x to 0, swap face, grow scale_x back.
function Card:flip(on_complete)
    if self.flipping then return end
    self.flipping = true

    local half = C.FLIP_DURATION / 2

    -- Phase 1: shrink to zero width
    Anim.to(self, { scale_x = 0 }, half, "ease_in_quad", function()
        -- At the midpoint, swap which side is showing
        self.face_up = not self.face_up

        -- Phase 2: grow back to full width
        Anim.to(self, { scale_x = C.CARD_SCALE }, half, "ease_out_quad", function()
            self.flipping = false
            if on_complete then on_complete() end
        end)
    end)
end

-- Instantly set face-up state without animation.
function Card:set_face_up(face_up)
    self.face_up = face_up
    self.scale_x = C.CARD_SCALE
    self.flipping = false
end

-- Draw the card. Called every frame from the owning state's draw().
function Card:draw()
    if not self.visible then return end

    local img = self:get_image()

    -- Draw centered on the card's x,y (top-left) accounting for scale.
    -- The origin offset centers the scale_x flip around the card's center.
    local ox = C.CARD_W / 2 -- origin x: center of the unscaled image
    local oy = 0             -- origin y: top of the image

    love.graphics.setColor(C.COLOR_WHITE)
    love.graphics.draw(
        img,
        self.x + C.CARD_DRAW_W / 2, -- screen x (center of where card should appear)
        self.y,                       -- screen y (top)
        0,                            -- rotation
        self.scale_x,                 -- scale x (animated during flip)
        self.scale_y,                 -- scale y
        ox,                           -- origin x (center of source image)
        oy                            -- origin y (top of source image)
    )

    -- Hover highlight
    if self.hovered and self.clickable then
        love.graphics.setColor(C.COLOR_HIGHLIGHT)
        love.graphics.rectangle(
            "fill",
            self.x, self.y,
            C.CARD_DRAW_W, C.CARD_DRAW_H,
            4, 4 -- rounded corners
        )
    end
end

-- Check if a screen point (mx, my) is inside this card's bounds.
function Card:contains_point(mx, my)
    return mx >= self.x
       and mx <= self.x + C.CARD_DRAW_W
       and my >= self.y
       and my <= self.y + C.CARD_DRAW_H
end

-- Debug string for logging.
function Card:__tostring()
    return self.rank .. " of " .. self.suit
end

return Card
