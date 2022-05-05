package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"telegram-bot-golang/env"
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
	url := "https://microsoft-translator-text.p.rapidapi.com/translate?to=ru&from=en&api-version=3.0&profanityAction=NoAction&textType=plain"

	payload := strings.NewReader("[\n    {\n        \"Text\": \"" + body.Message.Text + "\"\n    }\n]")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "microsoft-translator-text.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", env.GetEnvVariable("MICROSOFT_API_TOKEN"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	var microsoftTranslateResponse []MicrosoftTranslate

	if err := json.NewDecoder(res.Body).Decode(&microsoftTranslateResponse); err != nil {
		fmt.Println("could not decode microsoft response", err)
	}

	stringTranslation := ""
	for i, response := range microsoftTranslateResponse {
		for _, translation := range response.Translations {
			if i != 0 {
				stringTranslation += ", "
			}
			stringTranslation += translation.Text
		}
	}
	return SendMessageReqBody{
		ChatID:    body.Message.Chat.ID,
		Text:      fmt.Sprintf("Hey, [%s](tg://user?id=%d), I got your message: %s, translate:%s", body.Message.From.FirstName, body.Message.From.ID, body.Message.Text, stringTranslation),
		ParseMode: "MarkdownV2",
	}
}

type MicrosoftTranslate struct {
	Translations []Translation `json:"translations"`
}
type Translation struct {
	Text string `json:"text"`
	To   string `json:"to"`
}
