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
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/helper"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
	rapidMicrosoft "telegram-bot-golang/service/translate/rapid-microsoft"
	telegramConfig "telegram-bot-golang/telegram/config"
)

const NextMessage = "/next_message"
const NextMessageSubCambridge = "/next_message_cambridge"
const NextMessageFullCambridge = "/next_message_full_cambridge"
const NextMessageFullMultitran = "/next_message_full_multitran"
const ShowRequestVoice = "/show_voice"
const ShowRequestPic = "/show_pic"
const SearchRequest = "/search"
const ShowFull = "/show_full"

const LangEn = "en"
const LangRu = "ru"
const CountryUk = "uk"
const CountryUs = "us"

func GetHelloIGotYourMSGRequest(query IncomingTelegramQueryInterface) RequestTelegramText {
	return MakeRequestTelegramText(
		query.GetChatText(),
		GetBaseMsg(query.GetUsername(), query.GetUserId())+
			GetIGotYourNewRequest(query.GetChatText()),
		query.GetChatId(),
		[]Keyboard{},
	)
}

func GetResultFromRapidMicrosoft(chatId int, chatText string, state string) RequestTelegramText {
	var from, to string

	if state == "" {
		if helper.IsEn(chatText) {
			from = LangEn
			to = LangRu
		} else {
			from = LangRu
			to = LangEn
		}
	} else if state == LangEn+"_"+LangRu {
		from = LangEn
		to = LangRu
	} else {
		from = LangRu
		to = LangEn
	}

	translate := rapidMicrosoft.GetTranslate(chatText, to, from)
	if helper.IsEmpty(translate) {
		return RequestTelegramText{}
	}
	return MakeRequestTelegramText(
		chatText,
		GetBlockWithRapidInfo(chatText, translate),
		chatId,
		[]Keyboard{},
	)
}

func GetResultFromCambridge(cambridgeInfo cambridge.Page, chatId int, chatText string) []RequestTelegramText {
	var messages []RequestTelegramText

	for _, option := range cambridgeInfo.Options {
		requests := GetCambridgeOptionBlock(chatId, option)
		if len(requests) > 0 {
			if len(messages) == 0 {
				messages = append(
					messages,
					MergeRequestTelegram(
						MakeRequestTelegramText(chatText, GetCambridgeHeaderBlock(cambridgeInfo.RequestText), chatId, []Keyboard{}),
						requests[0],
					),
				)
				if len(requests) > 1 {
					messages = append(messages, requests[1:]...)
				}
			} else {
				messages = append(messages, requests...)
			}
		}

	}
	return messages
}

func GetResultFromMultitran(info multitran.Page, chatId int, chatText string) []RequestTelegramText {
	var messages []RequestTelegramText
	requests := GetMultitranOptionBlock(chatId, info)
	if len(requests) > 0 {
		messages = append(
			messages,
			MergeRequestTelegram(
				MakeRequestTelegramText(
					chatText,
					GetMultitranHeaderBlock(info.RequestText),
					chatId,
					[]Keyboard{},
				),
				requests[0],
			),
		)
	}
	if len(requests) > 1 {
		messages = append(messages, requests[1:]...)
	}
	return messages
}

