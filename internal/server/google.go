package server

import (
	"encoding/json"
	forum "forum/internal"
	"forum/internal/database"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var clientid = "637515183608-fpugh5na18o1f5mmrb29e0pqfmbf65mt.apps.googleusercontent.com"
var clientsecret = "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"
var scopes = []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}

//GogleOuath ..
func (handle *Handle) GogleOuath(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/google" {
		NotFound(w)
		w.WriteHeader(404)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	URL, err := url.Parse("https://accounts.google.com/o/oauth2/auth")
	if err != nil {
		return
	}
	parameters := url.Values{}
	parameters.Add("client_id", clientid)
	parameters.Add("scope", strings.Join(scopes, " "))
	parameters.Add("redirect_uri", "http://localhost:8080/login/google/callback")
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//GoogleCallbackHandler ..
func (handle *Handle) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	permissoncode := r.FormValue("code")
	if permissoncode == "" || len(permissoncode) != 73 {
		http.Error(w, "Can not sign in with GOOGLE", http.StatusUnauthorized)
		return
	}

	tokenURL := "https://oauth2.googleapis.com/token"

	accesstoken, err := exchange(permissoncode, tokenURL)
	if err != nil {
		http.Error(w, "Can not sign in with GOOGLE", http.StatusUnauthorized)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accesstoken))
	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Can not read from response from GOOGLE", http.StatusBadRequest)
		return
	}

	user := getInfo(string(response))
	user.Confirm = "GOOGLE"
	id := uuid.NewV4()
	cookie := &http.Cookie{
		Name:    "session",
		Value:   id.String(),
		Expires: time.Now().Add(60 * time.Minute),
		Path:    "/",
		MaxAge:  3600,
	}
	if user.Login == "" || user.Name == "" {
		http.Error(w, "Can not sign in with GOOGLE", http.StatusUnauthorized)
		return
	}
	if err := database.InsertUserBySocialNetworks(handle.DB, user); err != nil {
		if err = database.CanLoginBySocialNetworks(handle.DB, user); err != nil {
			http.Redirect(w, r, "/login", 301)
			log.Println("Can not sign in by GOOGLE")
			return
		}
	}
	http.SetCookie(w, cookie)

	database.InsertCookieIntoDB(handle.DB, user.Login, cookie)
	http.Redirect(w, r, "/", 301)
}

func getInfo(info string) forum.User {
	user := forum.User{}
	info = info[2 : len(info)-3]
	checkinfo := strings.Split(info, "\n")
	for _, v := range checkinfo {
		v = strings.Trim(v, " ")
		ans := strings.Split(v, ":")
		ans[0] = strings.Trim(ans[0], `"`)
		if ans[0] == "email" {
			user.Login = ans[1][2 : len(ans[1])-2]
		}
		if ans[0] == "given_name" {
			user.Name = ans[1][2 : len(ans[1])-2]
		}
	}
	return user
}

func exchange(permissoncode, tokenurl string) (string, error) {
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {permissoncode},
	}
	v.Set("redirect_uri", "http://localhost:8080/login/google/callback")

	return retrieveToken(v, tokenurl)
}

func retrieveToken(v url.Values, tokenurl string) (string, error) {
	req, err := newTokenRequest(tokenurl, v)
	if err != nil {
		return "", err
	}

	accesstoken, err := doTokenRoundTrip(req)
	if err != nil {
		return "", err
	}
	return accesstoken, nil

}

func doTokenRoundTrip(req *http.Request) (string, error) {
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, resperr := ioutil.ReadAll(resp.Body)
	if resperr != nil {
		return "", resperr
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
		return ghresp.AccessToken, err
	}

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken, nil
}

func newTokenRequest(tokenURL string, v url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(url.QueryEscape(clientid), url.QueryEscape(clientsecret))
	return req, nil
}
