package main

import "github.com/JMitchell159/InvTrackBot/bot"

func main() {
	bot.Start()

	<-make(chan struct{})
}
