package server

import (
	"html/template"
	"log"
	"net/http"
)

//NotFound ..
func NotFound(w http.ResponseWriter) {
	n, err := template.ParseFiles("./templates/html/pagenotfound.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	if err = n.Execute(w, nil); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
}
