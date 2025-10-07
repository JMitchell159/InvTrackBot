package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JMitchell159/InvTrackBot/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func (st *state) register(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/* register server Syntax:
	!register server
	*/
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

	/* register game Syntax:
	!register game <game_name>
	*/
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

	/*Register item Syntax:
	!register item name*/
	if strings.ToLower(cmdArgs[0]) == "item" {
		_, err := st.db.GetItem(context.Background(), cmdArgs[1])
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("%s item already exists.", cmdArgs[1]))
			if err != nil {
				fmt.Println("Failed sending duplicate item response:", err)
			}
			return
		}
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Something went wrong while trying to fetch item.")
			if err != nil {
				fmt.Println("Failed sending failed item fetch response:", err)
			}
			return
		}
		item, err := st.db.CreateItem(context.Background(), database.CreateItemParams{
			Name:      cmdArgs[1],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			_, err = s.ChannelMessageSend(e.ChannelID, "Something went wrong while trying to register item.")
			if err != nil {
				fmt.Println("Failed sending failed item register response:", err)
			}
			return
		}
		_, err = s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Registered %s item.", item.Name))
		if err != nil {
			fmt.Println("Failed sending register item response:", err)
		}
	}
}
