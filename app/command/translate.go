package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func ChangeTranslateTransition(command string, body telegram.WebhookMessage) telegram.RequestTelegramText {
	if command == AutoTranslateCommand {
		redis.Del(fmt.Sprintf(redis.TranslateTransitionKey, body.GetChatId(), body.GetUserId()))
	} else {
		redis.Set(fmt.Sprintf(redis.TranslateTransitionKey, body.GetChatId(), body.GetUserId()), Transitions()[command].key)
	}
	return telegram.RequestTelegramText{Text: telegram.GetBaseMsg(body.GetUsername(), body.GetUserId()) + telegram.GetChangeTranslateMsg(Transitions()[command].Desc), ChatId: body.GetChatId()}
}
func Transitions() map[string]struct {
	key  string
	Desc string
} {
	return map[string]struct {
		key  string
		Desc string
	}{
		RuEnCommand:          {key: "ru_en", Desc: "RU -> EN"},
		EnRuCommand:          {key: "en_ru", Desc: "EN -> RU"},
		AutoTranslateCommand: {key: "auto_translate", Desc: "AUTO"},
	}
}
