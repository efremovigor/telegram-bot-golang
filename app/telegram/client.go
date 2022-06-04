package telegram

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func SendVoice(chatId int, word string) {

	out, err := os.Create("./cache/media/cambridge/hello.mp3")
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()
	resp, err := http.Get("https://dictionary.cambridge.org/media/english-russian/uk_pron/u/ukh/ukhef/ukheft_029.mp3")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendAudio", env.GetEnvVariable("TELEGRAM_API_TOKEN"))

	payload := strings.NewReader(fmt.Sprintf("{\"performer\":\"Hello\",\"title\":\"Hello\",\"chat_id\":%d,\"audio\":\"./cache/media/cambridge/hello.mp3\",\"duration\":null,\"disable_notification\":false,\"reply_to_message_id\":null}", chatId))

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

type T struct {
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
		Date     int `json:"date"`
		Document struct {
			FileName     string `json:"file_name"`
			MimeType     string `json:"mime_type"`
			FileId       string `json:"file_id"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
		} `json:"document"`
	} `json:"result"`
}

//{"ok":true,"result":{"message_id":841,"from":{"id":5125700707,"is_bot":true,"first_name":"EnglishHelper","username":"IdontSpeakBot"},"chat":{"id":184357122,"first_name":"Igor","last_name":"Efremov","username":"Igor198811","type":"private"},"date":1654357898,"document":{"file_name":"ukheft_029.ogg","mime_type":"audio/ogg","file_id":"BQACAgQAAxkDAAIDSWKbf4rOWkrezgXn9ZZSvqqWNF7NAAIGAwACJpTkUF3cWGDxH4YgJAQ","file_unique_id":"AgADBgMAAiaU5FA","file_size":8769}}}
type T2 struct {
	Title               string      `json:"title"`
	ChatId              int         `json:"chat_id"`
	Audio               string      `json:"audio"`
	Duration            interface{} `json:"duration"`
	DisableNotification bool        `json:"disable_notification"`
	ReplyToMessageId    interface{} `json:"reply_to_message_id"`
}
