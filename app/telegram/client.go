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
	"telegram-bot-golang/service/dictionary/multitran"
	rapid_microsoft "telegram-bot-golang/service/translate/rapid-microsoft"
	"telegram-bot-golang/statistic"
	telegramConfig "telegram-bot-golang/telegram/config"
)

func GetHelloIGotYourMSGRequest(body WebhookMessage) SendMessageReqBody {
	return GetTelegramRequest(
		body.GetChatId(),
		GetBaseMsg(body.GetUsername(), body.GetUserId())+
			GetIGotYourNewRequest(body.GetChatText()),
	)
}

func GetResultFromRapidMicrosoft(body WebhookMessage, state string) SendMessageReqBody {
	var from, to string
	if state == "" || state == "en_ru" {
		from = "en"
		to = "ru"
	} else {
		from = "ru"
		to = "en"
	}

	translate := rapid_microsoft.GetTranslate(body.GetChatText(), to, from)
	if helper.IsEmpty(translate) {
		return SendMessageReqBody{}
	}
	return GetTelegramRequest(body.GetChatId(), GetBlockWithRapidInfo(translate))
}

func GetResultFromCambridge(cambridgeInfo cambridge.CambridgeInfo, body WebhookMessage) []SendMessageReqBody {
	statistic.Consider(cambridgeInfo.RequestText, body.GetUserId())
	var messages []SendMessageReqBody
	messages = append(messages, GetTelegramRequest(body.GetChatId(), GetCambridgeHeaderBlock(cambridgeInfo)))
	for _, option := range cambridgeInfo.Options {
		messages = append(messages, GetCambridgeOptionBlock(body.GetChatId(), option)...)
	}
	return messages
}

func GetResultFromMultitran(info multitran.Page, body WebhookMessage) []SendMessageReqBody {
	statistic.Consider(info.RequestText, body.GetUserId())
	var messages []SendMessageReqBody
	messages = append(messages, GetTelegramRequest(body.GetChatId(), GetMultitranHeaderBlock(info)))
	messages = append(messages, GetMultitranOptionBlock(body.GetChatId(), info)...)
	return messages
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
	return strings.NewReplacer(
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
	).Replace(text)
}

func sendMessage(response SendMessageReqBody) {
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

func sendVoices(chatId int, info cambridge.CambridgeInfo) {

	if voiceId, err := redis.Get(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.RequestText, "uk")); err == nil && len([]rune(voiceId)) > 0 {
		fmt.Println("find key uk voice in cache")
		sendVoiceFromCache(chatId, "uk", voiceId, info)
	} else {
		sendVoice(chatId, "uk", info)
	}

	if voiceId, err := redis.Get(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.RequestText, "us")); err == nil && len([]rune(voiceId)) > 0 {
		fmt.Println("find key us voice in cache")
		sendVoiceFromCache(chatId, "us", voiceId, info)
	} else {
		sendVoice(chatId, "us", info)
	}
}

func sendVoice(chatId int, country string, info cambridge.CambridgeInfo) {
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
	defer helper.CloseConnection(res.Body)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", filepath.Base("audio.mp3"))
	_, err = io.Copy(part, res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = writer.WriteField("performer", country)

	var title string
	if !helper.IsEmpty(info.Options[0].Text) {
		title = info.Options[0].Text
	} else {
		title = info.RequestText
	}
	_ = writer.WriteField("title", title)
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
		redis.Set(fmt.Sprintf(redis.WordVoiceTelegramKeys, info.RequestText, country), audioResponse.Result.Audio.FileId)
	}
}

func sendVoiceFromCache(chatId int, country string, audioId string, info cambridge.CambridgeInfo) {
	var title string
	if !helper.IsEmpty(info.Options[0].Text) {
		title = info.Options[0].Text
	} else {
		title = info.RequestText
	}
	request := SendEarlierVoiceRequest{Performer: country, Title: title, Audio: audioId, ChatId: chatId}
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
