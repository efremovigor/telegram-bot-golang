package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func SayHello(body telegram.WebhookMessage, listener telegram.TelegramListener) {
	listener.Msg <- telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
	)
}

func General(body telegram.WebhookMessage, listener telegram.TelegramListener) {
	fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
	state, _ := redis.Get(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()))
	listener.Msg <- telegram.GetHelloIGotYourMSGRequest(body)
	listener.Msg <- telegram.GetResultFromRapidMicrosoft(body, state)
	telegram.GetResultFromCambridge(body, listener)
}

func Help(body telegram.WebhookMessage, listener telegram.TelegramListener) {
	telegram.SendMessage(telegram.GetTelegramRequest(
		body.GetChatId(),
		"*List of commands available to you:*\n"+
			telegram.GetRowSeparation()+
			"*"+telegram.DecodeForTelegram(RuEnCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc))+
			"*"+telegram.DecodeForTelegram(EnRuCommand)+fmt.Sprintf("* \\- Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc))+
			"*"+telegram.DecodeForTelegram(HelpCommand)+"* \\- Show all the available commands\n"+
			"*"+telegram.DecodeForTelegram(GetAllTopCommand)+"* \\- To see the most popular requests for translation or explanation  \n"+
			"*"+telegram.DecodeForTelegram(GetMyTopCommand)+"* \\- To see your popular requests for translation or explanation  \n",
	))
}
