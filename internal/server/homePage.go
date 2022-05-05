package server

import (
	"database/sql"
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//HomePage ..
func (handle *Handle) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	t, err := template.ParseFiles("./templates/html/homePage.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	posts, err := database.SelectPosts(handle.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	categories := []string{"cats", "IT", "cybersport", "github", "Anime", "Alem"}

	homeInfo := forum.Homepage{
		Posts:      posts,
		Categories: categories,
	}

	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		if err := t.Execute(w, homeInfo); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		return
	}

	homeInfo, err = UserHomePageInSession(handle.DB, cookie, posts, categories)
	if err != nil {
		log.Println(err)
		return
	}

	if err := t.Execute(w, homeInfo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}

//UserHomePageInSession ..
func UserHomePageInSession(db *sql.DB, cookie *http.Cookie, posts []forum.Post, categories []string) (forum.Homepage, error) {
	var idArr []int
	var likedpostNotification []forum.PostNotification
	var dislikedpostNotification []forum.PostNotification
	var commentedPostNotification []forum.CommentNotification
	var homeInfo forum.Homepage
	var visit bool = false
	var flag3 bool
	var countOfNotifications int

	user, err := database.GetUserFromDB(db, cookie)
	if err != nil {
		return homeInfo, err
	}
	flag1 := database.IsPostRaitingTableEmpty(db, user)
	flag2 := database.IsCommentsTableEmpty(db, user)

	if !flag1 {
		likedpostNotification, idArr, err = database.SelectIncomingLikedPostNotifications(db, user, visit)
		if err != nil {
			log.Println(err)
			return homeInfo, err
		}
		database.InsertPostNotifyToHistory(db, true)
		err = database.UpdateIncomingPostNotifications(db, user, idArr)
		if err != nil {
			log.Println(err)
			return homeInfo, err
		}
		dislikedpostNotification, idArr, err = database.SelectIncomingDislikedPostNotifications(db, user, visit)
		if err != nil {
			log.Println(err)
			return homeInfo, err
		}
		database.InsertPostNotifyToHistory(db, false)
		err = database.UpdateIncomingPostNotifications(db, user, idArr)
		if err != nil {
			log.Println(err)
			return homeInfo, err
		}
	}

	if !flag2 {
		commentedPostNotification, idArr, err = database.SelectIncomingCommentedPostNotifications(db, user, visit)
		if err != nil {
			log.Println(err)
			return homeInfo, err
		}
		database.InsertCommentNotifyToHistory(db)
		err = database.UpdateIncomingCommentPostNotifications(db, user, idArr)
		if err != nil {
			log.Println(err)
			return homeInfo, err
		}
	}

	if len(likedpostNotification) != 0 || len(dislikedpostNotification) != 0 || len(commentedPostNotification) != 0 {
		countOfNotifications = len(likedpostNotification) + len(dislikedpostNotification) + len(commentedPostNotification)
		flag3 = true
	}

	homeInfo = forum.Homepage{
		Posts:                    posts,
		Categories:               categories,
		InSession:                true,
		LikedPostNotification:    likedpostNotification,
		DislikedPostNotification: dislikedpostNotification,
		CommentNotification:      commentedPostNotification,
		IsNotification:           flag3,
		CountOfNotifications:     countOfNotifications,
	}
	return homeInfo, nil
}
