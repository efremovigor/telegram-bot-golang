package telegram

import (
	"encoding/json"
	"fmt"
)

func HandleRequests(listener Listener) {
	for {
		select {
		case request := <-listener.Message:
			switch request.Type {
			case "text":
				message := request.Message

				if infoInJson, err := json.Marshal(message); err == nil {
					fmt.Println(string(infoInJson))
				} else {
					fmt.Println(err)
				}
				sendMessage(message.(RequestTelegramText))
			case "voice":
				voice := request.Message.(CambridgeRequestTelegramVoice)
				sendVoices(voice.ChatId, voice.Info)
			}
		}
	}
}
