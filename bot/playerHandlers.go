package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/JMitchell159/InvTrackBot/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func (st *state) addPlayer(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/* Syntax:
	!addPlayer <game_name> <player_name>
	*/
	if len(cmdArgs) < 2 {
		sendMessage(s, e.ChannelID, "The add command takes 2 arguments in this order, player name & game name.", "Failed sending required add arguments response:")
		return
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(cmdArgs[1]) {
		sendMessage(s, e.ChannelID, "Player name can only contain alphanumeric characters.", "Failed sending invalid name response:")
		return
	}
	game, err := st.db.GetGameByName(context.Background(), database.GetGameByNameParams{
		Name:     cmdArgs[0],
		ServerID: e.GuildID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		sendMessage(s, e.ChannelID, "Specified game does not exist in this server.", "Failed sending invalid game response:")
		return
	}
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching game.", "Failed sending fetching error response:")
		return
	}
	_, err = st.db.GetPlayerByName(context.Background(), database.GetPlayerByNameParams{
		Name:   cmdArgs[1],
		GameID: game.ID,
	})
	if err == nil {
		sendMessage(s, e.ChannelID, "That player has already been added to this game.", "Failed sending duplicate player response:")
		return
	}
	player, err := st.db.CreatePlayer(context.Background(), database.CreatePlayerParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmdArgs[1],
		GameID:    game.ID,
	})
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while trying to add player.", "Failed sending failed player add response:")
		return
	}
	sendMessage(s, e.ChannelID, fmt.Sprintf("Added player %s w/ ID: %s.", player.Name, player.ID.String()), "Failed sending add player response:")
}

func (st *state) listPlayers(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/*Syntax:
	!listPlayers <game_name>*/
	players, err := st.db.GetPlayersByGame(context.Background(), database.GetPlayersByGameParams{
		Name:     cmdArgs[0],
		ServerID: e.GuildID,
	})
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while trying to fetch players.", "Failed sending failed player fetch response:")
		return
	}
	if len(players) == 0 {
		sendMessage(s, e.ChannelID, "There are no players in this game.", "Failed sending no players response:")
		return
	}
	msg := fmt.Sprintf("Players in %s:\n", cmdArgs[0])
	for _, p := range players {
		msg = fmt.Sprintf("%s- %s (ID: %s)\n", msg, p.Name, p.ID.String())
	}
	sendMessage(s, e.ChannelID, msg, "Failed to send players list response:")
}
