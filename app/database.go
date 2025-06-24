package app

import (
	"database/sql"

	"github.com/fardannozami/ticket-booking/helper"
)

func NewSqlDb() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ticket-booking")
	helper.PanicIfError(err)

	return db
}
