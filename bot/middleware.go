package bot

import (
	"github.com/JMitchell159/InvTrackBot/config"
	"github.com/bwmarrin/discordgo"
)

func providePrefix(cfg *config.Config, msgHandler func(*discordgo.Session, *discordgo.MessageCreate, string)) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, e *discordgo.MessageCreate) {
		msgHandler(s, e, cfg.BotPrefix)
	}
}
