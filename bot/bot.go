package bot

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/JMitchell159/InvTrackBot/config"
	"github.com/JMitchell159/InvTrackBot/internal/database"
	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func Start() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Failed reading configuration:", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Println("Failed opening database connection:", err)
	}

	dbQueries := database.New(db)

	s := &state{
		db:  dbQueries,
		cfg: cfg,
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

	goBot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	goBot.StateEnabled = true

	BotId = u.ID

	goBot.AddHandler(provideState(s, messageHandler))

	err = goBot.Open()
	if err != nil {
		fmt.Println("Failed opening connection to Discord:", err)
		return
	}

	fmt.Println("Bot is now connected!")
}

func messageHandler(s *discordgo.Session, e *discordgo.MessageCreate, st *state) {
	if e.Author.ID == BotId {
		return
	}

	prefix := st.cfg.BotPrefix

	if strings.HasPrefix(e.Content, prefix) {
		args := strings.Fields(e.Content)[strings.Index(e.Content, prefix):]
		cmd := args[0][len(prefix):]
		if len(args) < 2 {
			_, err := s.ChannelMessageSend(e.ChannelID, "All commands must have at least 2 fields, the command name and at least one qualifier.")
			if err != nil {
				fmt.Println("Failed sending required arguments response:", err)
			}
			return
		}
		arguments := args[1:]

		switch cmd {
		case "register":
			st.register(s, e, arguments)
		case "addPlayer":
			st.addPlayer(s, e, arguments)
		case "addItem":
			st.addItem(s, e, arguments)
		default:
			_, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Unknown command %q.", cmd))
			if err != nil {
				fmt.Println("Failed sending Unknown Command response:", err)
			}
		}
	}
}
