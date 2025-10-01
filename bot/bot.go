package bot

import (
	"fmt"
	"strings"

	"github.com/JMitchell159/InvTrackBot/config"
	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Failed reading configuration:", err)
		return
	}

	goBot, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("Failed initializing Discord Session:", err)
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println("Failed getting current User:", err)
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println("Failed opening connectionto Discord:", err)
		return
	}

	fmt.Println("Bot is now connected!")
}

func messageHandler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.ID == BotId {
		return
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Failed reading configuration:", err)
		return
	}

	prefix := cfg.BotPrefix
	if strings.HasPrefix(e.Content, prefix) {
		args := strings.Fields(e.Content)[strings.Index(e.Content, prefix):]
		cmd := args[0][len(prefix):]
		//		Not in use right now, but will be used later
		//		arguments := args[1:]

		switch cmd {
		case "ping":
			_, err := s.ChannelMessageSend(e.ChannelID, "Pong!")
			if err != nil {
				fmt.Println("Failed sending Pong response:", err)
			}
		default:
			_, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Unknown command %q.", cmd))
			if err != nil {
				fmt.Println("Failed sending Unknown Command response:", err)
			}
		}
	}
}
