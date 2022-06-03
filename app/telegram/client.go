package telegram

import (
	"strings"
	"telegram-bot-golang/service/dictionary/cambridge"
	rapid_microsoft "telegram-bot-golang/service/translate/rapid-microsoft"
	"telegram-bot-golang/statistic"
)

func Reply(body WebhookReqBody, state string) SendMessageReqBody {
	var from, to string
	if state == "" || state == "en_ru" {
		from = "en"
		to = "ru"
	} else {
		from = "ru"
		to = "en"
	}
	cambridgeInfo := cambridge.Get(body.Message.Text)
	if cambridgeInfo.IsValid() {
		statistic.Consider(cambridgeInfo.Text, body.GetUserId())
	}

	return GetTelegramRequest(
		body.GetChatId(),
		GetBaseMsg(body.GetUsername(), body.GetUserId())+
			GetIGotYourNewRequest(body.GetChatText())+
			GetBlockWithRapidInfo(rapid_microsoft.GetTranslate(body.Message.Text, to, from))+
			GetBlockWithCambridge(cambridgeInfo),
	)
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
		"+", "\\+",
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
