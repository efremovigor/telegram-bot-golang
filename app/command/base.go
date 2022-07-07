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

func General(query telegram.IncomingTelegramQueryInterface) {
	state, _ := redis.Get(fmt.Sprintf(redis.TranslateTransitionKey, query.GetChatId(), query.GetUserId()))
	var collector telegram.Collector
	collector.Add(
		"text",
		telegram.GetResultFromRapidMicrosoft(query, state),
	)

	if page := cambridge.Get(query.GetChatText()); page.IsValid() {
		collector.Add(
			"text",
			handleCambridgePage(page, query.GetUserId(), query.GetChatId(), query.GetChatText())...,
		)
	}

	if search := cambridge.Search(query.GetChatText()); search.IsValid() {
		var buttons []telegram.Keyboard
		for _, founded := range search.Founded {
			buttons = append(buttons, telegram.Keyboard{Text: "🔍 " + founded.Word, CallbackData: telegram.SearchRequest + " cambridge " + founded.Word})
		}
		collector.Add(
			"text",
			telegram.MakeRequestTelegramText(
				search.RequestWord,
				telegram.DecodeForTelegram("Additional various 🔽"),
				query.GetChatId(),
				[]telegram.Keyboard{},
			),
		)

		fmt.Println(helper.ToJson(search))
	}

	if page := multitran.Get(query.GetChatText()); page.IsValid() {
		collector.Add("text", telegram.GetResultFromMultitran(page, query)...)
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextRequestMessageKey, query.GetUserId(), query.GetChatText()), query.GetChatText(), collector.Messages)
}

func handleCambridgePage(page cambridge.Page, userId int, chatId int, chatText string) (messages []telegram.RequestTelegramText) {
	messages = telegram.GetResultFromCambridge(page, chatId, chatText)
	statistic.Consider(chatText, userId)
	return
}

func saveMessagesQueue(key string, chatText string, messages []telegram.RequestChannelTelegram) {
	redis.Del(key)
	redis.SetStruct(key, telegram.UserRequest{Request: chatText, Output: messages}, time.Hour*24)
}

func GetSubCambridge(query telegram.IncomingTelegramQueryInterface) {
	cambridgeFounded, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeSearchValue, query.GetChatText()))
	if err != nil {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Don't find cambridge key - word:%s", query.GetChatText()))
		return
	}
	if page := cambridge.DoRequest(query.GetChatText(), cambridge.Url+cambridgeFounded, ""); page.IsValid() {
		var collector telegram.Collector
		collector.Add("text", handleCambridgePage(page, query.GetUserId(), query.GetChatId(), query.GetChatText())...)
		saveMessagesQueue(
			fmt.Sprintf(redis.NextCambridgeRequestMessageKey, query.GetUserId(), query.GetChatText()),
			query.GetChatText(),
			collector.Messages,
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

func GetNextMessage(userId int, word string) (message telegram.RequestChannelTelegram, err error) {
	var request telegram.UserRequest
	key := fmt.Sprintf(redis.NextRequestMessageKey, userId, word)
	state, err := redis.Get(key)
	if err != nil {
		key = fmt.Sprintf(redis.NextCambridgeRequestMessageKey, userId, word)
		state, err = redis.Get(key)
		if err != nil {
			fmt.Println(fmt.Sprintf("[Strange behaivor] Don't find key when we getting next message: word:%s", word))
			return
		}
	}
	if err := json.Unmarshal([]byte(state), &request); err != nil {
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
