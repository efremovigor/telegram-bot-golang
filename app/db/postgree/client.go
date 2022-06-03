package postgree

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"telegram-bot-golang/env"
)

const DatetimeLayer = "2006-01-02 15:04:05.999999"

func getDbConnectSource() string {
	return "host=db user=" + env.GetEnvVariable("DB_USER") +
		" password=" + env.GetEnvVariable("DB_PW") +
		" dbname=" + env.GetEnvVariable("DB_NAME") +
		" sslmode=disable"
}

func GetDbConnection() (db *sql.DB) {
	db, err := sql.Open("postgres", getDbConnectSource())
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	return
}

func CloseConnection(connect *sql.DB) {
	err := connect.Close()
	if err != nil {
		fmt.Println("postgree:error close connection:" + err.Error())
	}
}
