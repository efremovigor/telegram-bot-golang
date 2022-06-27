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
	key := fmt.Sprintf(redis.NextRequestMessageKey, query.GetUserId(), query.GetChatText())
	redis.Del(key)
	state, _ := redis.Get(fmt.Sprintf(redis.TranslateTransitionKey, query.GetChatId(), query.GetUserId()))
	messages := []telegram.RequestChannelTelegram{
		telegram.NewRequestChannelTelegram(
			"text",
			telegram.MergeRequestTelegram(
				telegram.GetHelloIGotYourMSGRequest(query),
				telegram.GetResultFromRapidMicrosoft(query, state),
			),
		),
	}

	if cambridgeInfo := cambridge.Get(query.GetChatText()); cambridgeInfo.IsValid() {
		for _, message := range telegram.GetResultFromCambridge(cambridgeInfo, query) {
			messages = append(messages, telegram.NewRequestChannelTelegram("text", message))
		}
		switch true {
		case helper.Len(cambridgeInfo.VoicePath.UK) > 0 && helper.Len(cambridgeInfo.VoicePath.US) > 0:
			messages = append(messages, telegram.NewRequestChannelTelegram("voice", telegram.NewCambridgeRequestTelegramVoice(cambridgeInfo.RequestText, query.GetChatId(), []string{telegram.CountryUk, telegram.CountryUs})))
		case helper.Len(cambridgeInfo.VoicePath.UK) > 0:
			messages = append(messages, telegram.NewRequestChannelTelegram("voice", telegram.NewCambridgeRequestTelegramVoice(cambridgeInfo.RequestText, query.GetChatId(), []string{telegram.CountryUk})))
		case helper.Len(cambridgeInfo.VoicePath.US) > 0:
			messages = append(messages, telegram.NewRequestChannelTelegram("voice", telegram.NewCambridgeRequestTelegramVoice(cambridgeInfo.RequestText, query.GetChatId(), []string{telegram.CountryUs})))
		}
		statistic.Consider(query.GetChatText(), query.GetUserId())
	}

	if cambridgeFounded := cambridge.Search(query.GetChatText()); cambridgeFounded.IsValid() {

	}

	if multitranInfo := multitran.Get(query.GetChatText()); multitranInfo.IsValid() {
		for _, message := range telegram.GetResultFromMultitran(multitranInfo, query) {
			messages = append(messages, telegram.NewRequestChannelTelegram("text", message))
		}
	}

	if requestTelegramInJson, err := json.Marshal(telegram.UserRequest{Request: query.GetChatText(), Output: messages}); err == nil {
		redis.Set(key, requestTelegramInJson, time.Hour*24)
	} else {
		fmt.Println(err)
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
		return
	}
	if err := json.Unmarshal([]byte(state), &request); err != nil {
		fmt.Println("Unmarshal request : " + err.Error())
	}

	if len(request.Output) > 0 {
		message = request.Output[0]
		if len(request.Output) > 1 {
			message.Buttons = []telegram.Keyboard{{Text: "more", CallbackData: telegram.NextRequestMessage + " " + word}}
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
