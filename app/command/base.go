package command

import (
	"encoding/json"
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/helper"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/statistic"
	"telegram-bot-golang/telegram"
	"time"
)

func SayHello(query telegram.IncomingTelegramQueryInterface) telegram.RequestTelegramText {
	return telegram.MakeRequestTelegramText(
		query.GetChatText(),
		telegram.DecodeForTelegram("Hello friend. How can I help you?"),
		query.GetChatId(),
		[]telegram.Keyboard{},
	)
}

func ListShortInfo(chatId int, userId int, chatText string) {
	state, _ := redis.Get(fmt.Sprintf(redis.TranslateTransitionKey, chatId, userId))
	collector := telegram.Collector{Type: telegram.ReasonTypeNextShortMessage}
	collector.Add(
		telegram.GetResultFromRapidMicrosoft(chatId, chatText, state),
	)

	if page := cambridge.Get(chatText); page.IsValid() {
		statistic.Consider(chatText, userId)
		collector.Add(
			telegram.GetCambridgeShortInfo(chatId, page)...,
		)
	}

	if search := cambridge.Search(chatText); search.IsValid() {
		var buttons []telegram.Keyboard
		for _, founded := range search.Founded {
			buttons = append(buttons, telegram.Keyboard{Text: "🔍 " + founded.Word, CallbackData: telegram.SearchRequest + " cambridge " + founded.Word})
		}
		collector.Add(
			telegram.MakeRequestTelegramText(
				search.RequestWord,
				telegram.DecodeForTelegram("✅ *Cambridge-dictionary-search*: \n\n Additional various 🔽"),
				chatId,
				buttons,
			),
		)
	}

	if page := multitran.Get(chatText); page.IsValid() {
		collector.Add(telegram.GetMultitranShortInfo(chatId, page)...)
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextShortInfoRequestMessageKey, userId, chatText), chatText, collector.GetMessageForSave())
}

func saveMessagesQueue(key string, chatText string, messages []telegram.RequestChannelTelegram) {
	redis.Del(key)
	redis.SetStruct(key, telegram.UserRequest{Request: chatText, Output: messages}, time.Hour*24)
}

func GetNextMessage(key string, word string) (message telegram.RequestChannelTelegram, err error) {
	var request telegram.UserRequest
	value, err := redis.Get(key)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Don't find key when we getting next message: word:%s", word))
		return
	}
	if err := json.Unmarshal([]byte(value), &request); err != nil {
		fmt.Println("Unmarshal request : " + err.Error())
	}

	if len(request.Output) > 0 {
		message = request.Output[0]
		fmt.Println(helper.ToJson(message))
		if len(request.Output) > 1 {
			request.Output = request.Output[1:]
			redis.SetStruct(key, request, time.Hour*24)
		} else {
			redis.Del(key)
		}
	} else {
		redis.Del(key)
	}

	return message, err
}

func GetCountMessages(key string) int {
	var request telegram.UserRequest
	value, err := redis.Get(key)
	if err != nil {
		return 0
	}
	if err := json.Unmarshal([]byte(value), &request); err != nil {
		return 0
	}
	return len(request.Output)
}

func Help(query telegram.IncomingTelegramQueryInterface) telegram.RequestTelegramText {
	return telegram.MakeRequestTelegramText(
		query.GetChatText(),
		"*List of commands available to you:*\n"+
			telegram.GetRowSeparation()+
			"*"+telegram.DecodeForTelegram(RuEnCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc))+
			"*"+telegram.DecodeForTelegram(EnRuCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc))+
			"*"+telegram.DecodeForTelegram(AutoTranslateCommand)+"* \\- Change translation automatic \n"+
			"*"+telegram.DecodeForTelegram(HelpCommand)+"* \\- Show all the available commands\n"+
			"*"+telegram.DecodeForTelegram(GetAllTopCommand)+"* \\- To see the most popular requests for translation or explanation  \n"+
			"*"+telegram.DecodeForTelegram(GetMyTopCommand)+"* \\- To see your popular requests for translation or explanation  \n",
		query.GetChatId(),
		[]telegram.Keyboard{},
	)
}
