package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func ChangeTranslateTransition(command string, query telegram.IncomingTelegramQueryInterface) telegram.RequestTelegramText {
	key := fmt.Sprintf(redis.TranslateTransitionKey, query.GetChatId(), query.GetUserId())
	if command == AutoTranslateCommand {
		redis.Del(key)
	} else {
		redis.Set(key, Transitions()[command].key, 0)
	}
	return telegram.MakeRequestTelegramText(
		query.GetChatText(),
		telegram.GetBaseMsg(query.GetUsername(), query.GetUserId())+telegram.GetChangeTranslateMsg(Transitions()[command].Desc),
		query.GetChatId(),
		[]telegram.Keyboard{},
	)
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
