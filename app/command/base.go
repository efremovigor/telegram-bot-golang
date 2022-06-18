package command

import (
	"encoding/json"
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/telegram"
)

func SayHello(body telegram.WebhookMessage) telegram.RequestTelegramText {
	return telegram.RequestTelegramText{
		Text:   telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
		ChatId: body.GetChatId(),
	}
}

func General(body telegram.WebhookMessage) {
	fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
	state, _ := redis.Get(fmt.Sprintf(redis.TranslateTransitionKey, body.GetChatId(), body.GetUserId()))
	messages := []telegram.RequestChannelTelegram{
		{Type: "text", Message: telegram.GetHelloIGotYourMSGRequest(body)},
		{Type: "text", Message: telegram.GetResultFromRapidMicrosoft(body, state)},
	}

	cambridgeInfo := cambridge.Get(body.GetChatText())
	if cambridgeInfo.IsValid() {
		for _, message := range telegram.GetResultFromCambridge(cambridgeInfo, body) {
			messages = append(messages, telegram.RequestChannelTelegram{Type: "text", Message: message})
		}
		messages = append(messages, telegram.RequestChannelTelegram{Type: "voice", Message: telegram.CambridgeRequestTelegramVoice{Info: cambridgeInfo, ChatId: body.GetChatId()}})
	}
	multitranInfo := multitran.Get(body.GetChatText())
	if multitranInfo.IsValid() {
		for _, message := range telegram.GetResultFromMultitran(multitranInfo, body) {
			messages = append(messages, telegram.RequestChannelTelegram{Type: "text", Message: message})
		}
	}
	if messagesTelegramInJson, err := json.Marshal(messages); err == nil {
		if requestTelegramInJson, err := json.Marshal(telegram.UserRequest{Request: body.GetChatText(), Output: messagesTelegramInJson}); err == nil {
			redis.Set(fmt.Sprintf(redis.NextRequestMessageKey, body.GetUserId()), requestTelegramInJson)
		} else {
			fmt.Println(err)
		}
	}
}

func GetNextMessage(userId int) (message telegram.RequestChannelTelegram, err error) {
	var request telegram.UserRequest
	var messages []telegram.RequestChannelTelegram
	state, _ := redis.Get(fmt.Sprintf(redis.NextRequestMessageKey, userId))
	if err := json.Unmarshal([]byte(state), &request); err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(request.Output, &messages); err != nil {
		fmt.Println(err)
	}

	message = messages[0]
	if len(messages[1:]) > 0 {
		if request.Output, err = json.Marshal(messages[1:]); err == nil {
			fmt.Println(err)
		}
		if infoInJson, err := json.Marshal(request); err == nil {
			redis.Set(fmt.Sprintf(redis.NextRequestMessageKey, userId), infoInJson)
		} else {
			fmt.Println(err)
		}
	} else {
		redis.Del(fmt.Sprintf(redis.NextRequestMessageKey, userId))
	}

	return message, err
}

func Help(body telegram.WebhookMessage) telegram.RequestTelegramText {
	return telegram.RequestTelegramText{
		Text: "*List of commands available to you:*\n" +
			telegram.GetRowSeparation() +
			"*" + telegram.DecodeForTelegram(RuEnCommand) + fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc)) +
			"*" + telegram.DecodeForTelegram(EnRuCommand) + fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc)) +
			"*" + telegram.DecodeForTelegram(AutoTranslateCommand) + "* \\- Change translation automatic \n" +
			"*" + telegram.DecodeForTelegram(HelpCommand) + "* \\- Show all the available commands\n" +
			"*" + telegram.DecodeForTelegram(GetAllTopCommand) + "* \\- To see the most popular requests for translation or explanation  \n" +
			"*" + telegram.DecodeForTelegram(GetMyTopCommand) + "* \\- To see your popular requests for translation or explanation  \n",
		ChatId: body.GetChatId(),
	}
}
