package statistic

import (
	"telegram-bot-golang/db/postgree/repository"
)

func Consider(key string, userId int) {
	repository.SaveNewWord(userId, key)
}
