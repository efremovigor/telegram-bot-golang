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
	EditedMessage struct {
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
	} `json:"edited_message"`
	UpdateId int `json:"update_id"`
}

func (body WebhookReqBody) GetChatId() int {
	if body.Message.Chat.ID != 0 {
		return body.Message.Chat.ID
	} else {
		return body.EditedMessage.Chat.ID
	}
}
func (body WebhookReqBody) GetChatText() string {
	if body.Message.Chat.ID != 0 {
		return body.Message.Text
	} else {
		return body.EditedMessage.Text
	}
}

func (body WebhookReqBody) GetUsername() string {
	if body.Message.Chat.ID != 0 {
		return body.Message.From.FirstName
	} else {
		return body.EditedMessage.From.FirstName
	}
}
func (body WebhookReqBody) GetUserId() int {
	if body.Message.Chat.ID != 0 {
		return body.Message.From.ID
	} else {
		return body.EditedMessage.From.ID
	}
}

type SendMessageReqBody struct {
	ChatID      int         `json:"chat_id"`
	Text        string      `json:"text"`
	ParseMode   string      `json:"parse_mode"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

type ReplyMarkup struct {
	Keyboard        [][]Keyboard `json:"keyboard"`
	OneTimeKeyboard bool         `json:"one_time_keyboard"`
	ResizeKeyboard  bool         `json:"resize_keyboard"`
}

type Keyboard struct {
	Text string `json:"text"`
}

type MicrosoftTranslate struct {
	Translations []Translation `json:"translations"`
}
type Translation struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

var Chats = map[int]map[int]string{}

func SayHello(body WebhookReqBody) SendMessageReqBody {
	state, exist := Chats[body.GetChatId()][body.GetUserId()]

	var from, to string

	if !exist || state == "en_ru" {
		from = "en"
		to = "ru"
	} else {
		from = "ru"
		to = "en"
	}

	url := fmt.Sprintf("https://microsoft-translator-text.p.rapidapi.com/translate?to=%s&from=%s&api-version=3.0&profanityAction=NoAction&textType=plain", to, from)

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
	return GetTelegramRequest(
		body.GetChatId(),
		GetBaseMsg(body.GetUsername(), body.GetUserId())+GetTranslateMsg(body.GetChatText(), stringTranslation),
	)
}

func GetBaseMsg(name string, id int) string {
	return fmt.Sprintf("Hey, [%s](tg://user?id=%d), ", name, id)
}
func GetTranslateMsg(base string, translate string) string {
	return fmt.Sprintf("I got your message: %s, translate:%s", DecodeForTelegram(base), DecodeForTelegram(translate))
}

func GetChangeTranslateMsg(translate string) string {
	return fmt.Sprintf("I changed translation:  %s", DecodeForTelegram(translate))
}

func GetTelegramRequest(chatId int, text string) SendMessageReqBody {
	return SendMessageReqBody{
		ChatID:      chatId,
		Text:        text,
		ParseMode:   "MarkdownV2",
		ReplyMarkup: ReplyMarkup{Keyboard: [][]Keyboard{{{Text: "Hello"}}, {{Text: "Привет"}}}, OneTimeKeyboard: true, ResizeKeyboard: true},
	}
}

func DecodeForTelegram(text string) string {
	replacer := strings.NewReplacer(
		">", "\\>",
		"<", "\\<",
		".", "\\.",
		"-", "\\-",
		"!", "\\!",
		"#", "\\#",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"=", "\\=",
	)
	return replacer.Replace(text)
}
