-- states/game_play.lua
-- The main gameplay state. Owns all game objects (TurnManager, visual cards),
-- handles input, orchestrates animations, and draws everything.
--
-- Flow:
-- 1. Deal 4 cards into the room (with animation)
-- 2. Player clicks cards to resolve them (take 3 of 4)
-- 3. For monsters: show combat choice buttons (barehanded vs weapon)
-- 4. After 3 resolved, carry over the 4th and deal next room
-- 5. Repeat until game over

local TurnManager = require("game.turn_manager")
local Anim        = require("lib.anim")
local Event       = require("lib.event")
local Assets      = require("core.assets")
local Button      = require("ui.button")
local HUD         = require("ui.hud")
local C           = require("core.constants")
local StateManager = require("core.state_manager")

local GamePlay = {}

function GamePlay:enter()
    -- Clear any leftover state
    Anim.cancel_all()
    Event.clear()

    -- Core game logic
    self.turn = TurnManager:new()

    -- Visual card tracking: all cards currently visible on screen
    -- (room cards, weapon card, slain stack, dungeon top, discard pile)
    self.room_cards = { nil, nil, nil, nil } -- indexed by slot (1-4)
    self.weapon_card = nil                    -- visual weapon card
    self.slain_cards = {}                     -- visual slain stack
    self.discard_top = nil                    -- visual top of discard
    self.flying_cards = {}                    -- cards mid-animation (removed from room, not yet landed)

    -- Combat choice UI (shown when clicking a monster)
    self.combat_buttons = {}
    self.choosing_combat = false
    self.combat_slot = nil

    -- Animation lock: when true, input is ignored (cards are moving)
    self.animating = false

    -- Card info tooltip: whichever card the mouse is over
    self.hovered_card = nil

    -- Game over transition
    self.game_over_pending = false
    self.game_over_timer = 0
    self.game_over_data = nil

    -- Subscribe to events
    self:setup_events()

    -- Deal the first room
    self:do_deal()
end

function GamePlay:exit()
    Anim.cancel_all()
    Event.clear()
end

-- === Event Subscriptions ===

function GamePlay:setup_events()
    Event.on("monster_fought", function(data)
        if data.damage > 0 then
            Assets.play_sfx("audio/casino/card-shove-1")
        else
            Assets.play_sfx("audio/casino/card-place-2")
        end
    end)

    Event.on("potion_used", function(data)
        if data.wasted then
            Assets.play_sfx("audio/casino/card-shove-2")
        else
            Assets.play_sfx("audio/casino/card-place-1")
        end
    end)

    Event.on("weapon_equipped", function(data)
        Assets.play_sfx("audio/casino/card-place-3")
    end)

    Event.on("room_avoided", function(data)
        Assets.play_sfx("audio/casino/card-shuffle")
    end)

    Event.on("game_over", function(data)
        -- Short delay before switching to game over screen
        self.game_over_pending = true
        self.game_over_timer = 1.0
        self.game_over_data = data
    end)
end

-- === Deal Animation ===

function GamePlay:do_deal()
    self.animating = true
    local slot_cards, carry_slot = self.turn:deal_room()

    -- Count how many new cards need to animate in (everything except carry-over)
    local new_card_slots = {}
    for i = 1, 4 do
        if slot_cards[i] and i ~= carry_slot then
            table.insert(new_card_slots, i)
        end
    end

    local pending_anims = #new_card_slots

    -- If there are no new cards to animate (edge case), just unlock immediately
    if pending_anims == 0 then
        for i = 1, 4 do
            self.room_cards[i] = slot_cards[i]
        end
        self.animating = false
        self:enable_room_interaction()
        return
    end

    -- Place all cards into visual slots
    for i = 1, 4 do
        self.room_cards[i] = slot_cards[i]
    end

    -- Carry-over card stays exactly where it is — no animation needed.
    -- New cards animate from the dungeon pile with a stagger.
    local stagger_i = 0
    for _, slot in ipairs(new_card_slots) do
        local card = slot_cards[slot]
        card.x = C.DUNGEON_X
        card.y = C.DUNGEON_Y
        card:set_face_up(false)
        card.visible = true

        local target_x = C.ROOM_START_X + (slot - 1) * C.ROOM_SPACING
        local delay = stagger_i * C.DEAL_STAGGER
        stagger_i = stagger_i + 1

        self:delayed_call(delay, function()
            Assets.play_sfx("audio/casino/card-slide-" .. math.random(1, 8))
            card:move_to(target_x, C.ROOM_Y, C.MOVE_DURATION, "ease_out_quad", function()
                card:flip(function()
                    pending_anims = pending_anims - 1
                    if pending_anims <= 0 then
                        self.animating = false
                        self:enable_room_interaction()
                    end
                end)
            end)
        end)
    end