func GetTelegramRequest(chatId int, text string, buttons []Keyboard) SendMessageReqBody {
	var keyboard [][]Keyboard
	fmt.Println(fmt.Sprintf("count of buttons: %d", len(buttons)))
	if len(buttons) > 0 {
		var buffer []Keyboard
		var bufferCount int
		for i, button := range buttons {

			if bufferCount+helper.Len(button.Text) > 40 {
				keyboard = append(keyboard, buffer)
				buffer = []Keyboard{}
				bufferCount = 0
			}

			bufferCount += helper.Len(button.Text)
			buffer = append(buffer, button)

			if i == len(buttons)-1 {
				keyboard = append(keyboard, buffer)
			}

		}
	} else {
		keyboard = [][]Keyboard{}
	}

	return SendMessageReqBody{
		ChatID:      chatId,
		Text:        text,
		ParseMode:   "MarkdownV2",
		ReplyMarkup: ReplyMarkup{Keyboard: keyboard},
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

func sendBaseInfo(telegramText RequestTelegramText) {
	request := GetTelegramRequest(telegramText.ChatId, telegramText.Text, telegramText.Buttons)
	sendBaseMessage(request)
}

func sendBaseMessage(request SendMessageReqBody) {
	if len([]rune(request.Text)) > 0 {
		toTelegram, err := json.Marshal(request)
		if err != nil {
			fmt.Println("error of serialisation telegram struct:" + string(toTelegram))
		}

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

func SendImage(chatId int, cache redis.PicFile) {

	if picId, err := redis.Get(fmt.Sprintf(redis.WordPicTelegramKey, cache.Word, 0)); err == nil && len([]rune(picId)) > 0 {
		sendPicFromCache(chatId, picId)
	} else {
		sendPic(chatId, cache)
	}
}

func SendVoices(chatId int, lang string, cache redis.VoiceFile) {

	if voiceId, err := redis.Get(fmt.Sprintf(redis.WordVoiceTelegramKey, cache.Word, lang)); err == nil && len([]rune(voiceId)) > 0 {
		sendVoiceFromCache(chatId, cache, voiceId)
	} else {
		sendVoice(chatId, cache)
	}
}

func sendVoice(chatId int, cache redis.VoiceFile) {
	res, err := http.Get(cambridge.Url + cache.Url)

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

	_ = writer.WriteField("performer", cache.Lang)
	_ = writer.WriteField("title", cache.Word)
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
	qwe, err := r.GetBody()
	body1, _ := ioutil.ReadAll(qwe)
	fmt.Println(string(body1))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer rapidMicrosoft.CloseConnection(res.Body)

	var fileResponse FileResponse
	if _, err := helper.ParseJson(res.Body, &fileResponse); err != nil && !fileResponse.Ok {
		fmt.Println("could not decode telegram response", err)
	} else {
		redis.Set(fmt.Sprintf(redis.WordVoiceTelegramKey, cache.Word, cache.Lang), fileResponse.Result.Audio.FileId, 0)
	}
}

func sendPic(chatId int, cache redis.PicFile) {

	res, err := http.Get(cambridge.Url + cache.Url)

	if err != nil {
		fmt.Println(err)
	}
	defer helper.CloseConnection(res.Body)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("photo", filepath.Base("pic.jpg"))
	_, err = io.Copy(part, res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = writer.WriteField("chat_id", strconv.Itoa(chatId))

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	r, _ := http.NewRequest("POST", telegramConfig.GetTelegramUrl("sendPhoto"), body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err = client.Do(r)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer rapidMicrosoft.CloseConnection(res.Body)

	var fileResponse FileResponse
	if _, err := helper.ParseJson(res.Body, &fileResponse); err != nil && !fileResponse.Ok {
		fmt.Println("could not decode telegram response", err)
	} else {
		redis.Set(fmt.Sprintf(redis.WordPicTelegramKey, cache.Word, 0), fileResponse.Result.Photo[0].FileId, 0)
	}
}

func sendVoiceFromCache(chatId int, cache redis.VoiceFile, audioId string) {

	request := SendEarlierVoiceRequest{Performer: cache.Lang, Title: cache.Word, Audio: audioId, ChatId: chatId, ReplyMarkup: ReplyMarkup{Keyboard: [][]Keyboard{}}}
	requestInJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}

	doJsonRequest(telegramConfig.GetTelegramUrl("sendAudio"), requestInJson)
}

func sendPicFromCache(chatId int, photoId string) {
	request := SendEarlierPhotoRequest{Photo: photoId, ChatId: chatId, ReplyMarkup: ReplyMarkup{Keyboard: [][]Keyboard{}}}
	requestInJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	doJsonRequest(telegramConfig.GetTelegramUrl("sendPhoto"), requestInJson)
}

func doJsonRequest(url string, json []byte) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(json)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println("bad response from telegram:" + res.Status + " Message:" + string(body) + "\n")
	}
	defer rapidMicrosoft.CloseConnection(res.Body)
}
