package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func sendMessage(s *discordgo.Session, channelID, discordMessage, errMessage string) {
	_, err := s.ChannelMessageSend(channelID, discordMessage)
	if err != nil {
		fmt.Println(errMessage, err)
	}
}

func parseUserInput(str string) ([]string, error) {
	splitQuotes := strings.Split(str, "\"")
	if len(splitQuotes) == 1 {
		return strings.Fields(str), nil
	}
	splitStart := strings.Fields(splitQuotes[0])
	splitStart = append(splitStart, splitQuotes[len(splitQuotes)-2])
	return splitStart, nil
}
