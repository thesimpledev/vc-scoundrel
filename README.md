# Scoundrel

A digital adaptation of the single-player rogue-like card game [Scoundrel](http://www.stfj.net/art/2011/Scoundrel.pdf), built in two game frameworks: [LOVE2D](https://love2d.org/) (Lua) and [Ebitengine](https://ebitengine.org/) (Go).

## About This Project

This project was **100% vibe coded with zero human intervention**. I gave AI pre-made assets, a rulebook, and said "build it." It worked.

I'm posting this as a companion piece to a [LinkedIn post I wrote](https://www.linkedin.com/feed/update/urn:li:activity:7430645021705412608/) about the difference between using AI as a force multiplier versus handing over both the vision and the execution. This game is an example of the latter — and an honest look at what that feels like.

## Why Two Implementations?

My main goal was to get a feel for how you would structure a project in LOVE2D. I've written games in MonoGame and Ebitengine before, so I was curious about LOVE and wanted something I could see as a mostly finished game in two languages.

## How to Play

You are a dungeon explorer with 20 HP. Each turn, 4 cards are dealt into a room. You must resolve 3 of them (in any order) and carry the 4th into the next room.

| Suit | Role | Effect |
|------|------|--------|
| Clubs & Spades | Monsters | Deal damage equal to their rank |
| Diamonds | Weapons | Reduce monster damage when equipped |
| Hearts | Potions | Restore HP (max 1 per room, 20 HP cap) |

Fight monsters barehanded (take full damage) or with a weapon (damage = monster value - weapon value). Weapons degrade after use and can only be used on weaker or equal monsters going forward.

**Win** by surviving until the deck is empty. **Lose** when your HP hits 0.

## Running

### LOVE2D (Lua)

```sh
cd love
love .
```

Requires [LOVE2D](https://love2d.org/) installed.

### Ebitengine (Go)

```sh
cd ebitengine
go run .
```

Requires [Go 1.24+](https://go.dev/).

## Credits

- **Game Design:** [Zach Gage](http://www.stfj.net/) and Kurt Bieg — creators of the original Scoundrel card game (2011)
- **Art Assets:** [Kenney](https://kenney.nl/) — card sprites and UI audio
