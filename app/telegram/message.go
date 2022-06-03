package telegram

import (
	"fmt"
	"telegram-bot-golang/db/postgree/model"
	"telegram-bot-golang/service/dictionary/cambridge"
)

func GetBaseMsg(name string, id int) string {
	return fmt.Sprintf("Hey, [%s](tg://user?id=%d)\n", name, id) +
		DecodeForTelegram("\n")
}
func GetIGotYourNewRequest(base string) string {
	return fmt.Sprintf(
		DecodeForTelegram("I got your message: ")+"*%s*\n", DecodeForTelegram(base))
}

func GetBlockWithRapidInfo(translate string) string {
	return fmt.Sprintf(
		DecodeForTelegram("Translates of rapid-microsoft: ")+"*%s*\n\n", DecodeForTelegram(translate))
}

func GetBlockWithCambridge(info cambridge.Info) string {
	mainBlock := DecodeForTelegram("Information from cambridge-dictionary: ")
	if info.IsValid() {
		mainBlock += fmt.Sprintf("*%s*", DecodeForTelegram(info.Text)) + "\n"
		mainBlock += GetFieldIfCan(info.Type, "Type")
		mainBlock += GetFieldIfCan(info.Transcription, "Transcription")
		if len(info.Explanation) > 0 {
			for n, explanation := range info.Explanation {
				if n > 0 {
					mainBlock += DecodeForTelegram("-+-+-+-+-+-") + "\n"
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
	} else {
		mainBlock += DecodeForTelegram("*-*") + "\n"
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
		return fmt.Sprintf("*Top %d words used*/n", n)
	}
	return fmt.Sprintf("*My %d words used*/n", n)
}
func GetRowRating(n int, statistic model.WordStatistic) string {
	return fmt.Sprintf("*%d*\\. %s \\- %d time/n", n, statistic.Word, statistic.Count)
}
