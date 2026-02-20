-- entities/weapon_slot.lua
-- Tracks the currently equipped weapon and the stack of monsters slain by it.
--
-- Key constraint: once a weapon kills a monster, it can only be used against
-- monsters with value LESS THAN OR EQUAL TO the last monster it killed.
-- This creates the rogue-like tension of weapon degradation.

local Class = require("lib.class")
local C     = require("core.constants")

local WeaponSlot = Class:extend()

function WeaponSlot:init()
    self.weapon = nil     -- the equipped Card (a diamond), or nil
    self.slain = {}       -- array of monster Cards killed by this weapon
end

-- Equip a new weapon. Returns the old weapon + slain stack for discarding.
function WeaponSlot:equip(weapon_card)
    local old_weapon = self.weapon
    local old_slain = self.slain

    self.weapon = weapon_card
    self.slain = {}

    return old_weapon, old_slain
end

-- Can this weapon be used against a monster of the given value?
-- First kill: always allowed. Subsequent kills: monster value must be <= last slain.
function WeaponSlot:can_use_against(monster_value)
    if not self.weapon then return false end
    if #self.slain == 0 then return true end

    local last_slain_value = self.slain[#self.slain].value
    return monster_value <= last_slain_value
end

-- Record a monster kill. The monster card goes on the slain stack.
function WeaponSlot:record_kill(monster_card)
    table.insert(self.slain, monster_card)
end

-- Calculate damage taken when using this weapon against a monster.
-- Damage = max(0, monster_value - weapon_value).
function WeaponSlot:calculate_damage(monster_value)
    if not self.weapon then return monster_value end
    return math.max(0, monster_value - self.weapon.value)
end

-- Is a weapon currently equipped?
function WeaponSlot:has_weapon()
    return self.weapon ~= nil
end

-- Get the weapon's value (for display). Returns 0 if no weapon.
function WeaponSlot:weapon_value()
    if self.weapon then return self.weapon.value end
    return 0
end

-- Clear the weapon slot entirely (e.g., on game reset).
function WeaponSlot:clear()
    self.weapon = nil
    self.slain = {}
end

return WeaponSlot
