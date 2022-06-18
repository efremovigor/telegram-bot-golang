package telegram

import (
	"strings"
	"telegram-bot-golang/service/dictionary/cambridge"
)

const MaxRequestSize = 3000

type Listener struct {
	Message chan RequestChannelTelegram
}

type RequestChannelTelegram struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

type CambridgeRequestTelegramVoice struct {
	Info   cambridge.CambridgeInfo `json:"info"`
	ChatId int                     `json:"chatId"`
}

type RequestTelegramText struct {
	Text   string `json:"text"`
	ChatId int    `json:"chatId"`
}

type WebhookMessage struct {
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

type AudioResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageId int `json:"message_id"`
		From      struct {
			Id        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date  int `json:"date"`
		Audio struct {
			Duration     int    `json:"duration"`
			FileName     string `json:"file_name"`
			MimeType     string `json:"mime_type"`
			Title        string `json:"title"`
			Performer    string `json:"performer"`
			FileId       string `json:"file_id"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
		} `json:"audio"`
	} `json:"result"`
}

func (body WebhookMessage) GetChatId() int {
	if body.Message.Chat.ID != 0 {
		return body.Message.Chat.ID
	} else {
		return body.EditedMessage.Chat.ID
	}
}
func (body WebhookMessage) GetChatText() string {
	if body.Message.Chat.ID != 0 {
		return strings.ToLower(strings.TrimSpace(body.Message.Text))
	} else {
		return strings.ToLower(strings.TrimSpace(body.EditedMessage.Text))
	}
}

func (body WebhookMessage) GetUsername() string {
	if body.Message.Chat.ID != 0 {
		return body.Message.From.FirstName
	} else {
		return body.EditedMessage.From.FirstName
	}
}
func (body WebhookMessage) GetUserId() int {
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

type SendEarlierVoiceRequest struct {
	Performer           string      `json:"performer"`
	Title               string      `json:"title"`
	ChatId              int         `json:"chat_id"`
	Audio               string      `json:"audio"`
	Duration            interface{} `json:"duration"`
	DisableNotification bool        `json:"disable_notification"`
	ReplyToMessageId    interface{} `json:"reply_to_message_id"`
}

type UserRequest struct {
	Request string `json:"request"`
	Output  []byte `json:"output"`
}
