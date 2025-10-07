package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JMitchell159/InvTrackBot/internal/database"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func (st *state) addPlayer(s *discordgo.Session, e *discordgo.MessageCreate, cmdArgs []string) {
	/* Syntax:
	!addPlayer <game_name> <player_name>
	*/
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
		Name:     cmdArgs[0],
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
		Name:   cmdArgs[1],
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
		Name:      cmdArgs[1],
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
