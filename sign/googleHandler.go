package sign

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//GoogleUserID is
type GoogleUserID struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `josn:"verified_email"`
	Picture       string `json:"picture"`
}

const (
	authGoogleLogin       = "/auth/google/login"
	authGoogleCallbackURL = "/auth/google/callback"
	oauthGoogleAPIURL     = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

var googleOauthConfig = oauth2.Config{
	RedirectURL: "http://localhost:8080" + authGoogleCallbackURL,
	// RedirectURL:  os.Getenv("DOMAIN_NAME") + authGoogleCallbackURL,
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	//요청 정보?
	Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

func getGoogleUserInfo(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code) //context : thread safe한 저장소
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}
	resp, err := http.Get(oauthGoogleAPIURL + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s", err.Error())
	}
	return ioutil.ReadAll(resp.Body)
}

//////
func googleOauthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")
	if oauthstate.Value != r.FormValue("state") {
		errMsg := fmt.Sprintf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state"))
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Store ID info into Session cookie
	var userInfo GoogleUserID
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := cookieStore.Get(r, "session")
	session.Values["id"] = userInfo.ID //key value 다른 정보 저장해도 됨.

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

//////
func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

////
func setGoogleHandle(mux *pat.Router) {
	mux.HandleFunc(authGoogleLogin, googleLoginHandler)
	mux.HandleFunc(authGoogleCallbackURL, googleOauthCallback)
}
