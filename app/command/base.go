package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	"telegram-bot-golang/telegram"
)

func SayHello(body telegram.WebhookMessage) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
	)
}

func General(body telegram.WebhookMessage) []telegram.RequestChannelTelegram {
	fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
	state, _ := redis.Get(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()))
	messages := []telegram.RequestChannelTelegram{
		{Type: "text", Message: telegram.GetHelloIGotYourMSGRequest(body)},
		{Type: "text", Message: telegram.GetResultFromRapidMicrosoft(body, state)},
	}

	cambridgeInfo := cambridge.Get(body.GetChatText())
	if cambridgeInfo.IsValid() {
		for _, message := range telegram.GetResultFromCambridge(cambridgeInfo, body) {
			messages = append(messages, telegram.RequestChannelTelegram{Type: "text", Message: message})
		}
		messages = append(messages, telegram.RequestChannelTelegram{Type: "voice", Message: telegram.CambridgeTelegramVoice{Info: cambridgeInfo, ChatId: body.GetChatId()}})
	}
	multitranInfo := multitran.Get(body.GetChatText())
	if multitranInfo.IsValid() {
		for _, message := range telegram.GetResultFromMultitran(multitranInfo, body) {
			messages = append(messages, telegram.RequestChannelTelegram{Type: "text", Message: message})
		}
	}

	return messages
}

func Help(body telegram.WebhookMessage) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		"*List of commands available to you:*\n"+
			telegram.GetRowSeparation()+
			"*"+telegram.DecodeForTelegram(RuEnCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc))+
			"*"+telegram.DecodeForTelegram(EnRuCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc))+
			"*"+telegram.DecodeForTelegram(AutoTranslateCommand)+"* \\- Change translation automatic \n"+
			"*"+telegram.DecodeForTelegram(HelpCommand)+"* \\- Show all the available commands\n"+
			"*"+telegram.DecodeForTelegram(GetAllTopCommand)+"* \\- To see the most popular requests for translation or explanation  \n"+
			"*"+telegram.DecodeForTelegram(GetMyTopCommand)+"* \\- To see your popular requests for translation or explanation  \n",
	)
}
