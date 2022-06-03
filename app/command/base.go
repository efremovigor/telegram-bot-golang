package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func SayHello(body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
	)
}

func General(body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
	state, _ := redis.Get(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()))
	return telegram.Reply(body, state)
}
