package statistic

import (
	"telegram-bot-golang/db/postgree/repository"
)

func Consider(newWord string, userId int) {
	repository.SaveNewWord(userId, newWord)
}
