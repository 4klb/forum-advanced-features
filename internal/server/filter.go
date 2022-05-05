package server

import (
	forum "forum/internal"
	"forum/internal/database"
	"html/template"
	"log"
	"net/http"
)

//Filter ..
func (handle *Handle) Filter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filter" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}
	switch r.Method {
	case "POST":
		tag := r.FormValue("filter")

		t, err := template.ParseFiles("./templates/html/filter.html")
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

		posts, err := database.SelectPostsByCategories(handle.DB, tag)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		HomeInfo := forum.Homepage{
			Posts:    posts,
			Category: tag,
		}
		cookie, _ := r.Cookie("session")
		if !database.IsUserInSession(handle.DB, cookie) {
			if err := t.Execute(w, HomeInfo); err != nil {
				w.WriteHeader(500)
				log.Println(err)
				return
			}
			return
		}

		HomeInfo.InSession = true
		if err := t.Execute(w, HomeInfo); err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}

	default:
		w.WriteHeader(405)
		return
	}
}
