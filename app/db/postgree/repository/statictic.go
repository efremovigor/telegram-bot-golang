package repository

import (
	"fmt"
	"telegram-bot-golang/db/postgree/model"
)

func SaveNewWord(userId int, newWord string) {
	var err error
	var word model.Word
	var statistic model.UserStatistic
	if word, err = model.GetWord(newWord); err != nil {
		word = model.NewWord(newWord)
		err = word.SaveWord()
		if err != nil {
			fmt.Println("postgree:error of save word:" + err.Error())
		}
		statistic = model.NewUserStatistic(word, userId)

		err = statistic.SaveUserStatistic()
		if err != nil {
			fmt.Println("postgree:error of save statistic:" + err.Error())
		}
		return
	} else {
		if statistic, err = model.GetUserStatistic(word, userId); err != nil {
			statistic = model.NewUserStatistic(word, userId)

			err = statistic.SaveUserStatistic()
			if err != nil {
				fmt.Println("postgree:error of save statistic:" + err.Error())
			}
			return
		} else {
			fmt.Println(err)
		}
	}
	statistic.IncrRequested()
	err = statistic.SaveUserStatistic()
	if err != nil {
		fmt.Println("postgree:error of save statistic:" + err.Error())
	}
}