end

-- Simple delayed callback using an anim on a dummy target.
function GamePlay:delayed_call(delay, fn)
    if delay <= 0 then
        fn()
        return
    end
    local dummy = { t = 0 }
    Anim.to(dummy, { t = 1 }, delay, "linear", fn)
end

-- === Card Interaction ===

-- Make room cards clickable after deal animation finishes.
function GamePlay:enable_room_interaction()
    for i = 1, 4 do
        local card = self.room_cards[i]
        if card then
            card.clickable = true
        end
    end
end

-- Disable all card interaction (during animations or combat choice).
function GamePlay:disable_room_interaction()
    for i = 1, 4 do
        local card = self.room_cards[i]
        if card then
            card.clickable = false
            card.hovered = false
        end
    end
end

-- Handle clicking on a room card.
function GamePlay:on_card_clicked(slot_index)
    if self.animating or self.choosing_combat then return end

    local card = self.room_cards[slot_index]
    if not card then return end

    local actions = self.turn:get_actions(slot_index)
    if #actions == 0 then return end

    -- If it's a monster with weapon option, show combat choice UI
    if card.is_monster and #actions > 1 then
        self:show_combat_choice(slot_index, actions)
        return
    end

    -- Otherwise resolve immediately
    local use_weapon = (actions[1] == "fight_with_weapon")
    self:do_resolve(slot_index, use_weapon)
end

-- === Combat Choice UI ===

function GamePlay:show_combat_choice(slot_index, actions)
    self.choosing_combat = true
    self.combat_slot = slot_index
    self:disable_room_interaction()

    local card = self.room_cards[slot_index]
    local btn_w, btn_h = 180, 40
    local btn_x = card.x + C.CARD_DRAW_W + 12
    local btn_y = card.y

    self.combat_buttons = {}

    -- Barehanded button
    local bare_dmg = card.value
    table.insert(self.combat_buttons, Button:new(
        "Barehanded (-" .. bare_dmg .. " HP)",
        btn_x, btn_y, btn_w, btn_h,
        function() self:choose_combat(false) end
    ))

    -- Weapon button (if available)
    for _, action in ipairs(actions) do
        if action == "fight_with_weapon" then
            local weap_dmg = math.max(0, card.value - self.turn.weapon_slot:weapon_value())
            table.insert(self.combat_buttons, Button:new(
                "Use Weapon (-" .. weap_dmg .. " HP)",
                btn_x, btn_y + btn_h + 8, btn_w, btn_h,
                function() self:choose_combat(true) end
            ))
        end
    end
end

function GamePlay:choose_combat(use_weapon)
    self.choosing_combat = false
    self.combat_buttons = {}
    self:do_resolve(self.combat_slot, use_weapon)
    self.combat_slot = nil
end

function GamePlay:hide_combat_choice()
    self.choosing_combat = false
    self.combat_buttons = {}
    self.combat_slot = nil
    self:enable_room_interaction()
end

-- === Resolve a Card ===

function GamePlay:do_resolve(slot_index, use_weapon)
    self.animating = true
    self:disable_room_interaction()

    local card = self.room_cards[slot_index]
    local result = self.turn:resolve_card(slot_index, use_weapon)
    if not result then
        self.animating = false
        self:enable_room_interaction()
        return
    end

    -- Animate the card to its destination based on what happened
    card.clickable = false
    card.hovered = false

    if result.action == "potion" then
        self:animate_to_discard(card, slot_index, function()
            self:after_resolve()
        end)
    elseif result.action == "weapon_equipped" then
        self:animate_equip_weapon(card, slot_index, result, function()
            self:after_resolve()
        end)
    elseif result.action == "monster_fought" then
        if result.used_weapon then
            self:animate_to_slain(card, slot_index, function()
                self:after_resolve()
            end)
        else
            self:animate_to_discard(card, slot_index, function()
                self:after_resolve()
            end)
        end
    end
