package command

import (
	"fmt"
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/telegram"
)

func GetTop10(query telegram.IncomingTelegramQueryInterface) telegram.RequestTelegramText {
	text := ""
	list, err := model.GetWordStatistics(10)
	if err == nil {
		text = telegram.GetRatingHeader(10, true)
		text += handleList(list)
	} else {
		fmt.Println(err)
	}
	return telegram.MakeRequestTelegramText(query.GetChatText(), text, query.GetChatId(), []telegram.Keyboard{})
}

func GetTop10ForUser(query telegram.IncomingTelegramQueryInterface) telegram.RequestTelegramText {
	text := ""
	list, err := model.GetWordStatisticsForUser(10, query.GetUserId())
	if err == nil {
		text = telegram.GetRatingHeader(10, false)
		text += handleList(list)
	} else {
		fmt.Println(err)
	}
	return telegram.MakeRequestTelegramText(query.GetChatText(), text, query.GetChatId(), []telegram.Keyboard{})
}

func handleList(list []model.WordStatistic) string {
	var text string
	for k, statistic := range list {
		text += telegram.GetRowRating(k+1, statistic)
	}
	return text
}
