package server

import (
	"database/sql"
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//Profile ..
func (handle *Handle) Profile(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("./templates/html/profile.html")
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}
		ditails := GetProfileInfo(handle.DB, cookie)
		if err := t.Execute(w, ditails); err != nil {
			w.WriteHeader(500)
			return
		}
	default:
		w.WriteHeader(405)
		return
	}
}

//GetProfileInfo ..
func GetProfileInfo(db *sql.DB, cookie *http.Cookie) forum.UserProfile {
	var ditails forum.UserProfile
	var flag bool

	user, err := database.GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return ditails
	}
	createdPosts, err := database.SelectPostByUser(db, user)
	if err != nil {
		log.Println(err)
		return ditails
	}
	likedPosts, err := database.SelectLikedPostsByUser(db, user)
	if err != nil {
		log.Println(err)
		return ditails
	}
	dislikedPosts, err := database.SelectDislikedPostsByUser(db, user)
	if err != nil {
		log.Println(err)
		return ditails
	}
	commentedPosts, err := database.SelectCommentedPostsByUser(db, user)
	if err != nil {
		log.Println(err)
		return ditails
	}
	notifications, err := database.SelectNotificationsHistory(db, user)
	if err != nil {
		log.Println(err)
		return ditails
	}
	if len(notifications) != 0 {
		flag = true
	}
	ditails = forum.UserProfile{
		User:           user,
		CreatedPosts:   createdPosts,
		LikedPosts:     likedPosts,
		DislikedPosts:  dislikedPosts,
		CommentedPosts: commentedPosts,
		Notifications:  notifications,
		IsNotification: flag,
	}
	return ditails
}
