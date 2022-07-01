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
	rapidMicrosoft "telegram-bot-golang/service/translate/rapid-microsoft"
	telegramConfig "telegram-bot-golang/telegram/config"
)

const NextRequestMessage = "/next_message"
const ShowRequestVoice = "/show_voice"
const SearchRequest = "/search"
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
	)
}

func GetResultFromRapidMicrosoft(query IncomingTelegramQueryInterface, state string) RequestTelegramText {
	var from, to string

	if state == "" {
		if helper.IsEn(query.GetChatText()) {
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

	translate := rapidMicrosoft.GetTranslate(query.GetChatText(), to, from)
	if helper.IsEmpty(translate) {
		return RequestTelegramText{}
	}
	return MakeRequestTelegramText(
		query.GetChatText(),
		GetBlockWithRapidInfo(translate),
		query.GetChatId(),
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
						MakeRequestTelegramText(chatText, GetCambridgeHeaderBlock(cambridgeInfo), chatId),
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

func GetResultFromMultitran(info multitran.Page, query IncomingTelegramQueryInterface) []RequestTelegramText {
	var messages []RequestTelegramText
	requests := GetMultitranOptionBlock(query.GetChatId(), info)
	if len(requests) > 0 {
		messages = append(
			messages,
			MergeRequestTelegram(
				MakeRequestTelegramText(
					query.GetChatText(),
					GetMultitranHeaderBlock(info),
					query.GetChatId(),
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
	if len(buttons) > 0 {
		var buffer []Keyboard
		for i, button := range buttons {

			if len(buffer) > 3 {
				keyboard = append(keyboard, buffer)
				buffer = []Keyboard{}
			}

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

func sendBaseInfo(telegramText RequestTelegramText, buttons []Keyboard) {
	request := GetTelegramRequest(telegramText.ChatId, telegramText.Text, buttons)
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

func SendVoices(chatId int, info cambridge.Page, lang string, hasMore bool) {

	if voiceId, err := redis.Get(fmt.Sprintf(redis.WordVoiceTelegramKey, info.RequestText, lang)); err == nil && len([]rune(voiceId)) > 0 {
		fmt.Println("find key " + lang + " voice in cache")
		sendVoiceFromCache(chatId, lang, voiceId, info, hasMore)
	} else {
		sendVoice(chatId, lang, info, hasMore)
	}
}

func sendVoice(chatId int, country string, info cambridge.Page, hasMore bool) {
	var path string
	switch country {
	case CountryUk:
		path = info.VoicePath.UK
	case CountryUs:
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

	if hasMore {
		_ = writer.WriteField("reply_markup", fmt.Sprintf("{\"inline_keyboard\":[[{\"text\":\"more\",\"callback_data\":\"%s %s\"}]]}", NextRequestMessage, info.RequestText))
	}

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

	buf, _ := ioutil.ReadAll(res.Body)
	b, err := io.ReadAll(ioutil.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		log.Fatalln(err)
	}

	var audioResponse AudioResponse
	if err = json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(b))).Decode(&audioResponse); err != nil && !audioResponse.Ok {
		fmt.Println("could not decode telegram response", err)
	} else {
		redis.Set(fmt.Sprintf(redis.WordVoiceTelegramKey, info.RequestText, country), audioResponse.Result.Audio.FileId, 0)
	}
}

func sendVoiceFromCache(chatId int, country string, audioId string, info cambridge.Page, hasMore bool) {
	var title string
	if !helper.IsEmpty(info.Options[0].Text) {
		title = info.Options[0].Text
	} else {
		title = info.RequestText
	}
	request := SendEarlierVoiceRequest{Performer: country, Title: title, Audio: audioId, ChatId: chatId, ReplyMarkup: ReplyMarkup{Keyboard: [][]Keyboard{}}}
	if hasMore {
		request.ReplyMarkup.SetKeyboard([]Keyboard{{Text: "more", CallbackData: NextRequestMessage + " " + title}})
	}
	requestInJson, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}

	req, _ := http.NewRequest("POST", telegramConfig.GetTelegramUrl("sendAudio"), strings.NewReader(string(requestInJson)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	qwe, err := req.GetBody()
	body1, _ := ioutil.ReadAll(qwe)
	fmt.Println(string(body1))

	if err != nil {
		fmt.Println(err)
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println("bad response from telegram:" + res.Status + " Message:" + string(body) + "\n")
	}
	defer rapidMicrosoft.CloseConnection(res.Body)
}
