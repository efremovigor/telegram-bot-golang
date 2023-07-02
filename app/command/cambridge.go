package command

import (
	"encoding/json"
	"fmt"
	"telegram-bot-golang/db/redis"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/statistic"
	"telegram-bot-golang/telegram"
)

func MakeCambridgeFullIfEmpty(chatId int, userId int, chatText string) {
	collector := telegram.Collector{Type: telegram.ReasonFullCambridgeMessage}
	if page := cambridge.Get(chatText); page.IsValid() {
		statistic.Consider(chatText, userId)
		collector.Add(
			handleCambridgePage(page, chatId, chatText)...,
		)
	}

	saveMessagesQueue(fmt.Sprintf(redis.NextFullInfoRequestMessageKey, "cambridge", userId, chatText), chatText, collector.GetMessageForSave())
}

func handleCambridgePage(page cambridge.Page, chatId int, chatText string) (messages []telegram.RequestTelegramText) {
	messages = telegram.GetResultFromCambridge(page, chatId, chatText)
	return
}

func GetSubCambridge(chatId int, userId int, chatText string) {
	cambridgeFounded, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeSearchValue, chatText))
	if err != nil {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Don't find cambridge key - word:%s", chatText))
		return
	}
	if page := cambridge.DoRequest(chatText, cambridge.Url+cambridgeFounded, ""); page.IsValid() {
		statistic.Consider(chatText, userId)
		collector := telegram.Collector{Type: telegram.ReasonSubCambridgeMessage}
		collector.Add(handleCambridgePage(page, chatId, chatText)...)
		saveMessagesQueue(
			fmt.Sprintf(redis.SubCambridgeMessageKey, userId, chatText),
			chatText,
			collector.GetMessageForSave(),
		)

		return
	} else {
		fmt.Println(fmt.Sprintf("[Strange behaivor] Aren't able to parse url:%s", cambridge.Url+cambridgeFounded))
	}
}

func SendVoice(query telegram.IncomingTelegramQueryInterface, lang string, hash string) {
	url, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeUniqVoiceLink, hash))
	if err != nil {
		fmt.Println(err)
		return
	}
	var cache redis.VoiceFile
	if err := json.Unmarshal([]byte(url), &cache); err != nil {
		fmt.Println(err)
		return
	}
	telegram.SendVoices(query.GetChatId(), lang, cache)
}

func SendImage(query telegram.IncomingTelegramQueryInterface, hash string) {
	url, err := redis.Get(fmt.Sprintf(redis.InfoCambridgeUniqPicLink, hash))
	if err != nil {
		fmt.Println(err)
		return
	}
	var cache redis.PicFile
	if err := json.Unmarshal([]byte(url), &cache); err != nil {
		fmt.Println(err)
		return
	}
	telegram.SendImage(query.GetChatId(), cache)
}
