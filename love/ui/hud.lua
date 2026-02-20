-- ui/hud.lua
-- Heads-up display: HP bar, deck count, weapon info, skip indicator,
-- and a card info bar at the bottom of the screen.

local C     = require("core.constants")
local Rules = require("game.rules")

local HUD = {}

-- Info bar dimensions
local INFO_BAR_H = 60
local INFO_BAR_Y = C.WINDOW_H - INFO_BAR_H

-- Draw the full HUD. Called every frame from game_play:draw().
function HUD.draw(hp, max_hp, deck_remaining, weapon_slot, can_avoid, phase, hovered_card, room)
    HUD.draw_hp_bar(hp, max_hp)
    HUD.draw_deck_count(deck_remaining)
    HUD.draw_weapon_info(weapon_slot)
    HUD.draw_skip_indicator(can_avoid)
    HUD.draw_card_info(hovered_card, weapon_slot, room, hp)
end

-- HP bar in the top-left corner
function HUD.draw_hp_bar(hp, max_hp)
    local x, y = 20, 20
    local w, h = 200, 24
    local fill = (hp / max_hp) * w

    -- Background
    love.graphics.setColor(C.COLOR_HP_BG)
    love.graphics.rectangle("fill", x, y, w, h, 4, 4)

    -- Fill (color shifts red when low)
    if hp <= 5 then
        love.graphics.setColor(C.COLOR_RED)
    else
        love.graphics.setColor(C.COLOR_HP_BAR)
    end
    love.graphics.rectangle("fill", x, y, fill, h, 4, 4)

    -- Border
    love.graphics.setColor(C.COLOR_WHITE)
    love.graphics.rectangle("line", x, y, w, h, 4, 4)

    -- Text
    love.graphics.setColor(C.COLOR_WHITE)
    love.graphics.print("HP: " .. hp .. " / " .. max_hp, x + 8, y + 4)
end

-- Deck count near the dungeon pile
function HUD.draw_deck_count(remaining)
    love.graphics.setColor(C.COLOR_WHITE)
    love.graphics.print("Dungeon: " .. remaining, C.DUNGEON_X, C.DUNGEON_Y - 22)
end

