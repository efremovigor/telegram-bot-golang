package telegram

import (
	"fmt"
	"strings"
	"telegram-bot-golang/service/dictionary/cambridge"
	rapid_microsoft "telegram-bot-golang/service/translate/rapid-microsoft"
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

func Reply(body WebhookReqBody, state string) SendMessageReqBody {
	var from, to string
	if state == "" || state == "en_ru" {
		from = "en"
		to = "ru"
	} else {
		from = "ru"
		to = "en"
	}

	return GetTelegramRequest(
		body.GetChatId(),
		GetBaseMsg(body.GetUsername(), body.GetUserId())+
			GetIGotYourNewRequest(body.GetChatText())+
			GetBlockWithRapidInfo(rapid_microsoft.GetTranslate(body.Message.Text, to, from))+
			GetBlockWithCambridge(cambridge.Get(body.Message.Text)),
	)
}

func GetBaseMsg(name string, id int) string {
	return fmt.Sprintf("Hey, [%s](tg://user?id=%d)\n", name, id) +
		DecodeForTelegram("-----\n")
}
func GetIGotYourNewRequest(base string) string {
	return fmt.Sprintf(
		"I got your message: %s\n", DecodeForTelegram(base))
}

func GetBlockWithRapidInfo(translate string) string {
	return fmt.Sprintf(
		DecodeForTelegram("Translates of rapid-microsoft:")+" %s\n\n", DecodeForTelegram(translate))
}

func GetBlockWithCambridge(info cambridge.Info) string {
	mainBlock := "Information from cambridge-dictionary:" + GetFieldIfCan(info.Text, "") + "\n"
	mainBlock += GetFieldIfCan(info.Type, "Type") + "\n"
	if len(info.Explanation) > 0 {
		mainBlock += "Explanations:\n"
		for n, explanation := range info.Explanation {
			mainBlock += fmt.Sprintf("%d", n) + ".\n"
			mainBlock += GetFieldIfCan(explanation.Level, "Level") + "\n"
			mainBlock += GetFieldIfCan(explanation.SemanticDescription, "Semantic") + "\n"
			mainBlock += GetFieldIfCan(explanation.Translate, "Translate") + "\n"
			if len(explanation.Example[0]) > 0 {
				mainBlock += GetFieldIfCan(explanation.Example[0], "Example") + "\n"
			}
		}
	}

	return mainBlock + "\n"
}

func GetFieldIfCan(value string, field string) string {
	if len([]rune(value)) > 0 {
		if len([]rune(field)) > 0 {
			return fmt.Sprintf("%s", DecodeForTelegram(value))

		}
		return fmt.Sprintf("%s: %s", field, DecodeForTelegram(value))
	}
	return ""
}

func GetChangeTranslateMsg(translate string) string {
	return fmt.Sprintf("I changed translation:  %s", DecodeForTelegram(translate))
}

func GetTelegramRequest(chatId int, text string) SendMessageReqBody {
	return SendMessageReqBody{
		ChatID:      chatId,
		Text:        text,
		ParseMode:   "MarkdownV2",
		ReplyMarkup: ReplyMarkup{Keyboard: [][]Keyboard{}, OneTimeKeyboard: true, ResizeKeyboard: true},
	}
}

func DecodeForTelegram(text string) string {
	replacer := strings.NewReplacer(
		">", "\\>",
		"<", "\\<",
		".", "\\.",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
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
