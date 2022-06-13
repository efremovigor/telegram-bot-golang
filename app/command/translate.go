package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func ChangeTranslateTransition(command string, body telegram.WebhookMessage, listener telegram.TelegramListener) {
	redis.Set(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()), Transitions()[command].key)

	listener.Msg <- telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.GetBaseMsg(body.GetUsername(), body.GetUserId())+telegram.GetChangeTranslateMsg(Transitions()[command].Desc),
	)
}
func Transitions() map[string]struct {
	key  string
	Desc string
} {
	return map[string]struct {
		key  string
		Desc string
	}{RuEnCommand: {key: "ru_en", Desc: "RU -> EN"}, EnRuCommand: {key: "en_ru", Desc: "EN -> RU"}}
}
