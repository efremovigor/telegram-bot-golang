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
				if infoInJson, err := json.Marshal(request); err == nil {
					fmt.Println(infoInJson)
				} else {
					fmt.Println(err)
				}
				sendMessage(request.Message.(RequestTelegramText))
			case "voice":
				voice := request.Message.(CambridgeRequestTelegramVoice)
				sendVoices(voice.ChatId, voice.Info)
			}
		}
	}
}
