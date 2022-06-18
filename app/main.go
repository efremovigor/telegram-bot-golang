package main

import (
	"encoding/json"
	"fmt"
	"telegram-bot-golang/http"
	"telegram-bot-golang/telegram"
)

func main() {
	str := "{\"Text\":\"Hey, [Igor](tg://user?id=184357122)\\n\\nI got your message: *get*\\n\",\"ChatId\":184357122}"
	var inter interface{}
	var inter1 telegram.RequestTelegramText
	inter = str
	qw := fmt.Sprintf("%s", inter)
	fmt.Println(string(qw))

	err := json.Unmarshal([]byte(qw), &inter1)
	if err != nil {
		fmt.Println(err)
	}
	qw1, _ := json.Marshal(inter1)
	fmt.Println(string(qw1))

	listener := telegram.Listener{Message: make(chan telegram.RequestChannelTelegram, 100)}
	go http.Handle(listener)
	go telegram.HandleRequests(listener)
	select {}
}
