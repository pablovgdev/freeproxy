package db

import "database/sql"

var db *sql.DB

func GetConnection() *sql.DB {
	if db != nil {
		return db
	}
	db, err := sql.Open("sqlite3_extended", "db.sqlite")
	if err != nil {
		panic(err)
	}
	return db
}
