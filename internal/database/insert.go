package database

import (
	"database/sql"
	"errors"
	"fmt"
	forum "forum/internal"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//InsertNewUserIntoDB ..
func InsertNewUserIntoDB(db *sql.DB, user forum.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("When generating hash from password InsertNewUserIntoDB() %w : ", err)
	}
	_, err = db.Exec("INSERT INTO SignInUser (Name, Login , Password) VALUES (?,?,?)", user.Name, user.Login, string(hash))
	if err != nil {
		return fmt.Errorf("Sorry name or login is already taken")
	}
	return nil
}

//InsertUserBySocialNetworks ..
func InsertUserBySocialNetworks(db *sql.DB, user forum.User) error {
	_, err := db.Exec("INSERT INTO SignInUser (Name , Login , Confirm) VALUES(?,?,?)", user.Name, user.Login, user.Confirm)
	if err != nil {
		return err
	}
	return nil
}

//InsertCookieIntoDB ..
func InsertCookieIntoDB(db *sql.DB, login string, cookie *http.Cookie) {
	var id int
	row := db.QueryRow("SELECT Id FROM SignInUser WHERE Login =?;", login)
	row.Scan(&id)

	_, err := db.Exec("DELETE FROM Cookie WHERE UserId =?", id)
	if err != nil {
		log.Println()
		return
	}

	_, err = db.Exec("INSERT INTO Cookie (Value, Expires,UserId)VALUES(?,?,?)", cookie.Value, cookie.Expires, id)
	if err != nil {
		log.Println(err)
		return
	}
}

func IsPostValid(db *sql.DB, post *forum.CreatePost) error {
	title := strings.Trim(post.Title, " ")

	text := strings.Trim(post.Text, " ")

	if CheckForSpace(title) {
		return errors.New("Title must not be empty")
	} else if CheckForSpace(text) {
		return errors.New("Text must not be empty")
	}

	if len(title) > 30 {
		return errors.New("Length of title must be less than 30 characters")
	}

	if len(post.Categories) == 0 {
		return errors.New("Please choose at least one category")
	}
	if title == "" {
		return errors.New("Invalid title")

	} else if text == "" {
		return errors.New("Invalid description")
	}
	return nil
}

func IsCommentValid(db *sql.DB, comment *forum.UpdateComment) error {

	text := strings.Trim(comment.Comment.Text, " ")

	if CheckForSpace(text) {
		return errors.New("Text must not be empty")
	}

	if text == "" {
		return errors.New("Invalid description")
	}
	return nil
}

//InsertPostIntoDB ..
func InsertPostIntoDB(db *sql.DB, post *forum.CreatePost, cookie *http.Cookie) error {
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
	res, err := db.Exec("INSERT INTO Post (Title , Post, AuthorID) VALUES(?,?,?)", post.Title, post.Text, user.ID)
	if err != nil {
		return err
	}

	postid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if err := InsertCategoriesToDB(db, post.Categories, postid); err != nil {
		_, err := db.Exec("DELETE FROM Post WHERE Id =?; ", postid)
		if err != nil {
			return err
		}
	}
	return nil
}

