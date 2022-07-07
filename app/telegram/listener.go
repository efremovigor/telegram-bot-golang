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
					sendBaseInfo(textRequest)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}
