package command

import (
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/telegram"
)

func GetTop10(query telegram.TelegramQueryInterface) telegram.RequestTelegramText {
	text := ""
	list, err := model.GetWordStatistics(10)
	if err == nil {
		text = telegram.GetRatingHeader(10, true)
		text += handleList(list)
	}
	return telegram.RequestTelegramText{Text: text, ChatId: query.GetChatId()}
}

func GetTop10ForUser(query telegram.TelegramQueryInterface) telegram.RequestTelegramText {
	text := ""
	list, err := model.GetWordStatisticsForUser(10, query.GetUserId())
	if err == nil {
		text = telegram.GetRatingHeader(10, false)
		text += handleList(list)
	}
	return telegram.RequestTelegramText{Text: text, ChatId: query.GetChatId()}
}

func handleList(list []model.WordStatistic) string {
	var text string
	for k, statistic := range list {
		text += telegram.GetRowRating(k+1, statistic)
	}
	return text
}
