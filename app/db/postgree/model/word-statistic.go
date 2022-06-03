package model

import (
	"fmt"
	"telegram-bot-golang/db/postgree"
)

type WordStatistic struct {
	Word  string
	Count int
}

func GetWordStatistics(limit int) (list []WordStatistic, err error) {

	connect := postgree.GetDbConnection()
	defer postgree.CloseConnection(connect)

	rows, err := connect.Query("SELECT w.name word , sum(us.requested) as count FROM user_statistic us LEFT JOIN word w on w.id = us.word_id GROUP BY w.name ORDER BY count DESC LIMIT $1", limit)
	if err != nil {
		fmt.Println("postgree:don't find any rows:" + err.Error())
		return
	}

	for rows.Next() {
		statistic := WordStatistic{}
		err = rows.Scan(&statistic.Word, &statistic.Count)
		if err != nil {
			fmt.Println("postgree:didn't able to write response:" + err.Error())
			return
		}
		list = append(list, statistic)
	}
	return
}

func GetWordStatisticsForUser(limit int, userId int) (list []WordStatistic, err error) {

	connect := postgree.GetDbConnection()
	defer postgree.CloseConnection(connect)

	rows, err := connect.Query("SELECT w.name word , sum(us.requested) as count FROM user_statistic us LEFT JOIN word w on w.id = us.word_id and us.user_id = $2 GROUP BY w.name ORDER BY count DESC LIMIT $1", limit, userId)
	if err != nil {
		fmt.Println("postgree:don't find any rows:" + err.Error())
		return
	}

	for rows.Next() {
		statistic := WordStatistic{}
		err = rows.Scan(&statistic.Word, &statistic.Count)
		if err != nil {
			fmt.Println("postgree:didn't able to write response:" + err.Error())
			return
		}
		list = append(list, statistic)
	}
	return
}
