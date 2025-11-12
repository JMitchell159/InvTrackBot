package bot

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (st *state) reset(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.ID != st.cfg.BotOwnerID {
		sendMessage(s, e.ChannelID, "You are not allowed to run this command.", "Failed to send unauthorized access response:")
		return
	}
	err := st.db.ResetServers(context.Background())
	if err != nil {
		sendMessage(s, e.ChannelID, "Failed to reset servers.", "Failed to send failed servers reset response:")
	} else {
		sendMessage(s, e.ChannelID, "Reset servers.", "Failed to send servers reset response:")
	}
	err = st.db.ResetItems(context.Background())
	if err != nil {
		sendMessage(s, e.ChannelID, "Failed to reset items.", "Failed to send failed items reset response:")
	} else {
		sendMessage(s, e.ChannelID, "Reset items.", "Failed to send items reset response:")
	}
}
