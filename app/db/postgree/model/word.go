package model

import (
	"telegram-bot-golang/db/postgree"
	"time"
)

type Word struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (word *Word) SetDefaultDates() {
	word.SetDefaultCreated()
	word.SetDefaultUpdated()
}

func (word *Word) SetDefaultCreated() {
	word.CreatedAt = time.Now().Format(postgree.DatetimeLayer)
}

func (word *Word) SetDefaultUpdated() {
	word.UpdatedAt = time.Now().Format(postgree.DatetimeLayer)
}

func NewWord(text string) (word Word) {
	word.Name = text
	word.SetDefaultDates()
	return
}

func GetWord(name string) (word Word, err error) {
	connect := postgree.GetDbConnection()
	defer postgree.CloseConnection(connect)
	row := connect.QueryRow("SELECT id, name, created_at , updated_at FROM word WHERE name = $1", name)
	err = row.Scan(&word.Id, &word.Name, &word.CreatedAt, &word.UpdatedAt)
	return
}

func (word *Word) SaveWord() (err error) {
	connect := postgree.GetDbConnection()
	defer postgree.CloseConnection(connect)
	if word.Id != 0 {
		word.SetDefaultUpdated()
		connect.QueryRow(`UPDATE word SET name=$2, updated_at = $3 WHERE id = $1`, word.Id, word.Name, word.UpdatedAt)
		return nil
	} else {
		sqlStatement := `INSERT INTO word (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`
		return connect.QueryRow(sqlStatement, word.Name, word.CreatedAt, word.UpdatedAt).Scan(&word.Id)
	}
}
