package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Token for both auth
var Token string = ""

// Channel ID
var Channel string = ""

// Server ID
var Server string = ""

func main() {

	line := 1
	file, err := os.Open("auth.txt")
	if err != nil {
		_, err = os.OpenFile("auth.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Println("Failed to load auth.txt, a file named auth.txt has been created.\nThe file should contain the following lines:\n\n<Discord bot token>\n<Discord text channel ID>\n<Discord server ID>")
		return
	}
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		switch line {
		case 1:
			Token = fscanner.Text()
		case 2:
			Channel = fscanner.Text()
		case 3:
			Server = fscanner.Text()
		}
		line = line + 1
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error using provided auth token.", err)
		return
	}

	dg.StateEnabled = true

	// Add a handler that manages reading of incoming messages.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Failed to open websocket to Discord.,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit. \nUse GNU screen or tmux to keep the heartbeat alive; \notherwise the heartbeat may randomly terminate.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func getRoleColour(userID string, s *discordgo.Session) int {
	member, _ := s.GuildMember(Server, userID)
	for _, roleID := range member.Roles {
		role, err := s.State.Role(Server, roleID)
		if err == nil {
			return role.Color
		}
	}
	return int(16777215)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == Channel {
		member, err := s.State.Member(Server, m.Author.ID)
		userColour := "#ffffff"

		member, err = s.GuildMember(Server, m.Author.ID)
		roleColour := getRoleColour(m.Author.ID, s)
		if roleColour != 0 {
			blue := roleColour & 0xFF
			green := (roleColour >> 8) & 0xFF
			red := (roleColour >> 16) & 0xFF
			userColour = "#" + strconv.FormatInt(int64(red), 16) + strconv.FormatInt(int64(green), 16) + strconv.FormatInt(int64(blue), 16)
		}

		serverNick := m.Author.Username
		if member.Nick != "" {
			serverNick = member.Nick
		}

		messageAppend := "return {[\"server\"] = \"[Discord]\", [\"colour\"] = \"" + userColour + "\", [\"nick\"] = \"<" + serverNick + ">\", [\"message\"] = \"" + strings.ReplaceAll(m.Content, "\n", " ") + "\"}\n"

		if m.Message.Attachments != nil {
			for k := range m.Message.Attachments {
				messageAppend = messageAppend + "return {[\"server\"] = \"[Discord]\", [\"colour\"] = \"" + userColour + "\", [\"nick\"] = \"<" + serverNick + ">\", [\"message\"] = \"" + m.Message.Attachments[k].URL + "\"}\n"
			}
		}

		// Lua may randomly delete this file so always attempt to re-create it.
		f, err := os.OpenFile("discord.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}

		f.WriteString(messageAppend)
		f.Close()
	}
}
