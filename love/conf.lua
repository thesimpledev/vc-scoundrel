-- conf.lua
-- LOVE2D reads this file BEFORE main.lua to configure the window.
-- This is the only place to set window properties before the game starts.
-- See: https://love2d.org/wiki/Config_Files

function love.conf(t)
    t.window.title = "Scoundrel"
    t.window.width = 1280
    t.window.height = 720
    t.window.resizable = false
    t.window.vsync = 1

    -- Disable modules we don't need (keeps the binary small)
    t.modules.physics = false
    t.modules.video = false
    t.modules.thread = false
end
