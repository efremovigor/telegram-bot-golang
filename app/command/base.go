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

func General(chatId int, userId int, chatText string) {
	state, _ := redis.Get(fmt.Sprintf(redis.TranslateTransitionKey, chatId, userId))
	collector := telegram.Collector{Type: telegram.ReasonTypeNextMessage}
	collector.Add(
		telegram.GetResultFromRapidMicrosoft(chatId, chatText, state),
	)

	if page := cambridge.Get(chatText); page.IsValid() {
		statistic.Consider(chatText, userId)
		collector.Add(
			handleCambridgePage(page, chatId, chatText)...,
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
				telegram.DecodeForTelegram("Additional various 🔽"),
				chatId,
				buttons,
			),
		)
	}

	if page := multitran.Get(chatText); page.IsValid() {
		collector.Add(telegram.GetResultFromMultitran(page, chatId, chatText)...)
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextMessageKey, userId, chatText), chatText, collector.GetMessageForSave())
}

func ListShortInfo(chatId int, userId int, chatText string) {
	var collector telegram.Collector

	if page := cambridge.Get(chatText); page.IsValid() {
		statistic.Consider(chatText, userId)
		collector.Add(
			handleCambridgePage(page, chatId, chatText)...,
		)
	}
}

func FullInfo(dictionary string, chatId int, userId int, chatText string) {
	var collector telegram.Collector
	switch dictionary {
	case "cambridge":
		collector.Type = telegram.ReasonFullCambridgeMessage
		if page := cambridge.Get(chatText); page.IsValid() {
			collector.Add(
				handleCambridgePage(page, chatId, chatText)...,
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
					telegram.DecodeForTelegram("Additional various 🔽"),
					chatId,
					buttons,
				),
			)
		}
	case "multitran":
		collector.Type = telegram.ReasonFullMultitranMessage
		if page := multitran.Get(chatText); page.IsValid() {
			collector.Add(telegram.GetResultFromMultitran(page, chatId, chatText)...)
		}
	case "wooordhunt":
		/** wooordhunt hundler */
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextFullInfoRequestMessageKey, dictionary, userId, chatText), chatText, collector.GetMessageForSave())

}

func handleCambridgePage(page cambridge.Page, chatId int, chatText string) (messages []telegram.RequestTelegramText) {
	messages = telegram.GetResultFromCambridge(page, chatId, chatText)
	return
}

func saveMessagesQueue(key string, chatText string, messages []telegram.RequestChannelTelegram) {
	redis.Del(key)
	redis.SetStruct(key, telegram.UserRequest{Request: chatText, Output: messages}, time.Hour*24)
}

func GetSubCambridge(chatId int, userId int, chatText string) {
	cambridgeFounded, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeSearchValue, chatText))
	if err != nil {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Don't find cambridge key - word:%s", chatText))
		return
	}
	if page := cambridge.DoRequest(chatText, cambridge.Url+cambridgeFounded, ""); page.IsValid() {
		statistic.Consider(chatText, userId)
		collector := telegram.Collector{Type: telegram.ReasonSubCambridgeMessage}
		collector.Add(handleCambridgePage(page, chatId, chatText)...)
		saveMessagesQueue(
			fmt.Sprintf(redis.SubCambridgeMessageKey, userId, chatText),
			chatText,
			collector.GetMessageForSave(),
		)

		return
	} else {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Aren't able to parse url:%s", cambridge.Url+cambridgeFounded))
	}
}

func SendVoice(query telegram.IncomingTelegramQueryInterface, lang string, hash string) {
	url, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeUniqVoiceLink, hash))
	if err != nil {
		fmt.Println(err)
		return
	}
	var cache redis.VoiceFile
	if err := json.Unmarshal([]byte(url), &cache); err != nil {
		fmt.Println(err)
		return
	}
	telegram.SendVoices(query.GetChatId(), lang, cache)
}

func SendImage(query telegram.IncomingTelegramQueryInterface, hash string) {
	url, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeUniqPicLink, hash))
	if err != nil {
		fmt.Println(err)
		return
	}
	var cache redis.PicFile
	if err := json.Unmarshal([]byte(url), &cache); err != nil {
		fmt.Println(err)
		return
	}
	telegram.SendImage(query.GetChatId(), cache)

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
