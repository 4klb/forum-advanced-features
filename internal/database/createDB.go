package database

import (
	"database/sql"
	"log"
)

func CreateDb(db *sql.DB, create chan int) {
	select {
	case <-create:
		if err := CreateTables(db); err != nil {
			log.Println(err)
			return
		}
	}
}

//CreateTables ..
func CreateTables(db *sql.DB) error {

	_, err := db.Exec(`CREATE TABLE SignInUser (
		id integer primary key autoincrement not null,
		Name text not null UNIQUE,
		Login     text not null UNIQUE,
		Password  text ,
		Confirm   text
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE Cookie (
		Value text not null,
		Expires DATETIME not null,
		UserId integer UNIQUE,
		FOREIGN KEY(UserId) REFERENCES SignInUser(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE Post (
		id integer primary key autoincrement not null,
		Title text not null,
		Post text not null,
		AuthorId integer,
		FOREIGN KEY(AuthorId) REFERENCES SignInUser(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE Categories (
		id integer primary key autoincrement not null,
		Name text not null,
		PostId integer,
		FOREIGN KEY(PostId) REFERENCES Post(id) ON DELETE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE PostRaiting (
		id integer primary key autoincrement not null,
		[Like] bool,
		DisLike bool,
		PostId integer,
		LikerId integer, 
		FOREIGN KEY(LikerId) REFERENCES SignInUser(id) ON DELETE CASCADE,
		FOREIGN KEY(PostId) REFERENCES Post(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE VisitedNotificationsPost (
		id integer primary key autoincrement not null,
		Visited bool DEFAULT false,
		PostRaitingId integer,
		FOREIGN KEY(PostRaitingId) REFERENCES PostRaiting(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE Comments (
		id integer primary key autoincrement not null,
		Comment text not null,
		CommentatorId integer,
		PostId integer,
		FOREIGN KEY(CommentatorId) REFERENCES SignInUser(id) ON DELETE CASCADE,
		FOREIGN KEY(PostId) REFERENCES Post(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE CommentRaiting (
		id integer primary key autoincrement not null,
		[Like] bool,
		DisLike bool,
		CommentsId integer,
		UserId integer,
		FOREIGN KEY(UserId) REFERENCES SignInUser(id) ON DELETE CASCADE,
		FOREIGN KEY(CommentsId) REFERENCES Comments(id) ON DELETE CASCADE
		)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE VisitedNotificationsComment (
		id integer primary key autoincrement not null,
		Visited bool DEFAULT false,
		CommentsId integer,
		FOREIGN KEY(CommentsId) REFERENCES Comments(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE UserNotificationsHistory (
		id integer primary key autoincrement not null,
		VisitedNotificationsPostId integer,
		VisitedNotificationsCommentId integer,
		FOREIGN KEY(VisitedNotificationsPostId) REFERENCES VisitedNotificationsPost(id) ON DELETE CASCADE,
		FOREIGN KEY(VisitedNotificationsCommentId) REFERENCES VisitedNotificationsComment(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	return nil
}
