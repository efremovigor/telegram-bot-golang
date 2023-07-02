package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/telegram"
)

func MakeMultitranFullIfEmpty(chatId int, userId int, chatText string) {
	collector := telegram.Collector{Type: telegram.ReasonFullMultitranMessage}
	if page := multitran.Get(chatText); page.IsValid() {
		collector.Add(telegram.GetResultFromMultitran(page, chatId, chatText)...)
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextFullInfoRequestMessageKey, "multitran", userId, chatText), chatText, collector.GetMessageForSave())
}
