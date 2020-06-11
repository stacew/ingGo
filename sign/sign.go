package sign

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
)

const (
	strSession    = "session"
	strOAuthState = "oauthstate"
)

//uuid로 SESSION_KEY 생성 후 환경변수 등록
const strEnvSessionKey = "SESSION_KEY"

var cookieStore = sessions.NewCookieStore([]byte(os.Getenv(strEnvSessionKey)))

const (
	//ConstPlatformID is
	ConstPlatformID = "platformID"
	//ConstPlatformType is
	ConstPlatformType = "platform"
)

//PlatformType is
type PlatformType int

const (
	//Google is
	Google PlatformType = iota
	//Facebook is
	Facebook
	//None is
	None
)

func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	expiration := time.Now().Add(1 * 24 * time.Hour)
	cookie := &http.Cookie{Name: strOAuthState, Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

//SetHandle is
func SetHandle(mux *pat.Router) {
	setGoogleHandle(mux)
}

// //CheckSign is
// func CheckSign(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
// 	//fGetSession
// 	fGetSession := func(r *http.Request) (string, string) {
// 		session, _ := cookieStore.Get(r, constSession)
// 		platformID := session.Values[ConstPlatformID]
// 		platformType := session.Values[ConstPlatformType]
// 		if platformID == nil || platformType == nil {
// 			return "", ""
// 		}
// 		return platformID.(string), platformType.(string)
// 	}

// 	platformID, platformType := fGetSession(r)
// 	if platformID != "" && platformType != "" {
// 		// if strings.Contains(r.URL.Path, "/sign") || strings.Contains(r.URL.Path, "/auth") {
// 		if strings.Contains(r.URL.Path, "/auth") {
// 			http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
// 		}
// 	}

// 	next(w, r)
// }
