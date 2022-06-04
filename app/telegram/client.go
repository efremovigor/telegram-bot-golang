package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"telegram-bot-golang/env"
	"telegram-bot-golang/service/dictionary/cambridge"
	rapid_microsoft "telegram-bot-golang/service/translate/rapid-microsoft"
	"telegram-bot-golang/statistic"
	telegramConfig "telegram-bot-golang/telegram/config"
)

func GetHelloIGotYourMSGRequest(body WebhookReqBody) {
	SendMessage(GetTelegramRequest(
		body.GetChatId(),
		GetBaseMsg(body.GetUsername(), body.GetUserId())+
			GetIGotYourNewRequest(body.GetChatText()),
	))

}

func GetResultFromRapidMicrosoft(body WebhookReqBody, state string) {
	var from, to string
	if state == "" || state == "en_ru" {
		from = "en"
		to = "ru"
	} else {
		from = "ru"
		to = "en"
	}
	SendMessage(
		GetTelegramRequest(
			body.GetChatId(),
			GetBlockWithRapidInfo(
				rapid_microsoft.GetTranslate(body.Message.Text, to, from),
			),
		))
}

func GetResultFromCambridge(body WebhookReqBody) {
	cambridgeInfo := cambridge.Get(body.Message.Text)
	if cambridgeInfo.IsValid() {
		statistic.Consider(cambridgeInfo.Text, body.GetUserId())
	}
	SendMessage(GetTelegramRequest(
		body.GetChatId(), GetBlockWithCambridge(cambridgeInfo)))
	SendVoices(body.GetChatId(), cambridgeInfo.VoicePath)
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

func SendMessage(response SendMessageReqBody) {
	if len([]rune(response.Text)) > 0 {
		toTelegram, err := json.Marshal(response)
		if err != nil {
			fmt.Println("error of serialisation telegram struct:" + string(toTelegram))
		}
		fmt.Println("----")
		fmt.Println("to telegram json:" + string(toTelegram))
		fmt.Println("+++")
		fmt.Println("+++")

		res, err := http.Post(telegramConfig.GetTelegramUrl(), "application/json", bytes.NewBuffer(toTelegram))
		if err != nil {
			fmt.Println("error of sending message to telegram:" + string(toTelegram))
		}
		if res.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println("bad response from telegram:" + res.Status + " Message:" + string(body) + "\n" + string(toTelegram))
		}
	}
}

func SendVoices(chatId int, voices cambridge.VoicePath) {

	if len([]rune(voices.UK)) > 0 {
		sendVoice(chatId, voices.UK)
	}
}

func sendVoice(chatId int, path string) {
	resp, err := http.Get(cambridge.Url + path)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", filepath.Base("audio.mp3"))
	io.Copy(part, resp.Body)

	_ = writer.WriteField("performer", "hello")
	_ = writer.WriteField("title", "hello")
	_ = writer.WriteField("chat_id", strconv.Itoa(chatId))
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	r, _ := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendAudio", env.GetEnvVariable("TELEGRAM_API_TOKEN")), body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body1, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body1))
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
