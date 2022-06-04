package command

import (
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/telegram"
)

func SayHello(body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		telegram.DecodeForTelegram("Hello Friend. How can I help you?"),
	)
}

func General(body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	fmt.Println(fmt.Sprintf("chat text: %s", body.GetChatText()))
	state, _ := redis.Get(fmt.Sprintf("chat_%d_user_%d", body.GetChatId(), body.GetUserId()))
	return telegram.Reply(body, state)
}

func Help(body telegram.WebhookReqBody) telegram.SendMessageReqBody {
	return telegram.GetTelegramRequest(
		body.GetChatId(),
		"*List of commands available to you:*\n"+
			telegram.GetRowSeparation()+
			RuEnCommand+fmt.Sprintf("Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[RuEnCommand].Desc))+
			EnRuCommand+fmt.Sprintf("Change translate of transition %s \n", telegram.DecodeForTelegram(Transitions()[EnRuCommand].Desc))+
			HelpCommand+"Show all the available commands\n"+
			GetAllTopCommand+"To see the most popular requests for translation or explanation  \n"+
			GetMyTopCommand+"To see your popular requests for translation or explanation  \n",
	)
}
