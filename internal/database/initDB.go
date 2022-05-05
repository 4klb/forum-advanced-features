package database

import (
	"database/sql"

	//Driver for sqlite3
	_ "github.com/mattn/go-sqlite3"
)

//DbInit ..
func DbInit() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./forum.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	return db, nil
}
