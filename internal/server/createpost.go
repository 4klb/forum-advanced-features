package server

import (
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//GetCreatePost ...
func (handle *Handle) GetCreatePost(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/login", 301)
		return
	}
	t, err := template.ParseFiles("./templates/html/createpost.html")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	categories := []string{"cats", "IT", "cybersport", "github", "Anime", "Alem"}

	unionEnterPost := forum.UnionEnterPost{
		ErrorVal:   forum.Error{},
		Categories: categories,
	}
	switch r.Method {
	case "GET":
		if r.URL.Path != "/createpost" {
			NotFound(w)
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			log.Println(err)
			return
		}
		if err := t.Execute(w, unionEnterPost); err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	default:
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
}

//PostCreatePost ...
func (handle *Handle) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if !database.IsUserInSession(handle.DB, cookie) {
		http.Redirect(w, r, "/postlogin", 301)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}

	if err := r.ParseForm(); err != nil {
		return
	}

	t, err := template.ParseFiles("./templates/html/createpost.html")
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	post := forum.CreatePost{
		Title:      r.Form["title"][0],
		Text:       r.Form["description"][0],
		Categories: r.Form["categories"],
	}

	categories := []string{"cats", "IT", "cybersport", "github", "Anime", "Alem"}

	unionEnterPost := forum.UnionEnterPost{
		ErrorVal:   forum.Error{},
		Categories: categories,
	}

	for v := range r.Form {
		if len(post.Categories) > 0 && v != "title" && v != "description" && v != "categories" {
			w.WriteHeader(400)
			return
		}
	}

	for _, v := range r.Form["categories"] {
		count := 0
		for _, j := range categories {
			if v != j {
				count++
			}
		}
		if count == 6 {
			w.WriteHeader(400)
			return
		}
	}
	if err := database.InsertPostIntoDB(handle.DB, &post, cookie); err != nil {
		unionEnterPost.ErrorVal.Err = true
		unionEnterPost.ErrorVal.MSG = err.Error()
		w.WriteHeader(400)
		if err := t.Execute(w, unionEnterPost); err != nil {
			w.WriteHeader(500)
		}
		return
	}
	http.Redirect(w, r, "/", 301)
}
