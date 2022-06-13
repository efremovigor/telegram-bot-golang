package telegram

func HandleRequests(listener TelegramListener) {
	for {
		select {
		case request := <-listener.Msg:
			switch request.Type {
			case "text":
				sendMessage(request.Msg.(SendMessageReqBody))
			case "voice":
				voice := request.Msg.(TelegramCambridgeVoice)
				sendVoices(voice.ChatId, voice.Info)
			}
		}
	}
}
