package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (st *state) listGames(s *discordgo.Session, e *discordgo.MessageCreate) {
	games, err := st.db.GetGamesByServer(context.Background(), e.GuildID)
	if err != nil {
		sendMessage(s, e.ChannelID, "Something went wrong while fetching games.", "Failed to send failed games fetch response:")
		return
	}
	if len(games) == 0 {
		sendMessage(s, e.ChannelID, "There are no games registered for this server.", "Failed to send no registered games response:")
		return
	}
	msg := "Games registered to this server:\n"
	for _, game := range games {
		msg = fmt.Sprintf("%s- %s\n", msg, game.Name)
	}
	sendMessage(s, e.ChannelID, msg, "Failed to send games list response:")
}
