# Discord for Minetest
____
######_Part of the SolarSail family._
### Notice:

This mod requires LuaJIT; attempts at running it as a standard Lua interpreter will result in Minetest failing to execute this mod.

### Preparation:

`git clone` or manually download a zip of this folder and extract it to `minetest/mods` resulting in `minetest/mods/discord`.

Install cURL. You can find this in your package manager, in the case of Windows, this is provided for you.

### Golang Bot Setup Procedure:

Install the latest version of [Golang](https://golang.org/).

Using the terminal, execute `go get github.com/bwmarrin/discordgo`

Afterwards, you can test that the Golang is installed and functioning correctly for your environment by running `go run main.go` inside the `discord` folder.

By default, the bot will complain of a missing `auth.txt` and generate an empty file for you.

The error, will always come out as this in the terminal:

```
$ go run main.go
Failed to load auth.txt, a file named auth.txt has been created.
The file should contain the following lines:

<Discord bot token>
<Discord text channel ID>
<Discord server ID>
```

The auth.txt contents should look like this example, but not exactly as the token, server ID and channel ID can differ dramatically:

```
abcdefghijklmop1234.5
1234567890987654321
0987654321234567890
```

### Minetest Mod Setup Procedure

Add this mod to the chosen world via `worldmods` or by including it from the `mods` directory.

Grant this mod insecure environment access via `security.trusted_mods = discord`.

### Licensing:

This mod (`main.go` and `init.lua`) is licensed as MIT.

Discord.lua is licensed as MIT and can be found here: https://github.com/videah/discord.lua

`luajit-request`'s license can be found at `discord/discord/luajit-request/LICENSE`