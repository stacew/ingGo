package sign

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
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
	fmt.Println(val)
	//목표
	//1. 유저가 Sign Up 하면, 내 도메인 브라우저 쿠키가 세션을 만들어야 한다.

	//2. 유저가 사이트에서 활동할 때, 브라우저에 sessionID가 있다면 w에 세팅해서 페이지의 모양을 좀 바꿔줘야 할 것 같다.

	//3, 서버의 쿠키스토어와 다르면?
	//고로 CheckSign()은 2번에서 필요하다.
	//*. slither.io는 어떻게 변경시킨건지 모르겠다. 화면만 바뀌는건가?

	next(w, r)
}

//SetHandle is
func SetHandle(mux *pat.Router) {
	mux.Get("/sign", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/sign.html", http.StatusTemporaryRedirect)
	})

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
