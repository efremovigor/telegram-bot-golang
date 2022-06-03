package model

import (
	"telegram-bot-golang/db/postgree"
	"time"
)

type UserStatistic struct {
	Id        int    `db:"id"`
	UserId    int    `db:"user_id"`
	WordId    int    `db:"word_id"`
	Requested int    `db:"requested"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func NewUserStatistic(word Word, userId int) (userStatistic UserStatistic) {
	userStatistic.UserId = userId
	userStatistic.WordId = word.Id
	userStatistic.Requested = 1
	userStatistic.SetDefaultDates()
	return
}

func (userStatistic *UserStatistic) IncrRequested() {
	userStatistic.Requested++
}

func (userStatistic *UserStatistic) SetDefaultDates() {
	userStatistic.SetDefaultCreated()
	userStatistic.SetDefaultUpdated()
}

func (userStatistic *UserStatistic) SetDefaultCreated() {
	userStatistic.CreatedAt = time.Now().Format(postgree.DatetimeLayer)
}

func (userStatistic *UserStatistic) SetDefaultUpdated() {
	userStatistic.UpdatedAt = time.Now().Format(postgree.DatetimeLayer)
}

func (userStatistic *UserStatistic) SaveUserStatistic() (err error) {
	connect := postgree.GetDbConnection()
	defer postgree.CloseConnection(connect)
	if userStatistic.Id != 0 {
		userStatistic.SetDefaultUpdated()
		connect.QueryRow(`UPDATE user_statistic SET requested = $2, updated_at = $3 WHERE id = $1`, userStatistic.Id, userStatistic.Requested, userStatistic.UpdatedAt)
		return nil
	} else {
		sqlStatement := `INSERT INTO user_statistic (user_id, word_id, requested, created_at, updated_at) VALUES ($1, $2, $3, $4, $5 ) RETURNING id`
		return connect.QueryRow(sqlStatement, userStatistic.UserId, userStatistic.WordId, userStatistic.Requested, userStatistic.CreatedAt, userStatistic.UpdatedAt).Scan(&userStatistic.Id)
	}
}

func GetUserStatistic(word Word, userId int) (userStatistic UserStatistic, err error) {
	connect := postgree.GetDbConnection()
	defer postgree.CloseConnection(connect)
	row := connect.QueryRow("SELECT id, user_id, word_id, requested, created_at , updated_at FROM user_statistic WHERE user_id = $1 and word_id = $2", userId, word.Id)
	userStatistic = UserStatistic{}
	err = row.Scan(&userStatistic.Id, &userStatistic.UserId, &userStatistic.WordId, &userStatistic.Requested, &userStatistic.CreatedAt, &userStatistic.UpdatedAt)
	return
}
