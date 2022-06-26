package command

import (
	"encoding/json"
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/statistic"
	"telegram-bot-golang/telegram"
)

func SayHello(query telegram.TelegramQueryInterface) telegram.RequestTelegramText {
	return telegram.RequestTelegramText{
		Text:   telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
		ChatId: query.GetChatId(),
	}
}

func General(query telegram.TelegramQueryInterface) {
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
		messages = append(messages, telegram.NewRequestChannelTelegram("voice", telegram.CambridgeRequestTelegramVoice{Info: cambridgeInfo, ChatId: query.GetChatId()}))
		statistic.Consider(query.GetChatText(), query.GetUserId())

	}

	if multitranInfo := multitran.Get(query.GetChatText()); multitranInfo.IsValid() {
		for _, message := range telegram.GetResultFromMultitran(multitranInfo, query) {
			messages = append(messages, telegram.NewRequestChannelTelegram("text", message))
		}
	}

	if requestTelegramInJson, err := json.Marshal(telegram.UserRequest{Request: query.GetChatText(), Output: messages}); err == nil {
		redis.Set(key, requestTelegramInJson)
	} else {
		fmt.Println(err)
	}
}

func GetNextMessage(userId int, word string) (message telegram.RequestChannelTelegram, err error) {
	var request telegram.UserRequest
	key := fmt.Sprintf(redis.NextRequestMessageKey, userId, word)
	state, _ := redis.Get(key)
	if err := json.Unmarshal([]byte(state), &request); err != nil {
		fmt.Println("Unmarshal request : " + err.Error())
	}

	if len(request.Output) > 0 {
		message = request.Output[0]
		if len(request.Output) > 1 {
			message.HasMore = true
			request.Output = request.Output[1:]
			if infoInJson, err := json.Marshal(request); err == nil {
				redis.Set(key, infoInJson)
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

func Help(query telegram.TelegramQueryInterface) telegram.RequestTelegramText {
	return telegram.RequestTelegramText{
		Text: "*List of commands available to you:*\n" +
			telegram.GetRowSeparation() +
			"*" + telegram.DecodeForTelegram(RuEnCommand) + fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc)) +
			"*" + telegram.DecodeForTelegram(EnRuCommand) + fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc)) +
			"*" + telegram.DecodeForTelegram(AutoTranslateCommand) + "* \\- Change translation automatic \n" +
			"*" + telegram.DecodeForTelegram(HelpCommand) + "* \\- Show all the available commands\n" +
			"*" + telegram.DecodeForTelegram(GetAllTopCommand) + "* \\- To see the most popular requests for translation or explanation  \n" +
			"*" + telegram.DecodeForTelegram(GetMyTopCommand) + "* \\- To see your popular requests for translation or explanation  \n",
		ChatId: query.GetChatId(),
	}
}
