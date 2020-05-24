package main

import (
	"log"
	"net/http"
	myapp "stacew/teamgoing/myApp"
	"stacew/teamgoing/port"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.Println(runtime.GOMAXPROCS(1))
	// log.Println(runtime.GOMAXPROCS(1))
	appHandler := myapp.MakeNewHandler()
	defer appHandler.Close()
	appHandler.Start()

	port := port.GetPort()
	log.Println("**************** ListenAndServe():" + port)
	log.Fatal(http.ListenAndServe(":"+port, appHandler)) //http start
}

//-ser //////////////////////////////
//NewServer(c *engineio.Options) (*Server, error)
//Close() error
//OnConnect(nsp string, f func(Conn) error)
//OnDisconnect(nsp string, f func(Conn, string))
//OnError(nsp string, f func(Conn, error))
//OnEvent(nsp, event string, f interface{})
//nsp에 사용자 이벤트 함수 등록
//ServeHTTP(w http.ResponseWriter, r *http.Request)
//???

//-ser.broad //////////////////////////////
//JoinRoom(namespace, room string, connection Conn) bool
//nsp에 방과 연결 추가
//LeaveRoom(namespace, room string, connection Conn) bool
//nsp에 방과 연결 제거
//LeaveAllRooms(namespace string, connection Conn) bool
//모든 방 및 연결 제거
//ClearRoom(namespace, room string) bool
//해당 nsp에 방 정리
//BroadcastToRoom(namespace, room, event string, args ...interface{}) bool
//해당 nsp에 방의 모두에게 이벤트와 인자를
//RoomLen(namespace, room string) int
//해당 nsp 방의 con 수
//Rooms(namespace string) []string
//해당 nsp에 모든 방 목록
//ForEach(namespace, room string, f EachFunc) bool
//방의 모두에게 특정 함수를 실행

//--con //////////////////////////////
//Close() error
//특정 con 연결 해제
//Emit(msg string, v ...interface{})
//특정 con 메시지

//Context() interface{}
//SetContext(v interface{})
//con 확장

//con 강제퇴장..?..?..
//Join(room string)
//Leave(room string)
//LeaveAll()
//Rooms() []string
