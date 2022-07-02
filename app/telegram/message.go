package telegram

import (
	"fmt"
	"strings"
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/helper"
	"telegram-bot-golang/service/dictionary/cambridge"
	"telegram-bot-golang/service/dictionary/multitran"
)

func GetBaseMsg(name string, id int) string {
	return fmt.Sprintf("Hey, [%s](tg://user?id=%d)\n", name, id) + "\n"
}
func GetIGotYourNewRequest(base string) string {
	return fmt.Sprintf(
		DecodeForTelegram("I got your message: ")+"*%s*\n", DecodeForTelegram(base))
}

func GetBlockWithRapidInfo(word string, translate string) string {
	return fmt.Sprintf(
		DecodeForTelegram("✅ *Rapid-microsoft*: ")+"*%s*\n", DecodeForTelegram(word)) + "\n" +
		GetFieldIfCan(translate, "💡 Translate")
}

func GetCambridgeHeaderBlock(word string) string {
	return fmt.Sprintf("✅ *Cambridge\\-dictionary*\\: *%s*", DecodeForTelegram(word)) + "\n\n"
}

func GetMultitranHeaderBlock(word string) string {
	return fmt.Sprintf("✅ *Multitran\\-dictionary*\\: *%s*", DecodeForTelegram(word)) + "\n\n"
}

func GetCambridgeOptionBlock(chatId int, info cambridge.Info) []RequestTelegramText {
	var messages []RequestTelegramText
	var mainBlock string
	mainBlock += fmt.Sprintf("*Word*\\: *%s*", DecodeForTelegram(info.Text))
	mainBlock += fmt.Sprintf("\\(%s\\)", DecodeForTelegram(info.Type)) + "\n"
	for lang, transcription := range info.Transcription {
		mainBlock += fmt.Sprintf("*%s*:\\[%s\\] ", strings.ToUpper(lang), DecodeForTelegram(transcription))
	}
	if len(info.Transcription) > 0 {
		mainBlock += "\n"
	}
	for n, explanation := range info.Explanation {
		if helper.Len(mainBlock) > MaxRequestSize {
			messages = append(messages,
				MakeRequestTelegramText(
					info.Text,
					mainBlock+"\n",
					chatId,
				),
			)
			mainBlock = ""
		}
		if n > 0 {
			mainBlock += GetRowSeparation()
		}
		mainBlock += GetFieldIfCan(explanation.Text, "Phrase")
		mainBlock += GetFieldIfCan(explanation.Level, "Level")
		mainBlock += GetFieldIfCan(explanation.SemanticDescription, "Semantic")
		mainBlock += GetFieldIfCan(explanation.Description, "Description")
		mainBlock += GetFieldIfCan(explanation.Translate, "💡 Translate")
		if len(explanation.Example) > 0 {
			mainBlock += "*Example*:\n"
		}
		for _, example := range explanation.Example {
			if helper.Len(mainBlock) > MaxRequestSize {
				messages = append(messages,
					MakeRequestTelegramText(
						info.Text,
						mainBlock+"\n",
						chatId,
					),
				)
				mainBlock = ""
			}
			mainBlock += DecodeForTelegram(example) + "\n"
		}
	}
	messages = append(messages, MakeRequestTelegramText(info.Text, mainBlock+"\n", chatId))
	return messages
}

func GetMultitranOptionBlock(chatId int, page multitran.Page) []RequestTelegramText {
	var messages []RequestTelegramText
	var mainBlock string
	for _, info := range page.Options {
		mainBlock += fmt.Sprintf("*Word*\\: *%s* \\[%s\\] \\(%s\\)", DecodeForTelegram(info.Text), DecodeForTelegram(info.Transcription), DecodeForTelegram(info.Type)) + "\n"
		for _, explanation := range info.Explanation {
			if helper.Len(mainBlock) > MaxRequestSize {
				messages = append(messages,
					MakeRequestTelegramText(
						info.Text,
						mainBlock+"\n",
						chatId,
					),
				)
				mainBlock = ""
			}
			mainBlock += GetFieldIfCan(explanation.Type, "Type")
			mainBlock += "*Explanation*:\n"

			for i, translate := range explanation.Text {
				if helper.Len(mainBlock) > MaxRequestSize {
					messages = append(messages,
						MakeRequestTelegramText(
							info.Text,
							mainBlock+"\n",
							chatId,
						),
					)
					mainBlock = ""
				}
				mainBlock += DecodeForTelegram(translate)

				if i < len(explanation.Text) {
					mainBlock += ", "
				}
			}
			mainBlock += "\n"
		}
	}
	messages = append(messages, MakeRequestTelegramText(
		page.RequestText,
		mainBlock+"\n",
		chatId,
	))
	return messages
}

func GetFieldIfCan(value string, field string) string {
	if len([]rune(value)) > 0 {
		return fmt.Sprintf("*%s*: %s", field, DecodeForTelegram(value)) + "\n"
	}
	return ""
}

func GetChangeTranslateMsg(translate string) string {
	return fmt.Sprintf("I changed translation:  %s", DecodeForTelegram(translate))
}

func GetRatingHeader(n int, all bool) string {
	if all {
		return fmt.Sprintf("*Top %d words used:*\n", n)
	}
	return fmt.Sprintf("*My %d words used:*\n", n)
}
func GetRowRating(n int, statistic model.WordStatistic) string {
	return fmt.Sprintf("*%d*\\. %s \\- %d\n", n, statistic.Word, statistic.Count)
}

func GetRowSeparation() string {
	return DecodeForTelegram("▫️◾️▫️◾️▫️◾️▫️") + "\n"
}
