package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/telegram"
)

func SayHello(body telegram.WebhookMessage) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
	)
}

func General(body telegram.WebhookMessage) []telegram.TelegramTree {
	fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
	state, _ := redis.Get(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()))
	msgs := []telegram.TelegramTree{
		{Type: "text", Msg: telegram.GetHelloIGotYourMSGRequest(body)},
		{Type: "text", Msg: telegram.GetResultFromRapidMicrosoft(body, state)},
	}

	cambridgeInfo := cambridge.Get(body.GetChatText())
	if cambridgeInfo.IsValid() {
		for _, msg := range telegram.GetResultFromCambridge(cambridgeInfo, body) {
			msgs = append(msgs, telegram.TelegramTree{Type: "text", Msg: msg})
		}
		msgs = append(msgs, telegram.TelegramTree{Type: "voice", Msg: telegram.TelegramCambridgeVoice{Info: cambridgeInfo, ChatId: body.GetChatId()}})
	}
	return msgs
}

func Help(body telegram.WebhookMessage) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		"*List of commands available to you:*\n"+
			telegram.GetRowSeparation()+
			"*"+telegram.DecodeForTelegram(RuEnCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc))+
			"*"+telegram.DecodeForTelegram(EnRuCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc))+
			"*"+telegram.DecodeForTelegram(HelpCommand)+"* \\- Show all the available commands\n"+
			"*"+telegram.DecodeForTelegram(GetAllTopCommand)+"* \\- To see the most popular requests for translation or explanation  \n"+
			"*"+telegram.DecodeForTelegram(GetMyTopCommand)+"* \\- To see your popular requests for translation or explanation  \n",
	)
}