end

function GamePlay:after_resolve()
    if self.turn.game_over then
        self.animating = false
        return
    end

    if self.turn.room:is_complete() or self.turn.phase == "room_complete" then
        -- Room is complete — find carry-over and deal next room
        self:handle_room_complete()
    else
        self.animating = false
        self:enable_room_interaction()
    end
end

function GamePlay:handle_room_complete()
    -- The turn_manager already stored the carry-over card and its slot.
    -- The carry-over stays in room_cards at its slot — don't clear it.
    -- Just nil out the resolved (empty) slots so they're ready for new cards.
    local has_carry = self.turn.carry_over ~= nil

    if not has_carry then
        self.room_cards = { nil, nil, nil, nil }
    end
    -- (If there IS a carry-over, it's still sitting in room_cards at its slot.
    --  The other 3 slots are already nil from animate_to_discard/slain.)

    -- If the deck still has cards (or we have a carry-over), deal next room
    if not self.turn.deck:is_empty() or has_carry then
        if not self.turn.game_over then
            self:delayed_call(0.3, function()
                self:do_deal()
            end)
        else
            self.animating = false
        end
    else
        -- Deck is empty and no more cards: game ends
        self.animating = false
        if not self.turn.game_over then
            self.turn:end_game()
        end
    end
end

-- === Card Animations ===

-- Add a card to the flying list (drawn separately while animating).
function GamePlay:start_flying(card)
    table.insert(self.flying_cards, card)
end

-- Remove a card from the flying list (it has landed at its destination).
function GamePlay:stop_flying(card)
    for i, c in ipairs(self.flying_cards) do
        if c == card then
            table.remove(self.flying_cards, i)
            return
        end
    end
end

-- Animate a card sliding to the discard pile.
function GamePlay:animate_to_discard(card, slot_index, on_complete)
    self.room_cards[slot_index] = nil
    self:start_flying(card)
    card:move_to(C.DISCARD_X, C.DISCARD_Y, C.MOVE_DURATION, "ease_in_out_quad", function()
        self:stop_flying(card)
        self.discard_top = card
        card:set_face_up(false) -- discard is face-down
        if on_complete then on_complete() end
    end)
end

-- Animate equipping a weapon (card slides to weapon slot).
function GamePlay:animate_equip_weapon(card, slot_index, result, on_complete)
    self.room_cards[slot_index] = nil
    self:start_flying(card)

    -- If there was an old weapon, animate it and slain stack to discard first
    if result.old_weapon then
        local old_cards = { result.old_weapon }
        for _, slain in ipairs(result.old_slain or {}) do
            table.insert(old_cards, slain)
        end

        local discard_pending = #old_cards
        for _, old_card in ipairs(old_cards) do
            self:start_flying(old_card)
            old_card:move_to(C.DISCARD_X, C.DISCARD_Y, C.MOVE_DURATION, "ease_in_out_quad", function()
                self:stop_flying(old_card)
                old_card:set_face_up(false)
                discard_pending = discard_pending - 1
                if discard_pending <= 0 then
                    self.discard_top = old_card
                end
            end)
        end
        self.weapon_card = nil
        self.slain_cards = {}
    end

    -- Move new weapon to weapon slot
    card:move_to(C.WEAPON_X, C.WEAPON_Y, C.MOVE_DURATION, "ease_out_back", function()
        self:stop_flying(card)
        self.weapon_card = card
        if on_complete then on_complete() end
    end)
end

-- Animate placing a monster on the slain stack.
function GamePlay:animate_to_slain(card, slot_index, on_complete)
    self.room_cards[slot_index] = nil
    local slain_index = #self.slain_cards

    -- Add to slain_cards immediately so it draws with the stack
    table.insert(self.slain_cards, card)

    local target_y = C.SLAIN_Y + slain_index * C.SLAIN_OFFSET_Y
    card:move_to(C.SLAIN_X, target_y, C.MOVE_DURATION, "ease_out_quad", function()
        if on_complete then on_complete() end
    end)
end

-- === Update ===

function GamePlay:update(dt)
    Anim.update_all(dt)

    -- Game over delay timer
    if self.game_over_pending then
        self.game_over_timer = self.game_over_timer - dt
        if self.game_over_timer <= 0 then
            self.game_over_pending = false
            StateManager.switch("game_over", self.game_over_data)
        end
    end
