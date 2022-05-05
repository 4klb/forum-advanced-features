package database

import (
	"database/sql"
	forum "forum/internal"
	"log"
	"net/http"
	"time"
)

//RemoveExpiredCookie ..
func RemoveExpiredCookie(db *sql.DB, create chan int) {
	for {
		_, err := db.Exec("DELETE FROM Cookie WHERE Expires <?", time.Now())
		if err != nil {
			create <- 0
		}
		time.Sleep(10 * time.Second)
	}
}

//IsUserInSession ..
func IsUserInSession(db *sql.DB, cookie *http.Cookie) bool {
	if cookie == nil {
		return false
	}
	var val string
	err := db.QueryRow("SELECT Value FROM Cookie WHERE Value = ?", cookie.Value).Scan(&val)
	if err == nil && err == sql.ErrNoRows {
		return false
	}
	if val == "" {
		return false
	}
	return true
}

//DeleteCookieFromDB ..
func DeleteCookieFromDB(db *sql.DB, cookie *http.Cookie) {
	_, err := db.Exec("DELETE FROM Cookie WHERE Value = ?", cookie.Value)
	if err != nil {
		log.Println(err)
		return
	}
}

//GetUserFromDB ..
func GetUserFromDB(db *sql.DB, cookie *http.Cookie) (forum.User, error) {
	var user forum.User

	row := db.QueryRow("SELECT SignInUser.Id, SignInUser.Name FROM SignInUser JOIN Cookie on UserId = id WHERE Value = ? ", cookie.Value)
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		log.Println(err)
		return user, err
	}

	// if err := row.Err(); err != nil {
	// 	return user, err
	// }

	return user, nil
}

//IsPreviousLikeComment ..
func IsPreviousLikedComment(db *sql.DB, yes, no bool, postID, userID int) bool {
	var like, dislike bool
	firstReq := "SELECT Like , DisLike FROM CommentRaiting WHERE CommentsId = ? AND UserId = ?"
	secondReq := "DELETE FROM CommentRaiting WHERE CommentsId = ? AND UserId = ?"
	row := db.QueryRow(firstReq, postID, userID)
	if err := row.Scan(&like, &dislike); err != nil {
		return true
	}
	_, err := db.Exec(secondReq, postID, userID)
	if err != nil {
		return false
	}
	if yes && like {
		return false
	} else if yes && !like {
		return true
	} else if no && dislike {
		return false
	} else if no && !dislike {
		return true
	}
	return true
}

//IsPreviousLikedPost ..
func IsPreviousLikedPost(db *sql.DB, yes, no bool, postID, userID int) bool {
	var like, dislike bool
	firstReq := "SELECT Like , DisLike FROM PostRaiting WHERE PostId = ? AND LikerId = ?"
	secondReq := "DELETE FROM PostRaiting WHERE PostId = ? AND LikerId = ?"
	row := db.QueryRow(firstReq, postID, userID)
	if err := row.Scan(&like, &dislike); err != nil {
		return true
	}
	_, err := db.Exec(secondReq, postID, userID)
	if err != nil {
		return false
	}
	if yes && like {
		return false
	} else if yes && !like {
		return true
	} else if no && dislike {
		return false
	} else if no && !dislike {
		return true
	}
	return true
}

//IsUserAuthorOfPost ..
func IsUserAuthorOfPost(db *sql.DB, postid int, user forum.User) bool {
	var title string
	row := db.QueryRow("SELECT Title FROM Post WHERE id =? AND AuthorId =?", postid, user.ID)
	if err := row.Scan(&title); err != nil {
		log.Println(err)
		return false
	}

	// if err := row.Err(); err != nil {
	// 	log.Println(err)
	// 	return false
	// }

	if title == "" {
		return false
	}

	return true
}

func IsUserAuthorOfComment(db *sql.DB, postid int, user forum.User) ([]forum.Comment, error) {
	var comments []forum.Comment
	rows, err := db.Query(`SELECT id, Comment FROM Comments
	WHERE PostId = ? AND CommentatorId = ?`, postid, user.ID)
	if err != nil {
		log.Println(err)
		return comments, err
	}
	for rows.Next() {
		var id int
		var comment string
		if err := rows.Scan(&id, &comment); err != nil {
			log.Println(err)
			return comments, err
		}
		// if err := rows.Err(); err != nil {
		// 	log.Println(err)
		// 	return false, err
		// }
		com := forum.Comment{
			ID:   id,
			Text: comment,
		}
		comments = append(comments, com)

		// if comment == "" {
		// 	return false, err
		// }
	}
	return comments, nil
}

/*

	RemoveExpiredCookie гоурутина которая каждые 30 секунд проверяет куки из БД
	если время жизни куки уже подошло к концу то удаляем ее из БД

	IsPreviousLikesDislikedPost

	DeleteCookieFromDB удаляем куки с БД если пользователь сам захотел выйти с сайта

	GetUserFromDB получаем айди пользователя исходя из его куки

	IsPreviousLikesDislikedPost проверяем последнее значение лайка и дизлайка
	пользователя из БД и делаем что то типа switch case

*/
