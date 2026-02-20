-- game/room.lua
-- Room state: tracks the 4 card slots and per-turn rules.
--
-- Each turn, 4 cards are dealt into the room. The player must take 3 of them
-- (one at a time), leaving the 4th to carry over into the next room.
--
-- Slots can be sparse (e.g., {nil, Card, nil, Card}) after cards are resolved,
-- so all iteration uses `for i = 1, 4` instead of `ipairs`.
--
-- Constraint: only one health potion may be used per room.

local Class = require("lib.class")

local Room = Class:extend()

function Room:init()
    self.slots = { nil, nil, nil, nil } -- exactly 4 slots, nil = empty
    self.resolved_count = 0             -- how many cards the player has taken this room
    self.potion_used = false            -- has a potion been used this turn?
    self.initial_count = 0              -- how many cards were dealt (for completion check)
end

-- Set up the room with cards at specific slot indices.
-- `slot_cards` is a table keyed by slot index (1-4), e.g., { [2] = card, [3] = card }.
function Room:set_cards(slot_cards)
    self.slots = { nil, nil, nil, nil }
    local count = 0
    for i = 1, 4 do
        if slot_cards[i] then
            self.slots[i] = slot_cards[i]
            count = count + 1
        end
    end
    self.resolved_count = 0
    self.potion_used = false
    self.initial_count = count
end

-- Mark a card as resolved (taken by the player). Returns the card.
function Room:resolve_card(index)
    local card = self.slots[index]
    self.slots[index] = nil
    self.resolved_count = self.resolved_count + 1
    return card
end

-- Is this room complete?
-- Standard room (4 cards): take 3, leave 1 carry-over.
-- Short room (< 4 cards, near end of dungeon): take all but 1, or take all if only 1.
function Room:is_complete()
    if self.initial_count >= 4 then
        return self.resolved_count >= 3
    elseif self.initial_count > 1 then
        return self.resolved_count >= self.initial_count - 1
    else
        return self.resolved_count >= 1
    end
end

-- Get the carry-over card and its slot index.
-- Returns card, slot_index (or nil, nil if empty).
function Room:get_remaining_card()
    for i = 1, 4 do
        if self.slots[i] then
            return self.slots[i], i
        end
    end
    return nil, nil
end

-- Get all current cards (non-nil slots). Used when avoiding a room.
function Room:get_all_cards()
    local cards = {}
    for i = 1, 4 do
        if self.slots[i] then
            table.insert(cards, self.slots[i])
        end
    end
    return cards
end

-- How many cards are still in the room?
function Room:card_count()
    local count = 0
    for i = 1, 4 do
        if self.slots[i] then count = count + 1 end
    end
    return count
end

-- Is a specific slot occupied?
function Room:has_card_at(index)
    return self.slots[index] ~= nil
end

-- Clear the room entirely.
function Room:clear()
    self.slots = { nil, nil, nil, nil }
    self.resolved_count = 0
    self.potion_used = false
    self.initial_count = 0
end

return Room
