package telegram

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/text/language"
)

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

func SayHello(body WebhookReqBody) SendMessageReqBody {
	ctx := context.Background()
	lang, err := language.Parse(body.Message.Text)
	if err != nil {
		fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		fmt.Errorf("error initialization translate client")
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{body.Message.Text}, lang, nil)
	if err != nil {
		fmt.Errorf("error initialization translate client")

	}
	if len(resp) == 0 {
		fmt.Errorf("Translate returned empty response to text: %s", body.Message.Text)
	}

	return SendMessageReqBody{
		ChatID:    body.Message.Chat.ID,
		Text:      fmt.Sprintf("Hey, [%s](tg://user?id=%d), I got your message: %s, translate:", body.Message.From.FirstName, body.Message.From.ID, body.Message.Text, resp[0].Text),
		ParseMode: "MarkdownV2",
	}
}
