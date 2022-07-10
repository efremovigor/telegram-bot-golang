package telegram

import (
	"encoding/json"
	"fmt"
	"strings"
	"telegram-bot-golang/helper"
)

const MaxRequestSize = 3000

type ReasonMessage int

const ReasonTypeNextMessage ReasonMessage = 1
const ReasonSubCambridgeMessage ReasonMessage = 2
const ReasonFullCambridgeMessage ReasonMessage = 3
const ReasonFullMultitranMessage ReasonMessage = 4

type Listener struct {
	Message chan RequestChannelTelegram
}

type Collector struct {
	Type     ReasonMessage
	Messages []RequestTelegramText
}

func (c *Collector) Add(messages ...RequestTelegramText) {
	for _, message := range messages {
		if len(c.Messages) > 0 {
			switch c.Type {
			case ReasonTypeNextMessage:
				c.Messages[len(c.Messages)-1].SetHasMore(NextMessage)
			case ReasonSubCambridgeMessage:
				c.Messages[len(c.Messages)-1].SetHasMore(NextMessageSubCambridge)
			case ReasonFullCambridgeMessage:
				c.Messages[len(c.Messages)-1].SetHasMore(NextMessageFullCambridge)
			case ReasonFullMultitranMessage:
				c.Messages[len(c.Messages)-1].SetHasMore(NextMessageFullMultitran)
			}
		}
		c.Messages = append(c.Messages, message)
	}
	fmt.Println(helper.ToJson(messages[0]))
}

func (c Collector) GetMessageForSave() (output []RequestChannelTelegram) {
	for _, message := range c.Messages {
		output = append(output, NewRequestChannelTelegram("text", message))
	}
	return
}

type RequestChannelTelegram struct {
	Type    string `json:"type"`
	Message []byte `json:"message"`
}

func NewRequestChannelTelegram(requestType string, request interface{}) RequestChannelTelegram {
	if requestInJson, err := json.Marshal(request); err == nil {
		return RequestChannelTelegram{Type: requestType, Message: requestInJson}
	}
	return RequestChannelTelegram{}
}

type RequestTelegramText struct {
	Word    string     `json:"word"`
	Text    string     `json:"text"`
	ChatId  int        `json:"chatId"`
	Buttons []Keyboard `json:"buttons"`
}

func MakeRequestTelegramText(word string, text string, chatId int, buttons []Keyboard) RequestTelegramText {
	return RequestTelegramText{
		Word:    word,
		Text:    text,
		ChatId:  chatId,
		Buttons: buttons,
	}
}

func (r *RequestTelegramText) SetHasMore(command string) {
	r.Buttons = append(r.Buttons, Keyboard{Text: "📚 more", CallbackData: command + " " + r.Word})
}

func MergeRequestTelegram(one RequestTelegramText, two RequestTelegramText) RequestTelegramText {
	one.Text += two.Text
	return one
}

type CallbackQuery struct {
	UpdateId      int `json:"update_id"`
	CallbackQuery struct {
		Id   string `json:"id"`
		From struct {
			Id           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Message struct {
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
			Date     int    `json:"date"`
			Text     string `json:"text"`
			Entities []struct {
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				Type   string `json:"type"`
				User   struct {
					Id           int    `json:"id"`
					IsBot        bool   `json:"is_bot"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					Username     string `json:"username"`
					LanguageCode string `json:"language_code"`
				} `json:"user"`
			} `json:"entities"`
			ReplyMarkup struct {
				InlineKeyboard [][]struct {
					Text         string `json:"text"`
					CallbackData string `json:"callback_data"`
				} `json:"inline_keyboard"`
			} `json:"reply_markup"`
		} `json:"message"`
		ChatInstance string `json:"chat_instance"`
		Data         string `json:"data"`
	} `json:"callback_query"`
}

type WebhookMessage struct {
	Message struct {
		Text      string `json:"text"`
		MessageId int    `json:"message_id"`
		Chat      struct {
			Id        int    `json:"id"`
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

type FileResponse struct {
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
		Photo []struct {
			FileId       string `json:"file_id"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
			Width        int    `json:"width"`
			Height       int    `json:"height"`
		} `json:"photo"`
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

type IncomingTelegramQueryInterface interface {
	IsValid() bool
	GetChatId() int
	GetChatText() string
	GetUsername() string
	GetUserId() int
	SetChatText(value string)
}

func (body CallbackQuery) IsValid() bool {
	if body.UpdateId != 0 {
		return true
	} else {
		return false
	}
}

func (body CallbackQuery) GetChatId() int {
	return body.CallbackQuery.Message.Chat.Id
}
func (body CallbackQuery) GetChatText() string {
	return body.CallbackQuery.Data
}

func (body CallbackQuery) GetUsername() string {
	return body.CallbackQuery.From.Username
}
func (body CallbackQuery) GetUserId() int {
	return body.CallbackQuery.From.Id
}

func (body *CallbackQuery) SetChatText(value string) {
	body.CallbackQuery.Data = value
}

func (body WebhookMessage) IsValid() bool {
	if body.Message.Chat.Id != 0 {
		return true
	} else {
		return false
	}
}

func (body WebhookMessage) GetChatId() int {
	if body.Message.Chat.Id != 0 {
		return body.Message.Chat.Id
	} else {
		return body.EditedMessage.Chat.ID
	}
}
func (body WebhookMessage) GetChatText() string {
	if body.Message.Chat.Id != 0 {
		return strings.ToLower(strings.TrimSpace(body.Message.Text))
	} else {
		return strings.ToLower(strings.TrimSpace(body.EditedMessage.Text))
	}
}

func (body *WebhookMessage) SetChatText(value string) {
	if body.Message.Chat.Id != 0 {
		body.Message.Text = value
	} else {
		body.EditedMessage.Text = value
	}
}

func (body WebhookMessage) GetUsername() string {
	if body.Message.Chat.Id != 0 {
		return body.Message.From.FirstName
	} else {
		return body.EditedMessage.From.FirstName
	}
}
func (body WebhookMessage) GetUserId() int {
	if body.Message.Chat.Id != 0 {
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
	Keyboard [][]Keyboard `json:"inline_keyboard"`
}

func (r *ReplyMarkup) SetKeyboard(buttons []Keyboard) {
	r.Keyboard = append(r.Keyboard, buttons)
}

type Keyboard struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type SendEarlierVoiceRequest struct {
	Performer           string      `json:"performer"`
	Title               string      `json:"title"`
	ChatId              int         `json:"chat_id"`
	Audio               string      `json:"audio"`
	Duration            interface{} `json:"duration"`
	DisableNotification bool        `json:"disable_notification"`
	ReplyToMessageId    interface{} `json:"reply_to_message_id"`
	ReplyMarkup         ReplyMarkup `json:"reply_markup"`
}

type SendEarlierPhotoRequest struct {
	ChatId      int         `json:"chat_id"`
	Photo       string      `json:"photo"`
	ReplyMarkup ReplyMarkup `json:"reply_markup"`
}

type UserRequest struct {
	Request string                   `json:"request"`
	Output  []RequestChannelTelegram `json:"output"`
}
