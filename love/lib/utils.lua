-- lib/utils.lua
-- Pure utility functions. No LOVE2D dependencies — usable in any Lua project.

local Utils = {}

-- Fisher-Yates shuffle: O(n), unbiased, in-place.
-- Mutates the table and also returns it for chaining convenience.
function Utils.shuffle(t)
    for i = #t, 2, -1 do
        local j = math.random(1, i)
        t[i], t[j] = t[j], t[i]
    end
    return t
end

-- Clamp a value between min and max (inclusive).
function Utils.clamp(value, min, max)
    if value < min then return min end
    if value > max then return max end
    return value
end

-- Linear interpolation: returns the value `t` percent of the way from `a` to `b`.
-- t=0 returns a, t=1 returns b, t=0.5 returns the midpoint.
function Utils.lerp(a, b, t)
    return a + (b - a) * t
end

-- Shallow copy a table. Useful for snapshotting state without sharing references.
function Utils.shallow_copy(t)
    local copy = {}
    for k, v in pairs(t) do
        copy[k] = v
    end
    return copy
end

-- Remove a value from an array-style table. Returns true if found and removed.
function Utils.remove_value(t, value)
    for i, v in ipairs(t) do
        if v == value then
            table.remove(t, i)
            return true
        end
    end
    return false
end

return Utils
