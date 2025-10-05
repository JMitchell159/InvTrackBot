package main

import (
	"github.com/JMitchell159/InvTrackBot/bot"
	_ "github.com/lib/pq"
)

func main() {
	bot.Start()

	<-make(chan struct{})
}
