package bot

import (
	"github.com/bwmarrin/discordgo"
)

func provideState(st *state, handler func(*discordgo.Session, *discordgo.MessageCreate, *state)) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, e *discordgo.MessageCreate) {
		handler(s, e, st)
	}
}
