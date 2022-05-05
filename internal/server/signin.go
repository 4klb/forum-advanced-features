package server

import (
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

//PostSignIn ..
func (handle *Handle) PostSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	if r.URL.Path != "/postlogin" {
		w.WriteHeader(404)
		return
	}
	er := forum.Error{}
	t, err := template.ParseFiles("./templates/html/login.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	if err := r.ParseForm(); err != nil {
		return
	}
	prevURL := r.Header.Get("Referer")
	if strings.Contains(prevURL, "postlogin") {
		if len(r.Form) != 2 {
			w.WriteHeader(400)
			return
		}
		for v := range r.Form {
			if v != "login" && v != "password" {
				w.WriteHeader(400)
				return
			}
		}
	}

	input := forum.User{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	}

	id := uuid.NewV4()
	cookie := &http.Cookie{
		Name:    "session",
		Value:   id.String(),
		Expires: time.Now().Add(60 * time.Minute),
		MaxAge:  3600,
	}

	if err = database.CanLogin(handle.DB, input); err == nil {
		log.Printf("User %v has added to session \n", input.Login)
		database.InsertCookieIntoDB(handle.DB, input.Login, cookie)
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", 301)
		return
	}

	w.WriteHeader(400)
	log.Println(err)
	er.Err = true
	er.MSG = err.Error()
	if err := t.Execute(w, er); err != nil {
		w.WriteHeader(500)
		return
	}
}

//GetSignIn ..
func (handle *Handle) GetSignIn(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/login.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	switch r.Method {
	case "GET":
		if r.URL.Path != "/login" {
			w.WriteHeader(404)
			return
		}
		er := forum.Error{}
		if err = t.Execute(w, er); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}
	default:
		w.WriteHeader(405)
		return
	}
}
