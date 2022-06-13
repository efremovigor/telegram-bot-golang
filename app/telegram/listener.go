package telegram

func HandleRequests(listener Listener) {
	for {
		select {
		case request := <-listener.Message:
			switch request.Type {
			case "text":
				sendMessage(request.Message.(SendMessageReqBody))
			case "voice":
				voice := request.Message.(CambridgeTelegramVoice)
				sendVoices(voice.ChatId, voice.Info)
			}
		}
	}
}
