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

	goBot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuildMembers
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
		args, err := parseUserInput(e.Content)
		if err != nil {
			sendMessage(s, e.ChannelID, "something went wrong while trying to parse input.", "Failed to send failed input parse response:")
			return
		}

		cmd := args[0][len(prefix):]
		switch cmd {
		case "reset":
			st.reset(s, e)
			return
		case "listGames":
			st.listGames(s, e)
			return
		}
		if len(args) < 2 {
			sendMessage(s, e.ChannelID, "All commands must have at least 2 fields, the command name and at least one qualifier.", "Failed sending required arguments response:")
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
		case "updateItem":
			st.updateItem(s, e, arguments)
		case "listPlayers":
			st.listPlayers(s, e, arguments)
		case "listItem":
			st.listItem(s, e, arguments)
		case "listItems":
			st.listItems(s, e, arguments)
		case "listInventory":
			st.listInventory(s, e, arguments)
		case "listInvByCat":
			st.listInventoryByCat(s, e, arguments)
		default:
			sendMessage(s, e.ChannelID, fmt.Sprintf("Unknown command %q.", cmd), "Failed sending Unknown Command response:")
		}
	}
}
