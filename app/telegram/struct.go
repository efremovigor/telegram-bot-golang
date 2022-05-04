package telegram

import "fmt"

type WebhookReqBody struct {
	Message struct {
		Text      string `json:"text"`
		MessageId int    `json:"message_id"`
		Chat      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Type      string `json:"type"`
			Username  string `json:"username"`
		} `json:"chat"`
		Date int `json:"date"`
		From struct {
			ID           int    `json:"id"`
			FirstName    string `json:"first_name"`
			IsBot        bool   `json:"is_bot"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			Username     string `json:"username"`
		} `json:"from"`
	} `json:"message"`
	UpdateId int `json:"update_id"`
}

type SendMessageReqBody struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func sayHello(body WebhookReqBody) SendMessageReqBody {
	return SendMessageReqBody{
		ChatID:    body.Message.Chat.ID,
		Text:      fmt.Sprintf("Hey, [%s](tg://user?id=%d), I got your message: %s", body.Message.From.FirstName, body.Message.From.ID, body.Message.Text),
		ParseMode: "MarkdownV2",
	}
}
