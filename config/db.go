package config

import (
	"database/sql"
	"fmt"
)

var Db *sql.DB

const (
	DB_USER     = "yp180102"
	DB_PASSWORD = "OiQr2df1DSy4Q0"
	DB_NAME     = "tokopedia-dev-db"
	DB_HOST     = "devel-postgre.tkpd"
	DB_PORT     = "5432"
)

func Init() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT)
	var err error
	Db, err = sql.Open("postgres", dbinfo)
	fmt.Println("Connected to database!")
	if err != nil {
		fmt.Println(err.Error())
	}
}
