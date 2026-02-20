-- game/scoring.lua
-- End-game score calculation. Pure functions, no LOVE2D calls.
--
-- Scoring rules from Scoundrel:
--   - If HP reached 0: find remaining monsters in dungeon, score = -(their total value)
--   - If survived: score = current HP
--   - If survived AND last card taken was a health potion: score = HP + potion value

local Scoring = {}

-- Calculate the final score.
-- `hp`         : player's current HP (0 if dead)
-- `deck`       : the Deck object (to count remaining monsters)
-- `last_card`  : the last card the player resolved (may be nil)
function Scoring.calculate(hp, deck, last_card)
    if hp <= 0 then
        -- Dead: negative score based on remaining monsters
        local remaining = deck:count_remaining_monster_value()
        return -remaining
    end

    -- Survived: base score is current HP
    local score = hp

    -- Bonus: if the very last card was a health potion, add its value
    if last_card and last_card.is_potion then
        score = score + last_card.value
    end

    return score
end

return Scoring