//InsertCategoriesToDB ..
func InsertCategoriesToDB(db *sql.DB, categories []string, postid int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	for _, category := range categories {
		_, err := tx.Exec("INSERT INTO Categories (Name,PostId) VALUES (?,?)", category, postid)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	return nil
}

//InsertCommentToDB ..
func InsertCommentToDB(db *sql.DB, comment string, postID int, cookie *http.Cookie) error {
	test := strings.Trim(comment, " ")
	if test == "" {
		return errors.New("Invalid comment")
	}
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	userID := user.ID
	res, err := db.Exec("INSERT INTO Comments (Comment, CommentatorId, PostId) VALUES (?,?,?)", comment, userID, postID)
	if err != nil {
		log.Println(err)
		return err

	}
	x, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rateid := int(x)
	err = InsertCommentedNotification(db, rateid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func InsertCommentedNotification(db *sql.DB, rateid int) error {
	_, err := db.Exec(`INSERT INTO VisitedNotificationsComment
	(Visited, CommentsId)
	VALUES(?,?)`, false, rateid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//InsertCommentNotification ..
func InsertCommentNotification(db *sql.DB, like, dislike string, id int, cookie *http.Cookie) {
	var yes bool
	var no bool
	if len(like) != 0 {
		yes = true
	} else if len(dislike) != 0 {
		no = true
	}
	if !yes && !no {
		return
	}

	req := "INSERT INTO CommentRaiting (Like, DisLike, CommentsId, UserId) VALUES (?,?,?,?)"
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return
	}
	userID := user.ID
	do := IsPreviousLikedComment(db, yes, no, id, userID)
	if do {
		_, err := db.Exec(req, yes, no, id, userID)
		if err != nil {
			return
		}
	}
	return
}

//InsertPostNotification ..
func InsertPostNotification(db *sql.DB, like, dislike string, id int, cookie *http.Cookie) {
	var yes bool
	var no bool
	if len(like) != 0 {
		yes = true
	} else if len(dislike) != 0 {
		no = true
	}
	if !yes && !no {
		return
	}
	req := "INSERT INTO PostRaiting (Like, DisLike, PostId, LikerId) VALUES (?,?,?,?)"

	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return
	}
	userID := user.ID
	do := IsPreviousLikedPost(db, yes, no, id, userID)
	if do {
		res, err := db.Exec(req, yes, no, id, userID)
		if err != nil {
			return
		}
		x, err := res.LastInsertId()
		if err != nil {
			return
		}
		rateid := int(x)
		InsertLikedPostNotification(db, rateid)
		return
	}
	return
}

func InsertLikedPostNotification(db *sql.DB, rateid int) {
	_, err := db.Exec(`INSERT INTO VisitedNotificationsPost 
	(Visited, PostRaitingId) 
	VALUES(?,?)`, false, rateid)
	if err != nil {
		log.Println(err)
		return
	}
}

func InsertIntoCaregories(post *forum.CreatePost, postid int, db *sql.DB, tx *sql.Tx) {
	for _, v := range post.Categories {
		_, err := db.Exec("INSERT INTO Categories (Name, PostId) VALUES(?,?)", v, postid)
		if err != nil {
			log.Println(err)
			tx.Rollback()
		}
	}
}

func InsertPostNotifyToHistory(db *sql.DB, data bool) error {
	var idArr []int
	rows, err := db.Query(`SELECT VisitedNotificationsPost.id 
	FROM VisitedNotificationsPost
	JOIN PostRaiting ON PostRaitingId = PostRaiting.id 
	WHERE Visited = 0 AND Like = ?`, data)
	if err != nil {
		log.Println(err)
		return err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
			return err
		}
		idArr = append(idArr, id)
	}
	for _, id := range idArr {
		_, err = db.Exec(`INSERT INTO UserNotificationsHistory
	(VisitedNotificationsPostId)
	VALUES (?)`, id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func InsertCommentNotifyToHistory(db *sql.DB) error {
	var idArr []int
	rows, err := db.Query(`SELECT id 
	FROM VisitedNotificationsComment
	WHERE Visited = 0`)
	if err != nil {
		log.Println(err)
		return err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
			return err
		}
		idArr = append(idArr, id)
	}
	for _, id := range idArr {
		_, err = db.Exec(`INSERT INTO UserNotificationsHistory
	(VisitedNotificationsCommentId)
	VALUES (?)`, id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func CheckForSpace(s string) bool {
	var checker bool
	var counter int
	if s == "" {
		checker = true
	}
	for _, v := range s {
		if v == 32 {
			counter++
		}
	}
	if counter == len(s) {
		checker = true
	}
	return checker
}

/*

	InsertNewUserIntoDB после того как проверили на уникальность
	имя , логин и пароль добавляем пользователя в БД

	InsertCookieIntoDB получаем айди пользователя по логину и затем
	добавляем в БД куки со значением(сами куки) время истечения и номер пользователя

	InsertPostIntoDB получаем айди пользователя из куки и присваеваем посту
	название , содержимое и номер пользователя который оставил пост

	InsertCategoriesToDB добавляем категории в БД по названию категории и
	номеру поста под которым находится категория

	InsertCommentToDB получаем ID поста под которым надо оставить коммент
	и добавляем в БД само значем значение коммента и кто оставил

	InsertLikeDislikeIntoDB смотрим значение кнопок лайка и дизлайка
	в зависимости от клика добавлям в БД по значение кнопки, но перед этим
	проверить предыдущие лайки пользователя под этим постом
*/
