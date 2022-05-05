package database

import (
	"database/sql"
	"errors"
	"fmt"
	forum "forum/internal"
	"log"
	"net/http"
)

func UpdateIncomingPostNotifications(db *sql.DB, user forum.User, idArr []int) error {
	for _, id := range idArr {
		_, err := db.Exec(`UPDATE VisitedNotificationsPost
		SET Visited = true 
		WHERE PostRaitingId  = ?`, id)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func UpdateIncomingCommentPostNotifications(db *sql.DB, user forum.User, idArr []int) error {
	for _, id := range idArr {
		fmt.Println(id)
		_, err := db.Exec(`UPDATE VisitedNotificationsComment
		SET Visited = true 
		WHERE CommentsId  = ?`, id)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func UpdatePost(db *sql.DB, post *forum.CreatePost, cookie *http.Cookie, postid int) error {
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	err = IsPostValid(db, post)
	if err != nil {
		log.Println(err)
		return err
	}
	err = EditPostInfo(db, post, postid, user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//EditPostInfo ..
func EditPostInfo(db *sql.DB, post *forum.CreatePost, postid int, user forum.User) error {
	_, err := db.Exec(`UPDATE Post
	Set Title = ?, Post = ?
	WHERE id = ? AND AuthorId = ?`, post.Title, post.Text, postid, user.ID)
	if err != nil {
		log.Println(err)
		return errors.New("Can not edit post")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	DeleteCategoryByPostId(postid, db, tx)
	InsertIntoCaregories(post, postid, db, tx)
	return nil
}

func UpdateComment(db *sql.DB, comment *forum.UpdateComment, cookie *http.Cookie) error {
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	err = IsCommentValid(db, comment)
	if err != nil {
		log.Println(err)
		return err
	}
	err = EditCommentInfo(db, comment, user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func EditCommentInfo(db *sql.DB, comment *forum.UpdateComment, user forum.User) error {
	_, err := db.Exec(`UPDATE Comments
	Set Comment = ?
	WHERE id = ? AND CommentatorId = ?`, comment.Comment.Text, comment.Comment.ID, user.ID)
	if err != nil {
		log.Println(err)
		return errors.New("Can not edit comment")
	}
	return nil
}
