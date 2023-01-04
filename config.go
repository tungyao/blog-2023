package main

import (
	"database/sql"
	"log"
)

func InitDb() {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Fatal(err)
	}
	Db = db
}
