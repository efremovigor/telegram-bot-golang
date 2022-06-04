package telegram

import (
	"fmt"
	"telegram-bot-golang/command"
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
			listExplanation := info.Explanation
			if len(listExplanation) > 5 {
				listExplanation = info.Explanation[0:4]
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
			if len(info.Explanation) > 5 {
				mainBlock += "\n" + DecodeForTelegram("...this is short info...") + "\n"
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

func GetHelpHeader() string {
	return "*List of commands available to you:*\n" +
		GetRowSeparation() +
		command.RuEnCommand + fmt.Sprintf("Change translate of transition %s \n", DecodeForTelegram(command.Transitions()[command.RuEnCommand].Desc)) +
		command.EnRuCommand + fmt.Sprintf("Change translate of transition %s \n", DecodeForTelegram(command.Transitions()[command.EnRuCommand].Desc)) +
		command.HelpCommand + "Show all the available commands\n" +
		command.GetAllTopCommand + "To see the most popular requests for translation or explanation  \n" +
		command.GetMyTopCommand + "To see your popular requests for translation or explanation  \n"
}
