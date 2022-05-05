package server

import (
	"database/sql"
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//CommentsHandler ..
func (handle *Handle) CommentsHandler(w http.ResponseWriter, r *http.Request) {
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
	comment := r.FormValue("comment")
	comment_btn := r.FormValue("comment_btn")
	for v := range r.Form {
		if v != "comment" && v != "comment_btn" {
			w.WriteHeader(400)
			return
		}
	}
	if comment_btn != "sent" {
		w.WriteHeader(400)
		return
	}
	if err := database.InsertCommentToDB(handle.DB, comment, postid, cookie); err != nil {
		handle.Post.ErrorVal.Err = true
		handle.Post.ErrorVal.MSG = err.Error()
	}
	http.Redirect(w, r, prevURL, 301)
}

//RateComment ..
func (handle *Handle) RateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	prevURL := r.Header.Get("Referer")

	var commentID int
	like := r.FormValue("likecom")
	dislike := r.FormValue("dislikecom")

	var counter int
	for v := range r.Form {
		counter++
		if v != "likecom" && v != "dislikecom" {
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

	if len(like) != 0 {
		commentID, _ = strconv.Atoi(like)
	} else if len(dislike) != 0 {
		commentID, _ = strconv.Atoi(dislike)
	}
	if commentID != 0 {
		database.InsertCommentNotification(handle.DB, like, dislike, commentID, cookie)
	}
	http.Redirect(w, r, prevURL, 301)
}

//EditComment ..
func (handle *Handle) EditComment(w http.ResponseWriter, r *http.Request) {
	log.Println("Got to EditComment")
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	if r.URL.Path != "/editcomment" {
		NotFound(w)
		w.WriteHeader(404)
		log.Println(404)
		return
	}

	t, err := template.ParseFiles("./templates/html/editcomment.html")
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

	commentId := r.FormValue("editcomment")
	commentid, err := strconv.Atoi(commentId)
	if err != nil {
		log.Println(404)
		w.WriteHeader(404)
		return
	}

	comments, err := database.SelectUserCommentsByComment(handle.DB, commentid, cookie)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	editComments := forum.EditComments{
		Comments: comments,
		ErrorVal: forum.Error{},
	}
	if err := t.Execute(w, editComments); err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
}

//ConfirmEditComment ..
func (handle *Handle) ConfirmEditComment(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/editcomment.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}

	cookie, _ := r.Cookie("session")

	commentId := r.FormValue("confirm")
	back := r.FormValue("back")

	if len(back) != 0 {
		postid := GetPostIdByCommentId(w, handle.DB, back)
		http.Redirect(w, r, "/post/"+postid, 301)
		return
	}

	commentid, err := strconv.Atoi(commentId)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	comment := forum.Comment{}

	for i, v := range r.Form {
		switch i {
		case "description":
			comment.Text = v[0]
		case "back":
			continue
		case "confirm":
			continue
		default:
			log.Println("An incoming value is not a required key")
			w.WriteHeader(400)
			return
		}
	}

	postid := GetPostIdByCommentId(w, handle.DB, commentId)

	updateComment := forum.UpdateComment{
		Comments: []forum.Comment{},
		Comment: forum.Comment{
			ID:   commentid,
			Text: comment.Text,
		},
		ErrorVal: forum.Error{},
		PostId:   postid,
	}

	if err = database.UpdateComment(handle.DB, &updateComment, cookie); err != nil {
		updateComment.ErrorVal.Err = true
		updateComment.ErrorVal.MSG = err.Error()
		w.WriteHeader(400)

		if err := t.Execute(w, updateComment); err != nil {
			// w.WriteHeader(500)
		}

		return
	}

	http.Redirect(w, r, "/post/"+postid, 301)
}

// DeleteComment ..
func (handle *Handle) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}

	if r.URL.Path != "/deletecomment" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}

	t, err := template.ParseFiles("./templates/html/deletecomment.html")
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

	commentId := r.FormValue("deletecomment")
	commentid, err := strconv.Atoi(commentId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}

	comment, err := database.SelectOneUserCommentsByComment(handle.DB, commentid, cookie)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if err := t.Execute(w, comment); err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
}

//ConfirmDeletePost ..
func (handle *Handle) ConfirmDeleteComment(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/deletecomment.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	cookie, _ := r.Cookie("session")

	commentId := r.FormValue("confirm")
	back := r.FormValue("back")

	if len(back) != 0 {
		postid := GetPostIdByCommentId(w, handle.DB, back)
		http.Redirect(w, r, "/post/"+postid, 301)
		return
	}

	commentid, err := strconv.Atoi(commentId)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	comment, err := database.SelectOneComment(handle.DB, commentid, cookie)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	postid := GetPostIdByCommentId(w, handle.DB, commentId)

	err = database.DeleteComment(handle.DB, commentid, cookie) //?
	if err != nil {
		if err := t.Execute(w, comment); err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		return
	}

	http.Redirect(w, r, "/post/"+postid, 301)
}

//GetPostIdByCommentId ..
func GetPostIdByCommentId(w http.ResponseWriter, db *sql.DB, data string) string {
	var postid string
	commentid, err := strconv.Atoi(data)
	if err != nil {
		w.WriteHeader(404)
		return postid
	}
	postId, err := database.SelectPostByComment(db, commentid)
	if err != nil {
		log.Println(err)
		return postid
	}
	postid = strconv.Itoa(postId)
	return postid
}
