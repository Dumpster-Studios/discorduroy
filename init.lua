-- prototype discord bridge for minetest
local discord_token = ""
local channel_id = ""

local line = 1

for lines in io.lines(minetest.get_modpath("discord") .. "/auth.txt") do
	if line == 1 then
		discord_token = tostring(lines)
	else
		channel_id = tostring(lines)
		break
	end
	line = line + 1
end

local ie = minetest.request_insecure_environment()
local old_require = require

require = ie.require

ie.package.path = package.path .. ";" .. minetest.get_modpath(minetest.get_current_modname()) .. "/?.lua"

local discord_c = require "discord.init"
require = old_require

discord = {}
mt_client = discord_c.Client:new()

discord.accepted_token = false
if mt_client:loginWithToken(discord_token) then
	mt_client:sendMessage('[Server] Server started.', channel_id)
	discord.accepted_token = true
end

function discord.send_message(message)
	if discord.accepted_token then
		mt_client:sendMessage(message, channel_id)
	end
end

minetest.register_on_chat_message(function(name, message)
	if discord.accepted_token then
		mt_client:sendMessage("<**" .. name .. "**> " .. minetest.strip_colors(message), channel_id)
	end
end)

minetest.register_on_joinplayer(function(player)
	if discord.accepted_token then
		mt_client:sendMessage("**" .. player:get_player_name() .. "**" .. " joined the game.", channel_id)
	end
end)

minetest.register_on_leaveplayer(function(player)
	if discord.accepted_token then
		mt_client:sendMessage("**" .. player:get_player_name() .. "**" .. " left the game.", channel_id)
	end
end)

minetest.register_on_shutdown(function()
	if discord.accepted_token then
		mt_client:sendMessage("[Server] Server shutting down.", channel_id)
	end
end)

-- Destroy the file on game boot:
local log_file = minetest.get_modpath("discord").."/discord.log"
local chat = {}
local increment = 0
local discord_colour = "#7289DA"
os.remove(log_file)

local function get_discord_messages()
	local log = io.open(log_file)
	if log == nil then
	else
		log:close()
		local num_lines = 0
		for line in io.lines(log_file) do
			-- Deserialise chat information coming from discord.go
			local desel = minetest.deserialize(line)
			if chat[num_lines] == nil then
				if desel.message == "" then
					-- Create an empty line if the deserialised message text is empty
					-- we do this because uploading an image or file with no text causes
					-- an anomalous message with an empty string
					chat[num_lines] = ""
				else
					-- Colourise the [Discord] text and the <user_name> of the Discord member
					-- according to their role.
					chat[num_lines] = minetest.colorize(discord_colour, desel.server) .. " " ..
										minetest.colorize(desel.colour, desel.nick) .. " " .. desel.message
				end
			end
			num_lines = num_lines + 1
		end

		for i=increment, #chat do
			-- Fixes empty lines being sent to chat
			if chat[i] == "" then
			else
				minetest.chat_send_all(chat[i])
			end
		end
		increment = #chat + 1
	end
	minetest.after(1, get_discord_messages)
end

minetest.after(2, get_discord_messages)