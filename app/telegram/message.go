package telegram

import (
	"fmt"
	"strings"
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/db/redis"
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

func GetCambridgeShortInfo(chatId int, page cambridge.Page) []RequestTelegramText {
	var messages []RequestTelegramText
	var buttons []Keyboard
	var mainBlock string
	if len(page.Options) == 0 {
		return messages
	}
	// для короткой информации достаточно одного
	mainBlock += GetCambridgeHeaderBlock(page.Options[0].Text)
	mainBlock += GetFieldIfCan(page.Options[0].Type, "Type")
	for lang, transcription := range page.Options[0].Transcription {
		mainBlock += fmt.Sprintf("*%s*:\\[%s\\] ", strings.ToUpper(lang), DecodeForTelegram(transcription)) + "\n"
		break
	}
	if len(page.Options[0].Explanation) > 0 {
		mainBlock += GetFieldIfCan(page.Options[0].Explanation[0].Text, "Phrase")
		mainBlock += GetFieldIfCan(page.Options[0].Explanation[0].Level, "Level")
		mainBlock += GetFieldIfCan(page.Options[0].Explanation[0].SemanticDescription, "📃 Semantic")
		mainBlock += GetFieldIfCan(page.Options[0].Explanation[0].Description, "📃 Description")
		mainBlock += GetFieldIfCan(page.Options[0].Explanation[0].Translate, "💡 Translate")
		if len(page.Options[0].Explanation[0].Example) > 0 {
			mainBlock += "📌" + DecodeForTelegram(page.Options[0].Explanation[0].Example[0]) + "\n"
		}
	}
	var hasImage, hasVoice bool
	for _, info := range page.Options {
		if helper.Len(info.Image) > 0 && !hasImage {
			hash := helper.MD5(info.Image)
			redis.SetStruct(fmt.Sprintf(redis.InfoCambridgeUniqPicLink, hash), redis.PicFile{Word: info.Text, Url: info.Image}, 0)
			buttons = append(buttons, Keyboard{Text: "🏞 picture", CallbackData: ShowRequestPic + " " + hash})
			hasImage = true
		}
		if helper.Len(info.VoicePath.US) > 0 && !hasVoice {
			hash := helper.MD5(info.VoicePath.US)
			redis.SetStruct(fmt.Sprintf(redis.InfoCambridgeUniqVoiceLink, hash), redis.VoiceFile{Lang: CountryUs, Word: info.Text, Url: info.VoicePath.US}, 0)
			buttons = append(buttons, Keyboard{Text: "🗣 " + CountryUs, CallbackData: ShowRequestVoice + " " + CountryUs + " " + hash})
			hasVoice = true
		}
		if hasImage == true && hasVoice == true {
			break
		}
	}

	messages = append(messages,
		MakeRequestTelegramText(
			page.Options[0].Text,
			mainBlock+"\n",
			chatId,
			append(buttons, Keyboard{Text: "🗂 full", CallbackData: NextMessageFullCambridge + " " + page.Options[0].Text}),
		))
	return messages
}

func GetCambridgeOptionBlock(chatId int, info cambridge.Info) []RequestTelegramText {
	var messages []RequestTelegramText
	var mainBlock string
	var buttons []Keyboard
	mainBlock += fmt.Sprintf("*Word*\\: *%s*", DecodeForTelegram(info.Text))
	mainBlock += fmt.Sprintf("\\(%s\\)", DecodeForTelegram(info.Type)) + "\n\n"
	for lang, transcription := range info.Transcription {
		mainBlock += fmt.Sprintf("*%s*:\\[%s\\] ", strings.ToUpper(lang), DecodeForTelegram(transcription))
	}
	if len(info.Transcription) > 0 {
		mainBlock += "\n"
	}

	if helper.Len(info.Image) > 0 {
		hash := helper.MD5(info.Image)
		redis.SetStruct(fmt.Sprintf(redis.InfoCambridgeUniqPicLink, hash), redis.PicFile{Word: info.Text, Url: info.Image}, 0)
		buttons = append(buttons, Keyboard{Text: "🏞 picture", CallbackData: ShowRequestPic + " " + hash})
	}

	if helper.Len(info.VoicePath.UK) > 0 {
		hash := helper.MD5(info.VoicePath.UK)
		redis.SetStruct(fmt.Sprintf(redis.InfoCambridgeUniqVoiceLink, hash), redis.VoiceFile{Lang: CountryUk, Word: info.Text, Url: info.VoicePath.UK}, 0)
		buttons = append(buttons, Keyboard{Text: "🗣 " + CountryUk, CallbackData: ShowRequestVoice + " " + CountryUk + " " + hash})
	}
	if helper.Len(info.VoicePath.US) > 0 {
		hash := helper.MD5(info.VoicePath.US)
		redis.SetStruct(fmt.Sprintf(redis.InfoCambridgeUniqVoiceLink, hash), redis.VoiceFile{Lang: CountryUs, Word: info.Text, Url: info.VoicePath.US}, 0)
		buttons = append(buttons, Keyboard{Text: "🗣 " + CountryUs, CallbackData: ShowRequestVoice + " " + CountryUs + " " + hash})
	}

	var explanationBlock string
	for n, explanation := range info.Explanation {
		explanationBlock += GetFieldIfCan(explanation.Text, "Phrase")
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
					[]Keyboard{},
				),
				MakeRequestTelegramText(
					info.Text,
					exampleBlock+"\n",
					chatId,
					[]Keyboard{},
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
					buttons,
				),
			)
			mainBlock = ""
			mainBlock = ""
			explanationBlock = ""
		} else {
			mainBlock += explanationBlock + exampleBlock
			if n != len(info.Explanation)-1 {
				mainBlock += "\n" + GetRowSeparation() + "\n"
			}
		}
	}
	messages = append(messages, MakeRequestTelegramText(info.Text, mainBlock+"\n", chatId, buttons))
	return messages
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
