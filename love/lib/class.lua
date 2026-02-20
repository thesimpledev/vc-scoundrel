-- lib/class.lua
-- Minimal metatable-based OOP system (~40 lines).
--
-- Usage:
--   local Class = require("lib.class")
--   local Enemy = Class:extend()
--   function Enemy:init(hp) self.hp = hp end
--   local goblin = Enemy:new(10)
--
-- How it works:
--   - Each class is a table that acts as the metatable for its instances.
--   - `__index` on the class means instances look up methods on the class table.
--   - `extend()` creates a new class whose `__index` chains back to the parent,
--     giving you single-inheritance for free.

local Class = {}
Class.__index = Class

-- Create a subclass. The new class inherits all methods from the parent.
-- Override methods by simply redefining them on the subclass.
function Class:extend()
    local cls = {}
    cls.__index = cls

    -- Chain the subclass to the parent so unresolved lookups walk up the tree.
    setmetatable(cls, { __index = self })

    -- Shortcut so subclasses can call Super.method(self, ...) if needed.
    cls.super = self

    return cls
end

-- Create a new instance. Calls `init(...)` if defined on the class.
function Class:new(...)
    local instance = setmetatable({}, self)
    if instance.init then
        instance:init(...)
    end
    return instance
end

return Class
