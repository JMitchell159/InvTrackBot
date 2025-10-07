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
			sendMessage(s, e.ChannelID, "This server has already been regitered.", "Failed sending duplicate server error response:")
			return
		}
		server, err := st.db.CreateServer(context.Background(), database.CreateServerParams{
			ID:        e.GuildID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Failed to register server.", "Failed sending failed registration response:")
			return
		}
		sendMessage(s, e.ChannelID, fmt.Sprintf("InventoryTracker resgistered for Server %s.", server.ID), "Failed sending registration response:")
		return
	}

	/* register game Syntax:
	!register game <game_name>
	*/
	if strings.ToLower(cmdArgs[0]) == "game" {
		if len(cmdArgs) < 2 {
			sendMessage(s, e.ChannelID, "Name must be specified when registering a game.", "Failed sending name required response:")
			return
		}
		if cmdArgs[1][0] == '@' {
			sendMessage(s, e.ChannelID, "Name argument cannot start with an @.", "Failed sending invalid name response:")
			return
		}
		_, err := st.db.GetServer(context.Background(), e.GuildID)
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "You must register the server first. The syntax for that command is '!register server'", "Failed sending server registration required response:")
			return
		}
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching server.", "Failed sending fetching error response:")
			return
		}
		_, err = st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
			Name:     cmdArgs[1],
			ServerID: e.GuildID,
		})
		if err == nil {
			sendMessage(s, e.ChannelID, "This game has already been registered.", "Failed sending duplicate game error response:")
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
			sendMessage(s, e.ChannelID, "Failed to register game.", "Failed sending failed game registration response:")
			return
		}
		sendMessage(s, e.ChannelID, fmt.Sprintf("Registered Game %s with ID: %s", game.Name, game.ID.String()), "Failed sending game registration response:")
		return
	}

	/*Register item Syntax:
	!register item name*/
	if strings.ToLower(cmdArgs[0]) == "item" {
		_, err := st.db.GetItem(context.Background(), cmdArgs[1])
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "This item already exists.", "Failed sending duplicate item response:")
			return
		}
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while trying to fetch item.", "Failed sending failed item fetch repsonse:")
			return
		}
		item, err := st.db.CreateItem(context.Background(), database.CreateItemParams{
			Name:      cmdArgs[1],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while trying to register item.", "Failed sending failed item register response:")
			return
		}
		sendMessage(s, e.ChannelID, fmt.Sprintf("Registered %s item.", item.Name), "Failed sending register item response:")
		return
	}

	sendMessage(s, e.ChannelID, fmt.Sprintf("Unknown command register %s.", cmdArgs[0]), "Failed sending Unknown Command response:")
}
