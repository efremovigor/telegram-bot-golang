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
	mainBlock += fmt.Sprintf("❗️*Word*\\: *%s*", DecodeForTelegram(info.Text))
	mainBlock += fmt.Sprintf("\\(%s\\)", DecodeForTelegram(info.Type)) + "\n\n"
	for lang, transcription := range info.Transcription {
		mainBlock += fmt.Sprintf("*%s*:\\[%s\\] ", strings.ToUpper(lang), DecodeForTelegram(transcription))
	}
	if len(info.Transcription) > 0 {
		mainBlock += "\n"
	}

	var explanationBlock string
	for n, explanation := range info.Explanation {
		if n > 0 {
			explanationBlock += "\n" + GetRowSeparation() + "\n"
		}
		explanationBlock += GetFieldIfCan(explanation.Text, "❗️Phrase") + "\n"
		explanationBlock += GetFieldIfCan(explanation.Level, "Level")
		explanationBlock += GetFieldIfCan(explanation.SemanticDescription, "📃 Semantic")
		explanationBlock += GetFieldIfCan(explanation.Description, "📃 Description")
		explanationBlock += GetFieldIfCan(explanation.Translate, "💡 Translate")
		var exampleBlock string
		if len(explanation.Example) > 0 {
			exampleBlock += "*Example*:\n"
		}
		for _, example := range explanation.Example {
			exampleBlock += "📌" + DecodeForTelegram(example) + "\n"
		}
		if helper.Len(mainBlock)+helper.Len(explanationBlock) > MaxRequestSize || helper.Len(exampleBlock) > MaxRequestSize {
			messages = append(messages,
				MakeRequestTelegramText(
					info.Text,
					mainBlock+explanationBlock+"\n",
					chatId,
				),
				MakeRequestTelegramText(
					info.Text,
					exampleBlock+"\n",
					chatId,
				),
			)
			mainBlock = ""
			mainBlock = ""
			explanationBlock = ""
		} else if helper.Len(mainBlock)+helper.Len(explanationBlock)+helper.Len(exampleBlock) > MaxRequestSize {
			messages = append(messages,
				MakeRequestTelegramText(
					info.Text,
					mainBlock+explanationBlock+exampleBlock+"\n",
					chatId,
				),
			)
			mainBlock = ""
			mainBlock = ""
			explanationBlock = ""
		} else {
			mainBlock += explanationBlock + exampleBlock
		}
	}
	messages = append(messages, MakeRequestTelegramText(info.Text, mainBlock+"\n", chatId))
	return messages
}

func GetMultitranOptionBlock(chatId int, page multitran.Page) []RequestTelegramText {
	var messages []RequestTelegramText
	var mainBlock string
	for _, info := range page.Options {
		mainBlock += fmt.Sprintf("❗️*Word*\\: *%s* \\[%s\\] \\(%s\\)", DecodeForTelegram(info.Text), DecodeForTelegram(info.Transcription), DecodeForTelegram(info.Type)) + "\n\n"
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
			mainBlock += "💡 *Explanation*:\n"

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

				if i < len(explanation.Text)-1 {
					mainBlock += ", "
				}
			}
			mainBlock += "\n\n" + GetRowSeparation() + "\n\n"
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
