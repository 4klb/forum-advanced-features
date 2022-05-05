package server

import (
	"database/sql"
	forum "forum/internal"
)

//Handle ..
type Handle struct {
	DB   *sql.DB
	Post forum.Post
}

//CreateHandle ..
func CreateHandle(db *sql.DB) *Handle {
	return &Handle{
		DB: db,
	}
}
