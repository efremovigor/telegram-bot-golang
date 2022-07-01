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
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
		query.GetChatId(),
	)
}

func General(query telegram.IncomingTelegramQueryInterface) {
	state, _ := redis.Get(fmt.Sprintf(redis.TranslateTransitionKey, query.GetChatId(), query.GetUserId()))
	messages := []telegram.RequestChannelTelegram{
		telegram.NewRequestChannelTelegram(
			"text",
			telegram.MergeRequestTelegram(
				telegram.GetHelloIGotYourMSGRequest(query),
				telegram.GetResultFromRapidMicrosoft(query, state),
			),
			[]telegram.Keyboard{},
		),
	}

	if page := cambridge.Get(query.GetChatText()); page.IsValid() {
		messages = append(messages, handleCambridgePage(page, query.GetUserId(), query.GetChatId(), query.GetChatText())...)
	}

	if search := cambridge.Search(query.GetChatText()); search.IsValid() {
		var buttons []telegram.Keyboard
		for _, founded := range search.Founded {
			buttons = append(buttons, telegram.Keyboard{Text: "🔍 " + founded.Word, CallbackData: telegram.SearchRequest + " cambridge " + founded.Word})
		}
		messages = append(messages, telegram.NewRequestChannelTelegram(
			"text",
			telegram.MakeRequestTelegramText(
				search.RequestWord,
				telegram.DecodeForTelegram("Maybe you look for it:"),
				query.GetChatId(),
			),
			buttons))

		fmt.Println(helper.ToJson(search))
	}

	if page := multitran.Get(query.GetChatText()); page.IsValid() {
		for _, message := range telegram.GetResultFromMultitran(page, query) {
			messages = append(messages, telegram.NewRequestChannelTelegram("text", message, []telegram.Keyboard{}))
		}
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextRequestMessageKey, query.GetUserId(), query.GetChatText()), query.GetChatText(), messages)
}

func handleCambridgePage(page cambridge.Page, userId int, chatId int, chatText string) (messages []telegram.RequestChannelTelegram) {
	for _, message := range telegram.GetResultFromCambridge(page, chatId, chatText) {
		messages = append(messages, telegram.NewRequestChannelTelegram("text", message, []telegram.Keyboard{}))
	}
	switch true {
	case helper.Len(page.VoicePath.UK) > 0 && helper.Len(page.VoicePath.US) > 0:
		messages = append(messages, telegram.NewRequestChannelVoiceTelegram(page.RequestText, chatId, []string{telegram.CountryUk, telegram.CountryUs}))
	case helper.Len(page.VoicePath.UK) > 0:
		messages = append(messages, telegram.NewRequestChannelVoiceTelegram(page.RequestText, chatId, []string{telegram.CountryUk}))
	case helper.Len(page.VoicePath.US) > 0:
		messages = append(messages, telegram.NewRequestChannelVoiceTelegram(page.RequestText, chatId, []string{telegram.CountryUs}))
	}
	statistic.Consider(chatText, userId)
	return
}

func saveMessagesQueue(key string, chatText string, messages []telegram.RequestChannelTelegram) {
	redis.Del(key)

	if requestTelegramInJson, err := json.Marshal(telegram.UserRequest{Request: chatText, Output: messages}); err == nil {
		redis.Set(key, requestTelegramInJson, time.Hour*24)
	} else {
		fmt.Println(err)
	}
}

func GetSubCambridge(query telegram.IncomingTelegramQueryInterface) {
	cambridgeFounded, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeSearchValue, query.GetChatText()))
	if err != nil {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Don't find cambridge key - word:%s", query.GetChatText()))
		return
	}
	if page := cambridge.DoRequest(query.GetChatText(), cambridge.Url+cambridgeFounded, ""); page.IsValid() {
		saveMessagesQueue(
			fmt.Sprintf(redis.NextCambridgeRequestMessageKey, query.GetUserId(), query.GetChatText()),
			query.GetChatText(),
			handleCambridgePage(page, query.GetUserId(), query.GetChatId(), query.GetChatText()),
		)

		return
	} else {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Aren't able to parse url:%s", cambridge.Url+cambridgeFounded))
	}
}

func GetVoice(query telegram.IncomingTelegramQueryInterface, lang string, word string) {
	if cambridgeInfo := cambridge.Get(word); cambridgeInfo.IsValid() {
		var request telegram.UserRequest
		key := fmt.Sprintf(redis.NextRequestMessageKey, query.GetUserId(), word)
		state, _ := redis.Get(key)
		if err := json.Unmarshal([]byte(state), &request); err != nil {
			fmt.Println("Unmarshal request : " + err.Error())
		}
		telegram.SendVoices(query.GetChatId(), cambridgeInfo, lang, len(request.Output) > 0)
	}
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
		if len(request.Output) > 1 {
			message.Buttons = append(message.Buttons, telegram.Keyboard{Text: "📚 more", CallbackData: telegram.NextRequestMessage + " " + word})
			request.Output = request.Output[1:]
			if infoInJson, err := json.Marshal(request); err == nil {
				redis.Set(key, infoInJson, time.Hour*24)
			} else {
				fmt.Println(err)
			}
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
	)
}
