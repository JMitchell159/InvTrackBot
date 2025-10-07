package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func sendMessage(s *discordgo.Session, channelID, discordMessage, errMessage string) {
	_, err := s.ChannelMessageSend(channelID, discordMessage)
	if err != nil {
		fmt.Println(errMessage, err)
	}
}
