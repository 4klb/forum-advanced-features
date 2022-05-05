package server

import (
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//GetSignOut ..
func (handle *Handle) GetSignOut(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	switch r.Method {
	case "GET":
		if r.URL.Path != "/signout" {
			w.WriteHeader(404)
			return
		}
		t, err := template.ParseFiles("./templates/html/signout.html")
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}
		if err := t.Execute(w, nil); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}
	default:
		w.WriteHeader(405)
		return
	}
}

//PostSignOut ..
func (handle *Handle) PostSignOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/postsignout" {
		w.WriteHeader(404)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	if err := r.ParseForm(); err != nil {
		return
	}
	if len(r.Form) != 1 {
		w.WriteHeader(400)
		return
	}
	choice := r.FormValue("choice")
	if choice == "yes" {
		database.DeleteCookieFromDB(handle.DB, cookie)
		log.Println("Delele cookie from DB happened successfully ")
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 301)
	} else if choice == "no" {
		http.Redirect(w, r, "/", 301)
		return
	} else {
		w.WriteHeader(400)
		log.Println("Do not change values!!!")
		return
	}
}
