package main

import (
	"telegram-bot-golang/http"
	"telegram-bot-golang/telegram"
)

func main() {
	listener := telegram.Listener{Message: make(chan telegram.RequestChannelTelegram, 100)}
	go http.Handle(listener)
	go telegram.HandleRequests(listener)
	select {}
}
