package main

import (
	"forum/internal/database"
	"forum/internal/server"
	"log"
	"net/http"
	// "net/mux"
)

func main() {
	db, err := database.DbInit()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	create := make(chan int)

	go database.RemoveExpiredCookie(db, create)

	go database.CreateDb(db, create)

	mux := http.NewServeMux()

	handle := server.CreateHandle(db)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("templates/assets/"))))
	mux.HandleFunc("/", handle.HomePage)

	mux.HandleFunc("/login", handle.GetSignIn)
	mux.HandleFunc("/postlogin", handle.PostSignIn)

	mux.HandleFunc("/signout", handle.GetSignOut)
	mux.HandleFunc("/postsignout", handle.PostSignOut)

	mux.HandleFunc("/registration", handle.GetSignUp)
	mux.HandleFunc("/postregistration", handle.PostSignUp)

	mux.HandleFunc("/post/", handle.ChosenPost)
	mux.HandleFunc("/createpost", handle.GetCreatePost)
	mux.HandleFunc("/postcreate", handle.PostCreatePost)
	mux.HandleFunc("/ratepost", handle.RatePost)

	mux.HandleFunc("/ratecomment", handle.RateComment)
	mux.HandleFunc("/placecomment", handle.CommentsHandler)

	mux.HandleFunc("/editpost", handle.EditPost)
	mux.HandleFunc("/confirmedit", handle.ConfirmEditPost)
	mux.HandleFunc("/deletepost", handle.DeletePost)
	mux.HandleFunc("/confirmdelete", handle.ConfirmDeletePost)

	mux.HandleFunc("/editcomment", handle.EditComment)
	mux.HandleFunc("/confirmeditcomment", handle.ConfirmEditComment)
	mux.HandleFunc("/deletecomment", handle.DeleteComment)
	mux.HandleFunc("/confirmdeletecomment", handle.ConfirmDeleteComment)

	mux.HandleFunc("/profile", handle.Profile)
	mux.HandleFunc("/filter", handle.Filter)

	mux.HandleFunc("/login/github", handle.GitHubOauth)
	mux.HandleFunc("/login/github/callback", handle.GithubCallbackHandler)

	mux.HandleFunc("/login/google", handle.GogleOuath)
	mux.HandleFunc("/login/google/callback", handle.GoogleCallbackHandler)

	log.Println("Server is Listening..." + "\n" + "http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
