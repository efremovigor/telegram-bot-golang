package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

var transitions = map[string]struct {
	key  string
	desc string
}{RuEnCommand: {key: "ru_en", desc: "RU -> EN"}, EnRuCommand: {key: "en_ru", desc: "EN -> RU"}}

func ChangeTranslateTransition(command string, body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	redis.Set(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()), transitions[command].key)

	return telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.GetBaseMsg(body.GetUsername(), body.GetUserId())+telegram.GetChangeTranslateMsg(transitions[command].desc),
	)
}
