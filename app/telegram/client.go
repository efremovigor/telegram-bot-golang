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
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/helper"
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

	translate := rapid_microsoft.GetTranslate(body.Message.Text, to, from)
	if helper.IsEmpty(translate) {
		return
	}
	SendMessage(GetTelegramRequest(body.GetChatId(), GetBlockWithRapidInfo(translate)))
}

func GetResultFromCambridge(body WebhookMessage) {
	cambridgeInfo := cambridge.Get(body.Message.Text)
	if !cambridgeInfo.IsValid() {
		return
	}
	statistic.Consider(cambridgeInfo.Text, body.GetUserId())
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

		res, err := http.Post(telegramConfig.GetTelegramUrl("sendMessage"), "application/json", bytes.NewBuffer(toTelegram))
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

	if voiceId, err := redis.Get(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.Text, "uk")); err == nil && len([]rune(voiceId)) > 0 {
		fmt.Println("find key uk voice in cache")
		sendVoiceFromCache(chatId, "uk", voiceId, info.Text)
	} else {
		sendVoice(chatId, "uk", info)
	}

	if voiceId, err := redis.Get(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.Text, "us")); err == nil && len([]rune(voiceId)) > 0 {
		fmt.Println("find key us voice in cache")
		sendVoiceFromCache(chatId, "us", voiceId, info.Text)
	} else {
		sendVoice(chatId, "us", info)
	}
}

func sendVoice(chatId int, country string, info cambridge.Info) {
	var path string
	switch country {
	case "uk":
		path = info.VoicePath.UK
	case "us":
		path = info.VoicePath.US
	}
	res, err := http.Get(cambridge.Url + path)

	if err != nil {
		fmt.Println(err)
	}
	defer rapid_microsoft.CloseConnection(res.Body)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", filepath.Base("audio.mp3"))
	_, err = io.Copy(part, res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = writer.WriteField("performer", country)
	_ = writer.WriteField("title", info.Text)
	_ = writer.WriteField("chat_id", strconv.Itoa(chatId))
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	r, _ := http.NewRequest("POST", telegramConfig.GetTelegramUrl("sendAudio"), body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err = client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rapid_microsoft.CloseConnection(res.Body)

	buf, _ := ioutil.ReadAll(res.Body)
	b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		log.Fatalln(err)
	}
	s, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))

	fmt.Println(string(s))

	var audioResponse AudioResponse
	if err = json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(b))).Decode(&audioResponse); err != nil && !audioResponse.Ok {
		fmt.Println("could not decode telegram response", err)
	} else {
		redis.Set(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.Text, country), audioResponse.Result.Audio.FileId)
	}
}

func sendVoiceFromCache(chatId int, country string, audioId string, word string) {
	request := SendEarlierVoiceRequest{Performer: country, Title: word, Audio: audioId, ChatId: chatId}
	requestInJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}

	req, _ := http.NewRequest("POST", telegramConfig.GetTelegramUrl("sendAudio"), strings.NewReader(string(requestInJson)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer rapid_microsoft.CloseConnection(res.Body)
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
