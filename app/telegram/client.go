package telegram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"telegram-bot-golang/env"
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
		"_", "\\_",
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

func SendVoice(chatId int) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendVoice", env.GetEnvVariable("TELEGRAM_API_TOKEN"))

	payload := strings.NewReader(fmt.Sprintf("{\"chat_id\":%d,\"voice\":\"BQACAgQAAxkDAAIDSWKbf4rOWkrezgXn9ZZSvqqWNF7NAAIGAwACJpTkUF3cWGDxH4YgJAQ\",\"duration\":null,\"disable_notification\":false,\"reply_to_message_id\":null}", chatId))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json")

	req.Header.Add("User-Agent", "Telegram Bot SDK - (https://github.com/irazasyed/telegram-bot-sdk)")

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)

	fmt.Println(string(body))

}
