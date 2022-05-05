package server

import (
	forum "forum/internal"
	"log"
	"net/http"
)

//PostFormChecker ..
func PostFormChecker(w http.ResponseWriter, r *http.Request, categories []string) (forum.CreatePost, bool) {
	post := forum.CreatePost{}

	for i, v := range r.Form {
		switch i {
		case "title":
			post.Title = v[0]
		case "description":
			post.Text = v[0]
		case "confirm":
			continue
		case "categories":
			for _, category := range v {
				post.Categories = append(post.Categories, category)
			}
		default:
			log.Println("An incoming value is not a required key")
			w.WriteHeader(400)
			return post, false
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
			return post, false
		}
	}

	for i := 0; i < len(post.Categories)-1; i++ {
		for j := i + 1; j < len(post.Categories); j++ {
			if post.Categories[i] == post.Categories[j] {
				log.Println("Categories are not unique")
				w.WriteHeader(400)
				return post, false
			}
		}
	}
	return post, true
}
