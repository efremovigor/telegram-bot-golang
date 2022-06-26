package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func ChangeTranslateTransition(command string, query telegram.TelegramQueryInterface) telegram.RequestTelegramText {
	if command == AutoTranslateCommand {
		redis.Del(fmt.Sprintf(redis.TranslateTransitionKey, query.GetChatId(), query.GetUserId()))
	} else {
		redis.Set(fmt.Sprintf(redis.TranslateTransitionKey, query.GetChatId(), query.GetUserId()), Transitions()[command].key, 0)
	}
	return telegram.RequestTelegramText{Text: telegram.GetBaseMsg(query.GetUsername(), query.GetUserId()) + telegram.GetChangeTranslateMsg(Transitions()[command].Desc), ChatId: query.GetChatId()}
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
