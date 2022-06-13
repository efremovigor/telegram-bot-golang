package main

import (
	"telegram-bot-golang/http"
	"telegram-bot-golang/telegram"
)

func main() {
	listener := telegram.TelegramListener{Msg: make(chan telegram.TelegramTree, 100)}
	go http.Handle(listener)
	go telegram.HandleRequests(listener)
	select {}
}
