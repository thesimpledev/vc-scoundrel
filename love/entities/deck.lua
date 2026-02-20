-- entities/deck.lua
-- The Dungeon deck: builds the 44-card Scoundrel deck, shuffles it,
-- and provides draw operations.
--
-- Scoundrel deck composition:
--   - All Clubs (2-A)    = 13 monster cards
--   - All Spades (2-A)   = 13 monster cards
--   - Diamonds 2-10      =  9 weapon cards (face cards + ace removed)
--   - Hearts 2-10        =  9 potion cards (face cards + ace removed)
--   Total: 44 cards

local Class = require("lib.class")
local Card  = require("entities.card")
local Utils = require("lib.utils")
local C     = require("core.constants")

local Deck = Class:extend()

function Deck:init()
    self.cards = {}
    self:build()
    self:shuffle()
end

-- Build the 44-card deck according to Scoundrel rules.
function Deck:build()
    self.cards = {}

    -- Clubs and Spades: all 13 ranks (monsters)
    for _, suit in ipairs({ C.SUIT_CLUBS, C.SUIT_SPADES }) do
        for _, rank in ipairs(C.RANKS) do
            table.insert(self.cards, Card:new(suit, rank))
        end
    end

    -- Diamonds and Hearts: only 2-10 (face cards and aces removed)
    local red_ranks = { "2", "3", "4", "5", "6", "7", "8", "9", "10" }
    for _, suit in ipairs({ C.SUIT_DIAMONDS, C.SUIT_HEARTS }) do
        for _, rank in ipairs(red_ranks) do
            table.insert(self.cards, Card:new(suit, rank))
        end
    end
end

function Deck:shuffle()
    Utils.shuffle(self.cards)
end

-- Draw one card from the top of the deck. Returns nil if empty.
function Deck:draw_card()
    if #self.cards == 0 then return nil end
    return table.remove(self.cards)
end

-- Place cards at the bottom of the deck (used when avoiding a room).
function Deck:place_at_bottom(cards)
    for _, card in ipairs(cards) do
        table.insert(self.cards, 1, card)
    end
end

-- How many cards remain?
function Deck:remaining()
    return #self.cards
end

-- Is the deck empty?
function Deck:is_empty()
    return #self.cards == 0
end

-- Count remaining monsters in the deck (for scoring when HP reaches 0).
function Deck:count_remaining_monster_value()
    local total = 0
    for _, card in ipairs(self.cards) do
        if card.is_monster then
            total = total + card.value
        end
    end
    return total
end

return Deck
