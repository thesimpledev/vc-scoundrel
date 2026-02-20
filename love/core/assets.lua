-- core/assets.lua
-- Centralized asset loader with lazy caching.
--
-- Instead of loading images at the point of use (which scatters LOVE2D calls
-- everywhere), all assets funnel through here. Each asset is loaded once on
-- first request and cached for subsequent calls.
--
-- Usage:
--   local Assets = require("core.assets")
--   local img = Assets.image("cards/cardHeartsA") -- omit extension
--   local snd = Assets.sound("audio/casino/card-slide-1")

local Assets = {}

-- Private caches
local image_cache = {}
local sound_cache = {}

-- Load (or return cached) image from assets/ directory.
-- `name` is the path relative to assets/ WITHOUT the .png extension.
function Assets.image(name)
    if not image_cache[name] then
        local path = "assets/" .. name .. ".png"
        image_cache[name] = love.graphics.newImage(path)
        -- Set nearest-neighbor filtering for crisp pixel art
        image_cache[name]:setFilter("nearest", "nearest")
    end
    return image_cache[name]
end

-- Load (or return cached) sound from assets/ directory.
-- `name` is the path relative to assets/ WITHOUT the .ogg extension.
-- `kind` is "static" (short SFX, default) or "stream" (music).
function Assets.sound(name, kind)
    kind = kind or "static"
    if not sound_cache[name] then
        local path = "assets/" .. name .. ".ogg"
        sound_cache[name] = love.audio.newSource(path, kind)
    end
    return sound_cache[name]
end

-- Play a sound effect (clones the source so overlapping plays work).
function Assets.play_sfx(name)
    local source = Assets.sound(name)
    local clone = source:clone()
    clone:play()
end

-- Flush all caches (useful on state transitions to free memory).
function Assets.clear()
    image_cache = {}
    sound_cache = {}
end

return Assets
