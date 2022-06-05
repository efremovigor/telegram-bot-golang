package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func GetHelloIGotYourMSGRequest(body WebhookMessage) {
	SendMessage(GetTelegramRequest(
		body.GetChatId(),
		GetBaseMsg(body.GetUsername(), body.GetUserId())+
			GetIGotYourNewRequest(body.GetChatText()),
	))

}

func GetResultFromRapidMicrosoft(body WebhookMessage, state string) {
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

func GetResultFromCambridge(body WebhookMessage) {
	cambridgeInfo := cambridge.Get(body.Message.Text)
	if cambridgeInfo.IsValid() {
		statistic.Consider(cambridgeInfo.Text, body.GetUserId())
	}
	SendMessage(GetTelegramRequest(
		body.GetChatId(), GetBlockWithCambridge(cambridgeInfo)))
	SendVoices(body.GetChatId(), cambridgeInfo)
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

func SendVoices(chatId int, info cambridge.Info) {

	if len([]rune(info.VoicePath.UK)) > 0 {
		sendVoice(chatId, "UK", info)
	}
	//if len([]rune(info.VoicePath.US)) > 0 {
	//	sendVoice(chatId, "US", info)
	//}
}

func sendVoice(chatId int, country string, info cambridge.Info) {
	var path string
	switch country {
	case "UK":
		path = info.VoicePath.UK
	case "US":
		path = info.VoicePath.US
	}
	resp, err := http.Get(cambridge.Url + path)

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", filepath.Base("audio.mp3"))
	io.Copy(part, resp.Body)

	_ = writer.WriteField("performer", country)
	_ = writer.WriteField("title", info.Text)
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

	buf, _ := ioutil.ReadAll(res.Body)
	b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		log.Fatalln(err)
	}

	var audioResponse AudioResponse
	if err = json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(b))).Decode(&audioResponse); err != nil && !audioResponse.Ok {
		fmt.Println("could not decode telegram response", err)
	} else {
		//redis.Set(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.Text, country), audioResponse.Result.Document.FileId)
		fmt.Println(audioResponse)

		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("audio", audioResponse.Result.Document.FileId)
		io.Copy(part, resp.Body)

		_ = writer.WriteField("performer", country)
		_ = writer.WriteField("title", info.Text)
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

		buf, _ := ioutil.ReadAll(res.Body)
		b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(b)

	}
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
