-- game/combat.lua
-- Pure damage calculation. No LOVE2D calls — testable in isolation.
--
-- Two combat modes:
--   1. Barehanded: take full monster damage
--   2. With weapon: weapon absorbs some damage, but weapon degrades

local Combat = {}

-- Fight a monster barehanded. Returns the damage taken (full monster value).
function Combat.barehanded_damage(monster_value)
    return monster_value
end

-- Fight a monster with a weapon. Returns the damage taken.
-- Damage = max(0, monster_value - weapon_value).
function Combat.weapon_damage(monster_value, weapon_value)
    return math.max(0, monster_value - weapon_value)
end

-- Apply damage to current HP. Returns new HP (clamped to 0).
function Combat.apply_damage(current_hp, damage)
    return math.max(0, current_hp - damage)
end

-- Apply healing from a potion. Returns new HP (clamped to max).
function Combat.apply_healing(current_hp, potion_value, max_hp)
    return math.min(max_hp, current_hp + potion_value)
end

return Combat