end

-- === Draw ===

function GamePlay:draw()
    -- Background
    love.graphics.setColor(C.COLOR_BG)
    love.graphics.rectangle("fill", 0, 0, C.WINDOW_W, C.WINDOW_H)

    -- Dungeon pile (draw card back if cards remain)
    if self.turn.deck:remaining() > 0 then
        love.graphics.setColor(C.COLOR_WHITE)
        local back_img = Assets.image("cards/" .. C.CARD_BACK)
        love.graphics.draw(back_img, C.DUNGEON_X, C.DUNGEON_Y, 0, C.CARD_SCALE, C.CARD_SCALE)
    end

    -- Discard pile
    if self.discard_top then
        self.discard_top:draw()
    end

    -- Weapon card
    if self.weapon_card then
        self.weapon_card:draw()
    end

    -- Slain monster stack (draw bottom-up so top card is visually on top)
    for _, card in ipairs(self.slain_cards) do
        card:draw()
    end

    -- Room cards
    for i = 1, 4 do
        local card = self.room_cards[i]
        if card then
            card:draw()
        end
    end

    -- Flying cards (mid-animation, drawn on top of everything)
    for _, card in ipairs(self.flying_cards) do
        card:draw()
    end

    -- Combat choice buttons
    for _, btn in ipairs(self.combat_buttons) do
        btn:draw()
    end

    -- HUD
    HUD.draw(
        self.turn.hp,
        C.MAX_HP,
        self.turn.deck:remaining(),
        self.turn.weapon_slot,
        self.turn:can_avoid(),
        self.turn.phase,
        self.hovered_card,
        self.turn.room
    )
end

-- === Input ===

function GamePlay:mousepressed(x, y, button)
    if button ~= 1 then return end

    -- Combat choice buttons get priority
    if self.choosing_combat then
        for _, btn in ipairs(self.combat_buttons) do
            btn:mousepressed(x, y, button)
        end
        return
    end

    if self.animating then return end

    -- Check room cards for clicks
    for i = 1, 4 do
        local card = self.room_cards[i]
        if card and card.clickable and card:contains_point(x, y) then
            self:on_card_clicked(i)
            return
        end
    end
end

function GamePlay:mousereleased(x, y, button)
    if button ~= 1 then return end

    if self.choosing_combat then
        for _, btn in ipairs(self.combat_buttons) do
            btn:mousereleased(x, y, button)
        end
    end
end

function GamePlay:mousemoved(x, y, dx, dy)
    -- Update hover state on combat buttons
    if self.choosing_combat then
        for _, btn in ipairs(self.combat_buttons) do
            btn:mousemoved(x, y)
        end
        return
    end

    if self.animating then
        self.hovered_card = nil
        return
    end

    -- Update hover state on room cards and track which card is hovered
    self.hovered_card = nil
    for i = 1, 4 do
        local card = self.room_cards[i]
        if card and card.clickable then
            local over = card:contains_point(x, y)
            card.hovered = over
            if over then
                self.hovered_card = card
            end
        end
    end
end

function GamePlay:keypressed(key)
    if self.animating then return end

    -- S to skip/avoid the current room
    if key == "s" and self.turn:can_avoid() then
        self:do_avoid_room()
    end

    -- Escape to cancel combat choice
    if key == "escape" and self.choosing_combat then
        self:hide_combat_choice()
    end

    -- Q to quit to menu
    if key == "q" then
        Anim.cancel_all()
        Event.clear()
        StateManager.switch("menu")
    end
end

-- === Avoid Room ===

function GamePlay:do_avoid_room()
    self.animating = true
    self:disable_room_interaction()

    -- Animate all 4 cards back to the dungeon pile
    local cards_to_move = {}
    for i = 1, 4 do
        if self.room_cards[i] then
            table.insert(cards_to_move, self.room_cards[i])
        end
    end

    local pending = #cards_to_move
    for _, card in ipairs(cards_to_move) do
        card:flip(function()
            card:move_to(C.DUNGEON_X, C.DUNGEON_Y, C.MOVE_DURATION, "ease_in_out_quad", function()
                card.visible = false
                pending = pending - 1
                if pending <= 0 then
                    self.room_cards = { nil, nil, nil, nil }
                    self.turn:avoid_room()
                    self:do_deal()
                end
            end)
        end)
    end
end

return GamePlay
