package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (st *state) reset(s *discordgo.Session, e *discordgo.MessageCreate) {
	err := st.db.ResetServers(context.Background())
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "Failed to reset servers.")
		if err != nil {
			fmt.Println("Failed to send failed servers reset response:", err)
		}
	} else {
		_, err = s.ChannelMessageSend(e.ChannelID, "Reset servers.")
		if err != nil {
			fmt.Println("Failed to send servers reset response:", err)
		}
	}
	err = st.db.ResetItems(context.Background())
	if err != nil {
		_, err = s.ChannelMessageSend(e.ChannelID, "Failed to reset items.")
		if err != nil {
			fmt.Println("Failed to send failed items reset response:", err)
		}
	} else {
		_, err = s.ChannelMessageSend(e.ChannelID, "Reset items.")
		if err != nil {
			fmt.Println("Failed to send items reset response:", err)
		}
	}
}
