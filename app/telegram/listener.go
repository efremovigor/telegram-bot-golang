package telegram

func HandleRequests(listener TelegramListener) {
	for {
		select {
		case request := <-listener.Msg:
			switch request.Type {
			case "text":
				SendMessage(request.Msg.(SendMessageReqBody))
			case "voice":
				voice := request.Msg.(TelegramCambridgeVoice)
				SendVoices(voice.ChatId, voice.Info)
			}
		}
	}
}
