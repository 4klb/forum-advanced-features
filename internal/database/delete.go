package database

import (
	"database/sql"
	"log"
	"net/http"
)

func DeleteCategoryByPostId(postid int, db *sql.DB, tx *sql.Tx) {
	_, err := db.Exec(`DELETE FROM Categories
	WHERE PostId = ?`, postid)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	}
}

//DeletePost ..
func DeletePost(db *sql.DB, postid int, cookie *http.Cookie) error {
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.Exec("DELETE FROM Post WHERE Post.id = ? AND AuthorId = ?", postid, user.ID)
	if err != nil {
		return err
	}
	return nil
}

//DeletePost ..
func DeleteComment(db *sql.DB, commentid int, cookie *http.Cookie) error {
	user, err := GetUserFromDB(db, cookie)
	_, err = db.Exec("DELETE FROM Comments WHERE id = ? AND CommentatorId = ?", commentid, user.ID)
	if err != nil {
		return err
	}
	return nil
}
