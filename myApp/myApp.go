package myapp

import (
	"log"
	"net/http"
	"stacew/teamgo/cipher"
	dm "stacew/teamgo/dataModel"
	"stacew/teamgo/sign"
	socketserver "stacew/teamgo/socketServer"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

//AppHandler is
type AppHandler struct {
	http.Handler   //embeded is-a같은 has-a 관계라는데, 이름 정해주면 안 됨...
	socketioServer *socketio.Server
	dmHandler      dm.DataHandlerInterface
}

func (m *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {

	cookie...
	encryptplatformID := w.Header().Get(sign.ConstPlatformID)
	platformType := w.Header().Get(sign.ConstPlatformType)
	if encryptplatformID != "" && platformType != "" {
		platformID := cipher.Decrypt(encryptplatformID)
		signPlatformType, _ := strconv.Atoi(platformType)
		userInfo := m.dmHandler.GetAndAddUserInfo(sign.PlatformType(signPlatformType), platformID)
		log.Println(userInfo.Name)
	}

	http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
}

//Close is
func (m *AppHandler) Close() {
	m.socketioServer.Close()
	m.dmHandler.Close()
}

//Start is
func (m *AppHandler) Start() {
	go m.socketioServer.Serve()
}

//MakeNewHandler is
func MakeNewHandler(dbConn string) *AppHandler {
	socketioServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	socketserver.NewGameServerAndMakeSocketHandler(socketioServer)
	// -----------------
	neg := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		// negroni.HandlerFunc(sign.CheckSign), 필요 없어짐... request header cookie 정보 이용하도록..
		negroni.NewStatic(http.Dir("a0001"))) //패치 후, 폴더 이름 변경
	// -----------------
	appHandler := &AppHandler{
		Handler:        neg,
		dmHandler:      dm.NewDataHandler(dbConn),
		socketioServer: socketioServer,
	}
	// -----------------
	mux := pat.New()
	neg.UseHandler(mux)
	// -----------------
	mux.Add("GET", "/socket.io/", socketioServer)
	//mux.Add("POST", "/socket.io/", socketioServer)//iris? gin 프레임 워크 예제에 post도 등록하던데 지금 여기서는 의미를 모르겠음
	// -----------------
	sign.SetHandle(mux)
	// -----------------
	//...expand...
	mux.Get("/", appHandler.indexHandler)
	// -----------------
	return appHandler
}
