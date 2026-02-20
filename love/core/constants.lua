-- core/constants.lua
-- All magic numbers live here. Changing a value here ripples everywhere.
-- Module pattern: `local M = {}; return M` — singleton, no Class needed.

local C = {}

-- Window
C.WINDOW_W = 1280
C.WINDOW_H = 720

-- Card images are 140x190 from the Kenney asset pack
C.CARD_W = 140
C.CARD_H = 190
C.CARD_SCALE = 0.75 -- render cards at 75% of native size
C.CARD_DRAW_W = C.CARD_W * C.CARD_SCALE -- 105
C.CARD_DRAW_H = C.CARD_H * C.CARD_SCALE -- 142.5
C.CARD_PAD = 16 -- spacing between cards

-- Card back image to use
C.CARD_BACK = "cardBack_blue1"

-- Player
C.MAX_HP = 20
C.START_HP = 20

-- Suits — Clubs/Spades are monsters, Diamonds are weapons, Hearts are potions
C.SUIT_CLUBS    = "Clubs"
C.SUIT_SPADES   = "Spades"
C.SUIT_DIAMONDS = "Diamonds"
C.SUIT_HEARTS   = "Hearts"

C.MONSTER_SUITS = { [C.SUIT_CLUBS] = true, [C.SUIT_SPADES] = true }
C.WEAPON_SUIT   = C.SUIT_DIAMONDS
C.POTION_SUIT   = C.SUIT_HEARTS

-- Rank values (used for damage / healing)
C.RANK_VALUES = {
    ["2"] = 2, ["3"] = 3, ["4"] = 4, ["5"] = 5,
    ["6"] = 6, ["7"] = 7, ["8"] = 8, ["9"] = 9, ["10"] = 10,
    ["J"] = 11, ["Q"] = 12, ["K"] = 13, ["A"] = 14,
}

-- Ordered list of all ranks (for deck building)
C.RANKS = { "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A" }

-- Layout positions — all coordinates are top-left of the card
-- Centered layout: dungeon left, room center, discard right

-- Dungeon pile (face-down draw pile) — far left
C.DUNGEON_X = 80
C.DUNGEON_Y = 260

-- Room slots (4 cards in a row) — centered horizontally
C.ROOM_Y = 80
C.ROOM_START_X = 340 -- leftmost room card
C.ROOM_SPACING = C.CARD_DRAW_W + C.CARD_PAD -- ~121 pixels between card starts

-- Discard pile — far right
C.DISCARD_X = 1100
C.DISCARD_Y = 260

-- Equipped weapon — below the room, left side
C.WEAPON_X = 420
C.WEAPON_Y = 300

-- Monsters slain by weapon — stacked below/on the weapon
C.SLAIN_X = 420
C.SLAIN_Y = 320 -- slightly offset so weapon rank peeks out
C.SLAIN_OFFSET_Y = 25 -- each slain card is offset down by this much

-- Animation
C.FLIP_DURATION = 0.25 -- seconds for a card flip
C.MOVE_DURATION = 0.3  -- seconds for a card slide
C.DEAL_STAGGER  = 0.1  -- delay between dealing each room card

-- Colors (RGBA 0-1)
C.COLOR_BG         = { 0.15, 0.45, 0.25, 1 } -- casino green
C.COLOR_WHITE      = { 1, 1, 1, 1 }
C.COLOR_BLACK      = { 0, 0, 0, 1 }
C.COLOR_RED        = { 0.9, 0.2, 0.2, 1 }
C.COLOR_GOLD       = { 1, 0.84, 0, 1 }
C.COLOR_HIGHLIGHT  = { 1, 1, 1, 0.3 } -- card hover overlay
C.COLOR_DIM        = { 0, 0, 0, 0.4 } -- disabled overlay
C.COLOR_HP_BAR     = { 0.2, 0.8, 0.3, 1 }
C.COLOR_HP_BG      = { 0.2, 0.2, 0.2, 1 }
C.COLOR_BUTTON     = { 0.25, 0.55, 0.35, 1 }
C.COLOR_BUTTON_HOV = { 0.35, 0.65, 0.45, 1 }
C.COLOR_BUTTON_TXT = { 1, 1, 1, 1 }

return C
