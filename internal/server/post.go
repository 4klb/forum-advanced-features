package server

import (
	forum "forum/internal"
	"forum/internal/database"

	// "forum/internal/server/checkforms"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//ChosenPost ..
func (handle *Handle) ChosenPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	t, err := template.ParseFiles("./templates/html/chosenpost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	var isPostOfUser bool
	var comments []forum.Comment

	cookie, _ := r.Cookie("session")
	id := r.URL.Path[6:]
	passed, postid := database.IsPostInDB(handle.DB, id)

	if database.IsUserInSession(handle.DB, cookie) {
		user, err := database.GetUserFromDB(handle.DB, cookie)
		if err != nil {
			log.Println(err)
			return
		}
		if database.IsUserAuthorOfPost(handle.DB, postid, user) {
			isPostOfUser = true
		}
		comments, err = database.IsUserAuthorOfComment(handle.DB, postid, user)
		if err != nil {
			log.Println(err)
			return
		}
	}

	if !passed {
		log.Println("page does not exist")
		w.WriteHeader(404)
		return
	}
	post, err := database.SelectOnePost(handle.DB, postid)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	for _, userComments := range comments {
		for i, allComments := range post.Comments {
			if userComments.ID == allComments.ID {
				post.Comments[i].IsCommentOfUser = true
			}
		}
	}

	chosenPost := forum.ChosenPost{
		Post:         post,
		IsPostOfUser: isPostOfUser,
	}

	if handle.Post.ErrorVal.Err != false {
		post.ErrorVal = handle.Post.ErrorVal
	}
	if err := t.Execute(w, chosenPost); err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	handle.Post.ErrorVal.Err = false
}

//RatePost ..
func (handle *Handle) RatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}

	if r.URL.Path != "/ratepost" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}

	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	prevURL := r.Header.Get("Referer")
	id := prevURL[27:]
	postid, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	like := r.FormValue("like")
	dislike := r.FormValue("dislike")

	var counter int

	for v := range r.Form {
		counter++
		if v != "like" && v != "dislike" {
			log.Println("An incoming value is not a required key")
			w.WriteHeader(400)
			return
		}
	}
	if counter == 0 {
		log.Println("Body is empty")
		w.WriteHeader(400)
		return
	}

	database.InsertPostNotification(handle.DB, like, dislike, postid, cookie)

	http.Redirect(w, r, prevURL, 301)
}

//EditPost ..
func (handle *Handle) EditPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	if r.URL.Path != "/editpost" {
		NotFound(w)
		w.WriteHeader(404)
		log.Println(404)
		return
	}

	t, err := template.ParseFiles("./templates/html/editpost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	postidval := r.FormValue("editpost")
	postid, err := strconv.Atoi(postidval)
	if err != nil {
		log.Println(404)
		w.WriteHeader(404)
		return
	}
	categories := []string{"cats", "IT", "cybersport", "github", "Anime", "Alem"}

	post, err := database.SelectOnePost(handle.DB, postid)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	unionEnterPost := forum.UnionEnterPost{
		Post:       post,
		ErrorVal:   forum.Error{},
		Categories: categories,
	}

	if err := t.Execute(w, unionEnterPost); err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
}

//ConfirmEditPost ..
func (handle *Handle) ConfirmEditPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/editpost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	cookie, _ := r.Cookie("session")

	postId := r.FormValue("confirm")
	back := r.FormValue("back")

	if len(back) != 0 {
		http.Redirect(w, r, "/post/"+back, 301)
		return
	}

	postid, err := strconv.Atoi(postId)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	categories := []string{"cats", "IT", "cybersport", "github", "Anime", "Alem"}

	post, checker := PostFormChecker(w, r, categories)
	if !checker {
		return
	}

	updatePost := forum.UpdatePost{
		PostId:     postid,
		Post:       post,
		ErrorVal:   forum.Error{},
		Categories: categories,
	}

	if err = database.UpdatePost(handle.DB, &post, cookie, postid); err != nil {
		updatePost.ErrorVal.Err = true
		updatePost.ErrorVal.MSG = err.Error()
		w.WriteHeader(400)

		if err := t.Execute(w, updatePost); err != nil {
			// w.WriteHeader(500) ?? superfluous
		}
		return
	}

	http.Redirect(w, r, "/post/"+postId, 301)
}

// DeletePost ..
func (handle *Handle) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}

	if r.URL.Path != "/deletepost" {
		NotFound(w)
		w.WriteHeader(404)
		return

	}

	t, err := template.ParseFiles("./templates/html/deletepost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	postId := r.FormValue("deletepost")
	postid, err := strconv.Atoi(postId)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	post, err := database.SelectOnePost(handle.DB, postid)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if err := t.Execute(w, post); err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

}

//ConfirmDeletePost ..
func (handle *Handle) ConfirmDeletePost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/deletepost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	cookie, _ := r.Cookie("session")

	postId := r.FormValue("confirm")
	back := r.FormValue("back")

	if len(back) != 0 {
		http.Redirect(w, r, "/post/"+back, 301)
		return
	}

	postid, err := strconv.Atoi(postId)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	post, err := database.SelectOnePost(handle.DB, postid)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	err = database.DeletePost(handle.DB, postid, cookie) //?
	if err != nil {
		if err := t.Execute(w, post); err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		return
	}

	http.Redirect(w, r, "/", 301)
}
