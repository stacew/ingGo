package myapp

import (
	"log"
	"net/http"
	"stacew/teamgoing/sign"
	socketmaker "stacew/teamgoing/socketMaker"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"

	socketio "github.com/googollee/go-socket.io"
)

//AppHandler is
type AppHandler struct {
	http.Handler   //embeded is-a같은 has-a 관계라는데, 이름 정해주면 안 됨...
	socketioServer *socketio.Server
	// dmHandler    dataModel.DataHandlerInterface
}

//Close is
func (m *AppHandler) Close() {
	m.socketioServer.Close()
}

//Start is
func (m *AppHandler) Start() {
	go m.socketioServer.Serve()
}

//MakeNewHandler is
func MakeNewHandler() *AppHandler {
	socketioServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatalln(err)
	}
	socketmaker.NewGameServerAndMakeSocketHandler(socketioServer)
	// -----------------
	neg := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(sign.CheckSign),
		negroni.NewStatic(http.Dir("a0001"))) //패치 후, 폴더 이름 변경
	// -----------------
	appHandler := &AppHandler{
		Handler: neg,
		// dmHandler: dataModel.NewDataHandler(dbConn),
		socketioServer: socketioServer,
	}
	// -----------------
	mux := pat.New()
	neg.UseHandler(mux)
	// -----------------
	mux.Add("GET", "/socket.io/", socketioServer)
	// mux.Add("POST", "/socket.io/", socketioServer) 예제에 post도 등록하는데 이유가..?
	// -----------------
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
	})
	// -----------------
	sign.SetHandle(mux)
	//...expand
	// -----------------
	return appHandler
}