-- Weapon info below the weapon slot
function HUD.draw_weapon_info(weapon_slot)
    local x = C.WEAPON_X
    local y = C.WEAPON_Y - 22

    love.graphics.setColor(C.COLOR_WHITE)
    if weapon_slot:has_weapon() then
        local weapon = weapon_slot.weapon
        love.graphics.print("Weapon: " .. weapon.rank .. " of " .. weapon.suit, x, y)

        if #weapon_slot.slain > 0 then
            local last = weapon_slot.slain[#weapon_slot.slain]
            love.graphics.print(
                "Last slain: " .. last.rank .. " (can fight \u{2264}" .. last.value .. ")",
                x, y + C.CARD_DRAW_H + #weapon_slot.slain * C.SLAIN_OFFSET_Y + 10
            )
        end
    else
        love.graphics.print("No weapon equipped", x, y)
    end
end

-- "Skip available" / "Must face this room" indicator
function HUD.draw_skip_indicator(can_avoid)
    local x = C.WINDOW_W - 220
    local y = 20

    if can_avoid then
        love.graphics.setColor(C.COLOR_GOLD)
        love.graphics.print("[ Press S to skip room ]", x, y)
    else
        love.graphics.setColor(0.6, 0.6, 0.6, 0.6)
        love.graphics.print("Cannot skip this room", x, y)
    end
end

-- Card info bar at the bottom of the screen.
-- Shows what the hovered card does in the current game context.
function HUD.draw_card_info(card, weapon_slot, room, hp)
    -- Background bar (always visible so layout is stable)
    love.graphics.setColor(0, 0, 0, 0.6)
    love.graphics.rectangle("fill", 0, INFO_BAR_Y, C.WINDOW_W, INFO_BAR_H)
    love.graphics.setColor(0.4, 0.4, 0.4, 0.5)
    love.graphics.rectangle("line", 0, INFO_BAR_Y, C.WINDOW_W, INFO_BAR_H)

    if not card then
        love.graphics.setColor(0.5, 0.5, 0.5, 0.5)
        love.graphics.print("Hover over a card to see what it does", 20, INFO_BAR_Y + 8)
        return
    end

    local name = card.rank .. " of " .. card.suit
    local line1, line2

    if card.is_monster then
        line1, line2 = HUD.monster_info(card, weapon_slot, room)
    elseif card.is_weapon then
        line1, line2 = HUD.weapon_info_text(card, weapon_slot)
    elseif card.is_potion then
        line1, line2 = HUD.potion_info(card, room, hp)
    end

    -- Card name (highlighted by type color)
    if card.is_monster then
        love.graphics.setColor(C.COLOR_RED)
    elseif card.is_weapon then
        love.graphics.setColor(0.6, 0.7, 1.0, 1)
    elseif card.is_potion then
        love.graphics.setColor(C.COLOR_HP_BAR)
    end
    love.graphics.print(name, 20, INFO_BAR_Y + 8)

    -- Description lines
    local name_w = love.graphics.getFont():getWidth(name)
    love.graphics.setColor(C.COLOR_WHITE)
    if line1 then
        love.graphics.print(line1, 20 + name_w + 20, INFO_BAR_Y + 8)
    end
    if line2 then
        love.graphics.setColor(0.8, 0.8, 0.8, 0.8)
        love.graphics.print(line2, 20, INFO_BAR_Y + 30)
    end
end

-- Build info text for a monster card.
function HUD.monster_info(card, weapon_slot, room)
    local line1 = "Monster  |  Deals " .. card.value .. " damage"
    local parts = {}

    -- Barehanded option (always available)
    table.insert(parts, "Barehanded: -" .. card.value .. " HP")

    -- Weapon option
    if weapon_slot:has_weapon() then
        if weapon_slot:can_use_against(card.value) then
            local weap_dmg = math.max(0, card.value - weapon_slot:weapon_value())
            if weap_dmg == 0 then
                table.insert(parts, "Weapon: no damage!")
            else
                table.insert(parts, "Weapon: -" .. weap_dmg .. " HP")
            end
        else
            local last = weapon_slot.slain[#weapon_slot.slain]
            table.insert(parts, "Weapon: blocked (monster " .. card.value .. " > last slain " .. last.value .. ")")
        end
    end

    return line1, table.concat(parts, "    ")
end

-- Build info text for a weapon card.
function HUD.weapon_info_text(card, weapon_slot)
    local line1 = "Weapon  |  Absorbs up to " .. card.value .. " damage per fight"
    local line2

    if weapon_slot:has_weapon() then
        local old = weapon_slot.weapon
        if card.value > old.value then
            line2 = "Replaces current weapon (" .. old.rank .. " of " .. old.suit .. ") - upgrade!"
        elseif card.value < old.value then
            line2 = "Replaces current weapon (" .. old.rank .. " of " .. old.suit .. ") - downgrade"
        else
            line2 = "Replaces current weapon (same strength, resets slain stack)"
        end
    else
        line2 = "You have no weapon - equipping this lets you fight monsters with reduced damage"
    end

    return line1, line2
end

-- Build info text for a potion card.
function HUD.potion_info(card, room, hp)
    local line1 = "Health Potion  |  Heals " .. card.value .. " HP (max " .. C.MAX_HP .. ")"
    local line2

    if room.potion_used then
        line2 = "Already used a potion this room - this one will be wasted!"
    else
        local actual_heal = math.min(card.value, C.MAX_HP - hp)
        if actual_heal <= 0 then
            line2 = "HP is already full - healing would be wasted"
        elseif actual_heal < card.value then
            line2 = "Would heal " .. actual_heal .. " HP (capped at max " .. C.MAX_HP .. ")"
        else
            line2 = "Would heal " .. actual_heal .. " HP  (" .. hp .. " -> " .. (hp + actual_heal) .. ")"
        end
    end

    return line1, line2
end

return HUD
