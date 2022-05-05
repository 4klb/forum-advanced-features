package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	forum "forum/internal"
	"forum/internal/database"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

var githubclientid string = "0cf3f46b1d9a1744789c"
var githubclientsecret string = "8efdb5e258377ae80ddf7a57c613fce7c9d89248"

//GitHubOauth ..
func (handle *Handle) GitHubOauth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/github" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubclientid,
		"http://localhost:8080/login/github/callback",
	)

	http.Redirect(w, r, redirectURL, 301)
}

//GithubCallbackHandler ..
func (handle *Handle) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Can not sign in with GITHUBS", http.StatusUnauthorized)
		return
	}

	githubAccessToken := getGithubAccessToken(code)
	if githubAccessToken == "" {
		http.Error(w, "Can not sign in with GITHUB", http.StatusUnauthorized)
		return
	}

	githubData := getGithubData(githubAccessToken)
	handle.loggedinHandler(w, r, githubData, githubAccessToken)
}

func (handle *Handle) loggedinHandler(w http.ResponseWriter, r *http.Request, githubData, token string) {
	if githubData == "" || len(token) != 40 {
		// Unauthorized users get an unauthorized message
		http.Error(w, "Can not sign in with GITHUB", http.StatusUnauthorized)
		return
	}
	// Set return type JSON
	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		w.WriteHeader(400)
		return
	}
	var user = forum.User{
		Confirm: "GITHUB",
	}
	if err := json.Unmarshal(prettyJSON.Bytes(), &user); err != nil {
		log.Println(err)
		return
	}
	if user.Name == "" {
		user.Name = user.Login
	}

	id := uuid.NewV4()
	cookie := &http.Cookie{
		Name:    "session",
		Value:   id.String(),
		Expires: time.Now().Add(60 * time.Minute),
		Path:    "/",
		MaxAge:  3600,
	}

	if user.Login == "" || user.Name == "" {
		http.Error(w, "Can not sign in with GITHUB", http.StatusUnauthorized)
		return
	}

	if err := database.InsertUserBySocialNetworks(handle.DB, user); err != nil {
		if err = database.CanLoginBySocialNetworks(handle.DB, user); err != nil {
			http.Redirect(w, r, "/login", 301)
			log.Println("Can not sign in by GITHUB")
			return
		}
	}
	http.SetCookie(w, cookie)

	database.InsertCookieIntoDB(handle.DB, user.Login, cookie)
	http.Redirect(w, r, "/", 301)
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		return ""
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return ""
	}

	// Read the response as a byte slice
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	// Convert byte slice to string and return
	return string(respbody)
}

func getGithubAccessToken(code string) string {

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     githubclientid,
		"client_secret": githubclientsecret,
		"code":          code,
	}

	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return ""
	}

	// Response body converted to stringified JSON
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse

	if err := json.Unmarshal(respbody, &ghresp); err != nil {
		return ""
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}
