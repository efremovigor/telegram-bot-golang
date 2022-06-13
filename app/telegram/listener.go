package telegram

func HandleRequests(listener TelegramListener) {
	for {
		select {
		case request := <-listener.Msg:
			SendMessage(request)
		}
	}
}
