package command

import (
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/telegram"
)

func GetTop10(body telegram.WebhookMessage) {
	text := ""
	list, err := model.GetWordStatistics(10)
	if err == nil {
		text = telegram.GetRatingHeader(10, true)
		text += handleList(list)
	}
	telegram.SendMessage(telegram.GetTelegramRequest(body.GetChatId(), text))
}

func GetTop10ForUser(body telegram.WebhookMessage) {
	text := ""
	list, err := model.GetWordStatisticsForUser(10, body.GetUserId())
	if err == nil {
		text = telegram.GetRatingHeader(10, false)
		text += handleList(list)
	}
	telegram.SendMessage(telegram.GetTelegramRequest(body.GetChatId(), text))
}

func handleList(list []model.WordStatistic) string {
	var text string
	for k, statistic := range list {
		text += telegram.GetRowRating(k+1, statistic)
	}
	return text
}
