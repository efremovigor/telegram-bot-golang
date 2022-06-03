package command

import (
	"telegram-bot-golang/telegram"
)

func Draft(body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
	)
}
