-- game/turn_manager.lua
-- Orchestrates the turn cycle: deal → resolve 3 → carry over → repeat.
-- This module manages the flow of a Scoundrel game but delegates rendering
-- and input to the game_play state.

local Class      = require("lib.class")
local Deck       = require("entities.deck")
local WeaponSlot = require("entities.weapon_slot")
local Room       = require("game.room")
local Combat     = require("game.combat")
local Rules      = require("game.rules")
local Scoring    = require("game.scoring")
local Event      = require("lib.event")
local C          = require("core.constants")

local TurnManager = Class:extend()

function TurnManager:init()
    self.deck = Deck:new()
    self.room = Room:new()
    self.weapon_slot = WeaponSlot:new()

    self.hp = C.START_HP
    self.avoided_last_room = false
    self.last_card = nil -- last card resolved (for scoring)

    -- The carry-over card from the previous room (4th unresolved card)
    self.carry_over = nil
    self.carry_over_slot = nil -- which slot (1-4) it was in

    -- Game flow state
    self.game_over = false
    self.victory = false
    self.score = 0

    -- Track the phase of the current turn
    -- "dealing" → "resolving" → "room_complete" → back to "dealing"
    self.phase = "ready"
end

-- Start a new turn: deal cards into the room.
-- Returns slot_cards (table keyed 1-4) and carry_over_slot (or nil).
-- The carry-over stays in its original slot; new cards fill the empty ones.
function TurnManager:deal_room()
    self.phase = "dealing"
    local slot_cards = {}
    local co_slot = nil

    -- Carry-over stays in its original slot
    if self.carry_over then
        co_slot = self.carry_over_slot
        slot_cards[co_slot] = self.carry_over
        self.carry_over = nil
        self.carry_over_slot = nil
    end

    -- Fill empty slots from the dungeon deck
    for i = 1, 4 do
        if not slot_cards[i] and not self.deck:is_empty() then
            slot_cards[i] = self.deck:draw_card()
        end
    end

    self.room:set_cards(slot_cards)
    self.phase = "resolving"

    Event.emit("room_dealt", { slot_cards = slot_cards, carry_over_slot = co_slot })
    return slot_cards, co_slot
end

-- Can the player avoid the current room?
function TurnManager:can_avoid()
    return Rules.can_avoid_room(self.avoided_last_room)
        and self.room:card_count() == 4  -- can only avoid a full 4-card room
        and self.phase == "resolving"
        and self.room.resolved_count == 0 -- haven't started resolving yet
end

-- Avoid the current room: all 4 cards go to bottom of dungeon.
function TurnManager:avoid_room()
    if not self:can_avoid() then return false end

    local cards = self.room:get_all_cards()
    self.room:clear()
    self.deck:place_at_bottom(cards)
    self.avoided_last_room = true

    Event.emit("room_avoided", { cards = cards })

    -- Don't auto-deal here — the visual layer (game_play) handles dealing
    -- so it can animate the cards moving back to the dungeon first.
    return true
end

-- Resolve a card at the given room slot index.
-- For monsters, `use_weapon` determines barehanded vs weapon combat.
-- Returns a result table describing what happened.
function TurnManager:resolve_card(slot_index, use_weapon)
    if self.phase ~= "resolving" then return nil end
    if not self.room:has_card_at(slot_index) then return nil end

    local card = self.room.slots[slot_index]
    local result = { card = card, slot = slot_index }

    if card.is_potion then
        result = self:resolve_potion(slot_index, card)
    elseif card.is_weapon then
        result = self:resolve_weapon(slot_index, card)
    elseif card.is_monster then
        result = self:resolve_monster(slot_index, card, use_weapon)
    end

    -- Check if room is complete after resolving
    if self.room:is_complete() then
        self:finish_room()
    end

    -- Check for game over
    if self.hp <= 0 then
        self:end_game()
    end

    return result
end

-- Resolve a potion card.
function TurnManager:resolve_potion(slot_index, card)
    local healed = 0
    local wasted = false

    if Rules.can_use_potion(self.room) then
        local old_hp = self.hp
        self.hp = Combat.apply_healing(self.hp, card.value, C.MAX_HP)
        healed = self.hp - old_hp
        self.room.potion_used = true
    else
        -- Second potion this turn: wasted, no healing
        wasted = true
    end

    self.room:resolve_card(slot_index)
    self.last_card = card

    local result = {
        action = "potion",
        card = card,
        slot = slot_index,
        healed = healed,
        wasted = wasted,
    }

    Event.emit("potion_used", result)
    return result
end

-- Resolve a weapon card (equip it).
function TurnManager:resolve_weapon(slot_index, card)
    local old_weapon, old_slain = self.weapon_slot:equip(card)

    self.room:resolve_card(slot_index)
    self.last_card = card

    local result = {
        action = "weapon_equipped",
        card = card,
        slot = slot_index,
        old_weapon = old_weapon,
        old_slain = old_slain,
    }

    Event.emit("weapon_equipped", result)
    return result
end

-- Resolve a monster card (fight it).
function TurnManager:resolve_monster(slot_index, card, use_weapon)
    local damage
    local used_weapon = false

    if use_weapon and Rules.can_use_weapon(self.weapon_slot, card) then
        -- Fight with weapon
        damage = Combat.weapon_damage(card.value, self.weapon_slot:weapon_value())
        self.weapon_slot:record_kill(card)
        used_weapon = true
    else
        -- Fight barehanded
        damage = Combat.barehanded_damage(card.value)
    end

    self.hp = Combat.apply_damage(self.hp, damage)
    self.room:resolve_card(slot_index)
    self.last_card = card

    local result = {
        action = "monster_fought",
        card = card,
        slot = slot_index,
        damage = damage,
        used_weapon = used_weapon,
    }

    Event.emit("monster_fought", result)
    return result
end

-- Called when all 3 cards have been resolved in a room.
function TurnManager:finish_room()
    -- The remaining card carries over to the next room, keeping its slot position
    self.carry_over, self.carry_over_slot = self.room:get_remaining_card()
    self.room:clear()
    self.avoided_last_room = false -- can avoid the next room

    self.phase = "room_complete"
    Event.emit("room_complete", {
        carry_over = self.carry_over,
        carry_over_slot = self.carry_over_slot,
    })

    -- If the deck is empty and there's no carry-over, game is over
    if self.deck:is_empty() and not self.carry_over then
        self:end_game()
    end
end

-- End the game and calculate the score.
function TurnManager:end_game()
    self.game_over = true
    self.victory = self.hp > 0
    self.score = Scoring.calculate(self.hp, self.deck, self.last_card)

    Event.emit("game_over", {
        victory = self.victory,
        score = self.score,
        hp = self.hp,
    })
end

-- Get available actions for a card at the given slot index.
-- Returns a table of action strings the player can take.
function TurnManager:get_actions(slot_index)
    if self.phase ~= "resolving" then return {} end
    if not self.room:has_card_at(slot_index) then return {} end

    local card = self.room.slots[slot_index]
    local actions = {}

    if card.is_potion then
        -- Potions can always be taken (even if wasted)
        table.insert(actions, "use_potion")
    elseif card.is_weapon then
        table.insert(actions, "equip_weapon")
    elseif card.is_monster then
        table.insert(actions, "fight_barehanded")
        if Rules.can_use_weapon(self.weapon_slot, card) then
            table.insert(actions, "fight_with_weapon")
        end
    end

    return actions
end

return TurnManager
