-- game/rules.lua
-- Validation logic: "is this action legal?" Pure functions, no side effects.
-- The turn_manager calls these before executing any action.

local C = require("core.constants")

local Rules = {}

-- Can the player use a potion this turn? (Only one per room.)
function Rules.can_use_potion(room)
    return not room.potion_used
end

-- Can the player fight this monster with the equipped weapon?
-- Requires: weapon exists AND monster value <= last slain monster value.
function Rules.can_use_weapon(weapon_slot, monster_card)
    if not weapon_slot:has_weapon() then return false end
    return weapon_slot:can_use_against(monster_card.value)
end

-- Can the player avoid this room? (Cannot avoid two rooms in a row.)
function Rules.can_avoid_room(avoided_last_room)
    return not avoided_last_room
end

-- Is the game over? (HP <= 0 or dungeon + room are both empty.)
function Rules.is_game_over(hp, deck, room)
    if hp <= 0 then return true end
    if deck:is_empty() and room:card_count() == 0 then return true end
    return false
end

-- Did the player win? (Made it through the entire dungeon alive.)
function Rules.is_victory(hp, deck, room)
    return hp > 0 and deck:is_empty() and room:card_count() == 0
end

-- Get the type of a card as a string (for UI/logic branching).
function Rules.card_type(card)
    if card.is_monster then return "monster" end
    if card.is_weapon  then return "weapon" end
    if card.is_potion  then return "potion" end
    return "unknown"
end

return Rules
