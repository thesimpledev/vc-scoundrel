package core

import "image/color"

// Window
const (
	WindowW = 1280
	WindowH = 720
)

// Card images are 140x190 from the Kenney asset pack
const (
	CardW     = 140
	CardH     = 190
	CardScale = 0.75
	CardDrawW = CardW * CardScale // 105
	CardDrawH = CardH * CardScale // 142.5
	CardPad   = 16
)

// Card back image to use
const CardBack = "cardBack_blue1"

// Player
const (
	MaxHP   = 20
	StartHP = 20
)

// Suits
const (
	SuitClubs    = "Clubs"
	SuitSpades   = "Spades"
	SuitDiamonds = "Diamonds"
	SuitHearts   = "Hearts"
)

var MonsterSuits = map[string]bool{
	SuitClubs:  true,
	SuitSpades: true,
}

const WeaponSuit = SuitDiamonds
const PotionSuit = SuitHearts

// RankValues maps rank strings to their numeric values for combat/healing.
var RankValues = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5,
	"6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
	"J": 11, "Q": 12, "K": 13, "A": 14,
}

// Ranks is the ordered list of all ranks for deck building.
var Ranks = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

// Layout positions — all coordinates are top-left of the card

// Dungeon pile (face-down draw pile)
const (
	DungeonX = 80
	DungeonY = 260
)

// Room slots (4 cards in a row)
const (
	RoomY      = 80
	RoomStartX = 340
	RoomSpacing = CardDrawW + CardPad // ~121
)

// Discard pile
const (
	DiscardX = 1100
	DiscardY = 260
)

// Equipped weapon
const (
	WeaponX = 420
	WeaponY = 300
)

// Slain monsters stack
const (
	SlainX       = 420
	SlainY       = 320
	SlainOffsetY = 25
)

// Animation timings (seconds)
const (
	FlipDuration = 0.25
	MoveDuration = 0.3
	DealStagger  = 0.1
)

// Colors (RGBA)
var (
	ColorBG        = color.NRGBA{38, 115, 64, 255}
	ColorWhite     = color.NRGBA{255, 255, 255, 255}
	ColorBlack     = color.NRGBA{0, 0, 0, 255}
	ColorRed       = color.NRGBA{230, 51, 51, 255}
	ColorGold      = color.NRGBA{255, 214, 0, 255}
	ColorHighlight = color.NRGBA{255, 255, 255, 77}
	ColorDim       = color.NRGBA{0, 0, 0, 102}
	ColorHPBar     = color.NRGBA{51, 204, 77, 255}
	ColorHPBG      = color.NRGBA{51, 51, 51, 255}
	ColorButton    = color.NRGBA{64, 140, 89, 255}
	ColorButtonHov = color.NRGBA{89, 166, 115, 255}
	ColorButtonTxt = color.NRGBA{255, 255, 255, 255}
)
