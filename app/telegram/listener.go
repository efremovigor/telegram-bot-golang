package telegram

import (
	"encoding/json"
	"fmt"
)

func HandleRequests(listener Listener) {
	for {
		select {
		case request := <-listener.Message:
			fmt.Println("message to telegram:" + string(request.Message))
			switch request.Type {
			case "text":
				var textRequest RequestTelegramText
				if err := json.Unmarshal(request.Message, &textRequest); err == nil {
					sendMessage(textRequest, request.HasMore)
				} else {
					fmt.Println(err)
				}
			case "voice":
				var voiceRequest CambridgeRequestTelegramVoice
				if err := json.Unmarshal(request.Message, &voiceRequest); err == nil {
					sendVoices(voiceRequest.ChatId, voiceRequest.Info, request.HasMore)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}
