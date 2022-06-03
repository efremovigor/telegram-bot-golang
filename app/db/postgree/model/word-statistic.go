package model

import "telegram-bot-golang/db/postgree"

type WordStatistic struct {
	Word  string
	Count int
}

func GetWordStatistics(limit int) (list []WordStatistic, err error) {

	connect := postgree.GetDbConnection()
	postgree.CloseConnection(connect)

	rows, err := connect.Query("SELECT w.name word , sum(us.requested) as count FROM user_statistic us LEFT JOIN word w on w.id = us.word_id GROUP BY w.name ORDER BY count LIMIt $1", limit)
	if err != nil {
		return
	}

	for rows.Next() {
		statistic := WordStatistic{}
		err = rows.Scan(&statistic.Word, &statistic.Count)
		if err != nil {
			return
		}
		list = append(list, statistic)
	}
	return
}

func GetWordStatisticsForUser(limit int, userId int) (list []WordStatistic, err error) {

	connect := postgree.GetDbConnection()
	postgree.CloseConnection(connect)

	rows, err := connect.Query("SELECT w.name word , sum(us.requested) as count FROM user_statistic us LEFT JOIN word w on w.id = us.word_id and us.user_id = $2 GROUP BY w.name ORDER BY count LIMIt $1", limit, userId)
	if err != nil {
		return
	}

	for rows.Next() {
		statistic := WordStatistic{}
		err = rows.Scan(&statistic.Word, &statistic.Count)
		if err != nil {
			return
		}
		list = append(list, statistic)
	}
	return
}
