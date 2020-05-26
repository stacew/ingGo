package sign

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
)

//uuid로 SESSION_KEY 생성 후 환경변수 등록
var cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

//CheckSign is
func CheckSign(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fGetSessionID := func(r *http.Request) string {
		session, _ := cookieStore.Get(r, "session")
		val := session.Values["id"]
		if val == nil {
			return ""
		}
		return val.(string)
	}

	val := fGetSessionID(r)
	if val != "" {
		w.Header().Add("sign", "true")

		if strings.Contains(r.URL.Path, "/sign") || strings.Contains(r.URL.Path, "/auth") {
			http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
		}
	}

	next(w, r)
}

//SetHandle is
func SetHandle(mux *pat.Router) {
	setGoogleHandle(mux)
}

//////////
func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	expiration := time.Now().Add(1 * 24 * time.Hour)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}
