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
					sendBaseInfo(textRequest, request.HasMore)
				} else {
					fmt.Println(err)
				}
			case "voice":
				var textRequest CambridgeRequestTelegramVoice
				if err := json.Unmarshal(request.Message, &textRequest); err == nil {
					sendVoiceMessage(textRequest)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}
