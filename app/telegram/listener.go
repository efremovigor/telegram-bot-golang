package telegram

func HandleRequests(listener Listener) {
	for {
		select {
		case request := <-listener.Message:
			switch request.Type {
			case "text":
				sendMessage(request.Message.(RequestTelegramText))
			case "voice":
				voice := request.Message.(CambridgeRequestTelegramVoice)
				sendVoices(voice.ChatId, voice.Info)
			}
		}
	}
}
