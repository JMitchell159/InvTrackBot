package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
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
		if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[1]) {
			sendMessage(s, e.ChannelID, "Game name can only contain alphanumeric characters.", "Failed sending invalid name response:")
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
	!register item <item_name> <game_name>*/
	if strings.ToLower(cmdArgs[0]) == "item" {
		if len(cmdArgs) < 3 {
			sendMessage(s, e.ChannelID, "Register item command requires 2 arguments, the item name and the game name.", "Failed sending required arguments response:")
			return
		}
		if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[1]) || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[2]) {
			sendMessage(s, e.ChannelID, "Item name and game name arguments can only contain alphanumeric characters.", "Failed sending invalid arguments response:")
			return
		}
		game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
			Name:     cmdArgs[2],
			ServerID: e.GuildID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "This game does not exist in this server.", "Failed sending non-existent game response:")
			return
		}
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching game.", "Failed sending fetching error response:")
			return
		}
		nameWithHash := fmt.Sprintf("%s[%s]", cmdArgs[1], game.ID.String())
		_, err = st.db.GetItem(context.Background(), nameWithHash)
		if err == nil {
			sendMessage(s, e.ChannelID, "This item already exists.", "Failed sending duplicate item response:")
			return
		}
		item, err := st.db.CreateItem(context.Background(), database.CreateItemParams{
			Name:      nameWithHash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while trying to register item.", "Failed sending failed item register response:")
			return
		}
		split := strings.Split(item.Name, "[")
		sendMessage(s, e.ChannelID, fmt.Sprintf("Registered %s item.", split[0]), "Failed sending register item response:")
		return
	}

	/*Register item w/ desc Syntax:
	!register itemDesc <item_name> <game_name> <description>*/
	if strings.ToLower(cmdArgs[0]) == "itemdesc" {
		if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[1]) || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[2]) {
			sendMessage(s, e.ChannelID, "Item name and game name arguments can only contain alphanumeric characters.", "Failed sending invalid arguments response:")
			return
		}
		game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
			Name:     cmdArgs[2],
			ServerID: e.GuildID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "This game does not exist in this server.", "Failed sending non-existent game response:")
			return
		}
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching game.", "Failed sending fetching error response:")
			return
		}
		nameWithHash := fmt.Sprintf("%s[%s]", cmdArgs[1], game.ID.String())
		_, err = st.db.GetItem(context.Background(), nameWithHash)
		if err == nil {
			sendMessage(s, e.ChannelID, "This item already exists.", "Failed sending duplicate item response:")
			return
		}
		item, err := st.db.CreateItemWDesc(context.Background(), database.CreateItemWDescParams{
			Name:      nameWithHash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Description: sql.NullString{
				Valid:  true,
				String: cmdArgs[2],
			},
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while trying to register item.", "Failed sending failed item register resposne:")
			return
		}
		split := strings.Split(item.Name, "[")
		sendMessage(s, e.ChannelID, fmt.Sprintf("Registered %s item with description: %s", split[0], item.Description.String), "Failed sending register item response:")
		return
	}

	/*Register item w/ cat Syntax:
	!register itemCat <item_name> <game_name> <category>*/
	if strings.ToLower(cmdArgs[0]) == "itemcat" {
		if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[1]) || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[2]) {
			sendMessage(s, e.ChannelID, "Item name and game name arguments can only contain alphanumeric characters.", "Failed sending invalid arguments response:")
			return
		}
		game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
			Name:     cmdArgs[2],
			ServerID: e.GuildID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "This game does not exist in this server.", "Failed sending non-existent game response:")
			return
		}
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching game.", "Failed sending fetching error response:")
			return
		}
		nameWithHash := fmt.Sprintf("%s[%s]", cmdArgs[1], game.ID.String())
		_, err = st.db.GetItem(context.Background(), nameWithHash)
		if err == nil {
			sendMessage(s, e.ChannelID, "This item already exists.", "Failed sending duplicate item response:")
			return
		}
		item, err := st.db.CreateItemWCat(context.Background(), database.CreateItemWCatParams{
			Name:      nameWithHash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Category: sql.NullString{
				Valid:  true,
				String: cmdArgs[2],
			},
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while trying to register item.", "Failed sending failed item register response:")
			return
		}
		split := strings.Split(item.Name, "[")
		sendMessage(s, e.ChannelID, fmt.Sprintf("Registered item %s with category: %s", split[0], item.Category.String), "Failed sending register item response:")
		return
	}

	/*Register full item Syntax:
	!register itemFull <item_name> <game_name> <category> <description>*/
	if strings.ToLower(cmdArgs[0]) == "itemfull" {
		if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[1]) || !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[2]) {
			sendMessage(s, e.ChannelID, "Item name and game name arguments can only contain alphanumeric characters.", "Failed sending invalid arguments response:")
			return
		}
		game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
			Name:     cmdArgs[2],
			ServerID: e.GuildID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			sendMessage(s, e.ChannelID, "This game does not exist in this server.", "Failed sending non-existent game response:")
			return
		}
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while fetching game.", "Failed sending fetching error response:")
			return
		}
		nameWithHash := fmt.Sprintf("%s[%s]", cmdArgs[1], game.ID.String())
		_, err = st.db.GetItem(context.Background(), nameWithHash)
		if err == nil {
			sendMessage(s, e.ChannelID, "This item already exists.", "Failed sending duplicate item response:")
			return
		}
		item, err := st.db.CreateItemFull(context.Background(), database.CreateItemFullParams{
			Name:      nameWithHash,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Category: sql.NullString{
				Valid:  true,
				String: cmdArgs[2],
			},
			Description: sql.NullString{
				Valid:  true,
				String: cmdArgs[3],
			},
		})
		if err != nil {
			sendMessage(s, e.ChannelID, "Something went wrong while trying to register item.", "Failed sending failed item register response:")
			return
		}
		split := strings.Split(item.Name, "[")
		sendMessage(s, e.ChannelID, fmt.Sprintf("Registered item %s with category: %s\nand description: %s", split[0], item.Category.String, item.Description.String), "Failed sending register item response:")
		return
	}

	sendMessage(s, e.ChannelID, fmt.Sprintf("Unknown command register %s.", cmdArgs[0]), "Failed sending Unknown Command response:")
}
