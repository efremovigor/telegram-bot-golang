package telegram

import (
	"fmt"
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/service/dictionary/cambridge"
)

func GetBaseMsg(name string, id int) string {
	return fmt.Sprintf("Hey, [%s](tg://user?id=%d)\n", name, id) + "\n"
}
func GetIGotYourNewRequest(base string) string {
	return fmt.Sprintf(
		DecodeForTelegram("I got your message: ")+"*%s*\n", DecodeForTelegram(base))
}

func GetBlockWithRapidInfo(translate string) string {
	return fmt.Sprintf(
		DecodeForTelegram("Rapid-microsoft: ")+"*%s*\n\n", DecodeForTelegram(translate))
}

func GetCambridgeHeaderBlock(info cambridge.CambridgeInfo) string {
	return fmt.Sprintf("Cambridge\\-dictionary\\: *%s*", info.RequestText) + "\n"
}

func GetCambridgeOptionBlock(info cambridge.Info) string {
	var mainBlock string
	mainBlock += fmt.Sprintf("*Word*\\: *%s* \\[%s\\] \\(%s\\)", DecodeForTelegram(info.Text), DecodeForTelegram(info.Transcription), DecodeForTelegram(info.Type)) + "\n"
	if len(info.Explanation) > 0 {
		listExplanation := info.Explanation
		if len(listExplanation) > 6 {
			listExplanation = info.Explanation[0:5]
		}
		for n, explanation := range listExplanation {
			if n > 0 {
				mainBlock += GetRowSeparation()
			}
			mainBlock += GetFieldIfCan(explanation.Level, "Level")
			mainBlock += GetFieldIfCan(explanation.SemanticDescription, "Semantic")
			mainBlock += GetFieldIfCan(explanation.Description, "Description")
			mainBlock += GetFieldIfCan(explanation.Translate, "Translate")
			if len(explanation.Example) > 0 {
				mainBlock += GetFieldIfCan(explanation.Example[0], "Example")
			}
		}
	}
	return mainBlock + "\n"
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
	return DecodeForTelegram("-+-+-+-+-+-") + "\n"
}
