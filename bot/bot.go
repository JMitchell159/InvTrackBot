package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JMitchell159/InvTrackBot/config"
	"github.com/JMitchell159/InvTrackBot/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
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
			return
		case "addPlayer":
			st.addPlayer(s, e, arguments)
			return
		default:
			_, err := s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Unknown command %q.", cmd))
			if err != nil {
				fmt.Println("Failed sending Unknown Command response:", err)
			}
		}

	}
}

func (st *state) register(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	// register server
	if strings.ToLower(cmdArgs[0]) == "server" {
		_, err := st.db.GetServer(context.Background(), e.GuildID)
		if err == nil {
			_, err := s.ChannelMessageSend(e.ChannelID, "This server has already been registered.")
			if err != nil {
				fmt.Println("Failed sending duplicate server error response:", err)
			}
			return
		}
		server, err := st.db.CreateServer(context.Background(), database.CreateServerParams{
			ID:        e.GuildID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to register InventoryTracker.")
			if err != nil {
				fmt.Println("Failed sending failed registration response:", err)
			}
			return
		}
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Tracker registered for Server %s.", server.ID))
		if err != nil {
			fmt.Println("Failed sending registration response:", err)
		}
		return
	}
	// register game
	if strings.ToLower(cmdArgs[0]) == "game" {
		if len(cmdArgs) < 2 {
			_, err := s.ChannelMessageSend(e.ChannelID, "Name must be specified when registering a game.")
			if err != nil {
				fmt.Println("Failed sending name required response:", err)
			}
			return
		}
		if cmdArgs[1][0] == '@' {
			_, err := s.ChannelMessageSend(e.ChannelID, "Name argument cannot start with an @.")
			if err != nil {
				fmt.Println("Failed sending name required response:", err)
			}
			return
		}
		_, err := st.db.GetServer(context.Background(), e.GuildID)
		if errors.Is(err, sql.ErrNoRows) {
			_, err := s.ChannelMessageSend(e.ChannelID, "You must register the server first. The syntax for that command is '!register server @[botName]'")
			if err != nil {
				fmt.Println("Failed sending server registration required response:", err)
			}
			return
		}
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Something went wrong when fetching server: %v", err))
			if err != nil {
				fmt.Println("Failed sending fetching error response:", err)
			}
		}
		_, err = st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
			Name:     cmdArgs[1],
			ServerID: e.GuildID,
		})
		if err == nil {
			_, err := s.ChannelMessageSend(e.ChannelID, "This game has already been registered.")
			if err != nil {
				fmt.Println("Failed sending duplicate game error response:", err)
			}
			return
		}
		game, err := st.db.CreateGame(context.Background(), database.CreateGameParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmdArgs[1],
			ServerID:  e.GuildID,
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Failed to register game.")
			if err != nil {
				fmt.Println("Failed sending failed game registration response:", err)
			}
			return
		}
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Registered Game %s with ID: %s", game.Name, game.ID.String()))
		if err != nil {
			fmt.Println("Failed sending game registration response:", err)
		}
	}
}

func (st *state) addPlayer(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	if cmdArgs == nil {
		_, err := s.ChannelMessageSend(e.ChannelID, "The add command takes 2 arguments in this order, player name & game name")
		if err != nil {
			fmt.Println("Failed sending required add arguments response:", err)
		}
		return
	}
	if len(cmdArgs) == 1 {
		_, err := s.ChannelMessageSend(e.ChannelID, "The add command takes 2 arguments in this order, player name & game name")
		if err != nil {
			fmt.Println("Failed sending required add arguments response:", err)
		}
		return
	}
	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[1],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		_, err = s.ChannelMessageSend(e.ChannelID, "Specified game does not exist in this server.")
		if err != nil {
			fmt.Println("Failed sending invalid game response:", err)
		}
		return
	}
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Something went wrong when fetching game: %v", err))
		if err != nil {
			fmt.Println("Failed sending fetching error response:", err)
		}
		return
	}
	_, err = st.db.GetPlayerByName(context.Background(), database.GetPlayerByNameParams{
		Name:   cmdArgs[0],
		GameID: game.ID,
	})
	if err == nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "That player has already been added to this game.")
		if err != nil {
			fmt.Println("Failed sending duplicate player response:", err)
		}
		return
	}
	player, err := st.db.CreatePlayer(context.Background(), database.CreatePlayerParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmdArgs[0],
		GameID:    game.ID,
	})
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Something went wrong when trying to add player: %v", err))
		if err != nil {
			fmt.Println("Failed sending botched add player response:", err)
		}
		return
	}
	_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Added player %s w/ ID: %s to game %s w/ ID: %s", player.Name, player.ID.String(), game.Name, game.ID.String()))
	if err != nil {
		fmt.Println("Failed sending add player response:", err)
	}
}
