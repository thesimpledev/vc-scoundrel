-- lib/anim.lua
-- Lightweight animation engine. Animates numeric fields on a target table
-- from their current values to goal values over a duration using easing.
--
-- Usage:
--   local Anim = require("lib.anim")
--   local a = Anim:new(card, { x = 300, y = 100 }, 0.5, "ease_out_quad")
--   a.on_complete = function() print("done!") end
--   -- In update loop:
--   a:update(dt)
--
-- The global list lets you fire-and-forget:
--   Anim.to(card, { x = 300 }, 0.5, "ease_out_quad", function() ... end)
--   Anim.update_all(dt)  -- call once per frame

local Class = require("lib.class")
local Utils = require("lib.utils")

local Anim = Class:extend()

-- Global list of active animations (fire-and-forget style)
local active_anims = {}

-- Easing functions: take t in [0,1], return eased t in [0,1].
-- These shape the "feel" of the animation (linear, snappy, bouncy, etc.)
local Easing = {}

function Easing.linear(t)
    return t
end

function Easing.ease_in_quad(t)
    return t * t
end

function Easing.ease_out_quad(t)
    return t * (2 - t)
end

function Easing.ease_in_out_quad(t)
    if t < 0.5 then
        return 2 * t * t
    else
        return -1 + (4 - 2 * t) * t
    end
end

function Easing.ease_out_back(t)
    local s = 1.70158
    t = t - 1
    return t * t * ((s + 1) * t + s) + 1
end

function Easing.ease_out_elastic(t)
    if t == 0 or t == 1 then return t end
    return math.pow(2, -10 * t) * math.sin((t - 0.075) * (2 * math.pi) / 0.3) + 1
end

-- Constructor: animate `target`'s fields toward `goals` over `duration` seconds.
-- `easing` is a string key into the Easing table (default: "ease_out_quad").
function Anim:init(target, goals, duration, easing)
    self.target = target
    self.goals = goals
    self.duration = duration
    self.easing_fn = Easing[easing or "ease_out_quad"] or Easing.ease_out_quad
    self.elapsed = 0
    self.done = false
    self.on_complete = nil -- optional callback

    -- Snapshot the starting values so we can interpolate from them
    self.start_values = {}
    for key, _ in pairs(goals) do
        self.start_values[key] = target[key]
    end
end

-- Advance the animation by `dt` seconds. Returns true when finished.
function Anim:update(dt)
    if self.done then return true end

    self.elapsed = self.elapsed + dt
    local t = Utils.clamp(self.elapsed / self.duration, 0, 1)
    local eased = self.easing_fn(t)

    -- Interpolate each field
    for key, goal in pairs(self.goals) do
        self.target[key] = Utils.lerp(self.start_values[key], goal, eased)
    end

    if t >= 1 then
        self.done = true
        if self.on_complete then
            self.on_complete()
        end
        return true
    end

    return false
end

-- === Global animation management (fire-and-forget) ===

-- Create an animation and add it to the global list. Returns the anim.
function Anim.to(target, goals, duration, easing, on_complete)
    local a = Anim:new(target, goals, duration, easing)
    a.on_complete = on_complete
    table.insert(active_anims, a)
    return a
end

-- Update all active global animations. Call this once per frame.
function Anim.update_all(dt)
    for i = #active_anims, 1, -1 do
        if active_anims[i]:update(dt) then
            table.remove(active_anims, i)
        end
    end
end

-- Cancel all active animations (e.g., on state transition).
function Anim.cancel_all()
    active_anims = {}
end

-- How many animations are currently running?
function Anim.active_count()
    return #active_anims
end

return Anim
