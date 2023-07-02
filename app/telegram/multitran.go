package telegram

import (
	"fmt"
	"telegram-bot-golang/helper"
	"telegram-bot-golang/service/dictionary/multitran"
)

func GetMultitranHeaderBlock(word string) string {
	return fmt.Sprintf("✅ *Multitran\\-dictionary*\\: *%s*", DecodeForTelegram(word)) + "\n\n"
}

func GetMultitranShortInfo(chatId int, page multitran.Page) []RequestTelegramText {
	var messages []RequestTelegramText
	var mainBlock string
	mainBlock += GetMultitranHeaderBlock(page.Options[0].Text)
	mainBlock += fmt.Sprintf("*Word*\\: *%s* \\[%s\\] \\(%s\\)", DecodeForTelegram(page.Options[0].Text), DecodeForTelegram(page.Options[0].Transcription), DecodeForTelegram(page.Options[0].Type)) + "\n\n"
	if len(page.Options[0].Explanation) > 0 {
		mainBlock += GetFieldIfCan(page.Options[0].Explanation[0].Type, "Type")
		for i, translate := range page.Options[0].Explanation[0].Text {
			mainBlock += DecodeForTelegram(translate)

			if i < len(page.Options[0].Explanation[0].Text)-1 {
				mainBlock += ", "
			}
			if i > 5 {
				break
			}
		}
	}

	messages = append(messages,
		MakeRequestTelegramText(
			page.Options[0].Text,
			mainBlock+"\n",
			chatId,
			[]Keyboard{{Text: "🗂 full", CallbackData: NextMessageFullMultitran + " " + page.Options[0].Text}},
		))
	return messages
}

func GetMultitranOptionBlock(chatId int, page multitran.Page) []RequestTelegramText {
	var messages []RequestTelegramText
	var mainBlock string
	for _, info := range page.Options {
		mainBlock += fmt.Sprintf("*Word*\\: *%s* \\[%s\\] \\(%s\\)", DecodeForTelegram(info.Text), DecodeForTelegram(info.Transcription), DecodeForTelegram(info.Type)) + "\n\n"
		for _, explanation := range info.Explanation {
			if helper.Len(mainBlock) > MaxRequestSize {
				messages = append(messages,
					MakeRequestTelegramText(
						info.Text,
						mainBlock+"\n",
						chatId,
						[]Keyboard{},
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
							[]Keyboard{},
						),
					)
					mainBlock = ""
				}
				mainBlock += DecodeForTelegram(translate)

				if i < len(explanation.Text)-1 {
					mainBlock += ", "
				}
			}
			mainBlock += "\n\n" + GetRowSeparation() + "\n"
		}
	}
	messages = append(messages, MakeRequestTelegramText(
		page.RequestText,
		mainBlock+"\n",
		chatId,
		[]Keyboard{},
	))
	return messages
}
