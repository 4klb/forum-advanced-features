package server

import (
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//GetSignUp  ..
func (handle *Handle) GetSignUp(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/html/registration.html")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	switch r.Method {
	case "GET":
		if r.URL.Path != "/registration" {
			NotFound(w)
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

//PostSignUp ..
func (handle *Handle) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	if r.URL.Path != "/postregistration" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}

	t, err := template.ParseFiles("./templates/html/registration.html")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if err := r.ParseForm(); err != nil {
		return
	}
	if len(r.Form) != 4 {
		w.WriteHeader(400)
		return
	}
	for v := range r.Form {
		if v != "name" && v != "login" && v != "password" && v != "passwordconfirm" {
			w.WriteHeader(400)
			return
		}
	}
	input := forum.User{
		Name:     r.FormValue("name"),
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
		Confirm:  r.FormValue("passwordconfirm"),
	}

	er := forum.Error{}

	if err := database.CanRegister(handle.DB, input); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		er.Err = true
		er.MSG = err.Error()
		if err := t.Execute(w, er); err != nil {
			w.WriteHeader(500)
			log.Println(err)
		}
		return
	}
	log.Printf("User %v has added to DB\n", input.Name)
	http.Redirect(w, r, "/", 301)
}
