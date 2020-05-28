package socketserver

import (
	"errors"
	gameserver "stacew/teamgoing/gameServer"
	"sync"

	"log"

	socketio "github.com/googollee/go-socket.io"
)

//NewGameServerAndMakeSocketHandler is
func NewGameServerAndMakeSocketHandler(socketioServer *socketio.Server) {
	indexGameNSP := "/"
	indexSocektServerInfo := &socketServerInfo{
		nsp:        indexGameNSP,
		mutex:      new(sync.RWMutex),
		conCount:   0,
		currentCon: 0,
	}

	gameServer := gameserver.NewGameServer(socketioServer, indexSocektServerInfo.nsp)
	indexSocektServerInfo.makeBaseSocket(socketioServer, gameServer)
	indexSocektServerInfo.makeRoomSocket(socketioServer, gameServer)
	indexSocektServerInfo.makeCustomSocket(socketioServer, gameServer)
}

//socketServerInfo is
type socketServerInfo struct {
	nsp        string
	mutex      *sync.RWMutex
	conCount   int
	currentCon int
}

// Undevelop is 서버 정보 확인용으로 추후 작업, socketServerInfo
func (m *socketServerInfo) Undevelop() {

}

func (m *socketServerInfo) makeCustomSocket(socketioServer *socketio.Server, gameServer *gameserver.MyGameServer) {
	socketioServer.OnEvent(m.nsp, "cShot", func(c socketio.Conn, msg string) {
		gameServer.CShot(c.ID(), msg)
	})

	// socketioServer.OnEvent(m.nsp, "cReqMsg", func(s socketio.Conn) string {
	// 	//Expand cReqMsg.
	// 	return "strconv.Itoa(m.currentCon)" + "cReqMsg Msg Test"
	// })
}

func (m *socketServerInfo) makeRoomSocket(socketioServer *socketio.Server, gameServer *gameserver.MyGameServer) {
	socketioServer.OnEvent(m.nsp, "cJoin", func(c socketio.Conn) {
		if c == nil {
			log.Println("[ConnNil] cJoin")
			return
		}

		m.mutex.Lock() //

		roomName, err := gameServer.CJoin(c.ID())
		if err != nil {
			log.Println("[Check Error] Join")
			return
		}
		c.Join(roomName)

		log.Println(c.Rooms())
		m.mutex.Unlock() //

		gameServer.BroadCastJoinAndStart(c.ID())
	})

	socketioServer.OnEvent(m.nsp, "cLeave", func(c socketio.Conn) {
		if c == nil {
			log.Println("[ConnNil] cLeave")
		}

		gameRoomName, bExist := gameServer.CLeave(c.ID())
		if bExist {
			c.Leave(gameRoomName) //
		}
	})
}

func (m *socketServerInfo) makeBaseSocket(socketioServer *socketio.Server, gameServer *gameserver.MyGameServer) {
	socketioServer.OnConnect(m.nsp, func(c socketio.Conn) error {

		if c == nil {
			return errors.New("[ConnNil] connect nil : " + m.nsp)
		}

		log.Println("OnConnect", c.ID())

		m.conCount++
		m.currentCon++
		return nil
	})

	socketioServer.OnError(m.nsp, func(c socketio.Conn, e error) {
		if c == nil {
			log.Println("[ConnNil] OnError", e, m.nsp)
			return
		}
		// gameRoomName, bExist := gameServer.CLeave(c.ID())
		// if bExist {
		// 	c.Leave(gameRoomName) //
		// }
		// m.currentCon-- <- OnDisconnect랑 중복으로 깍임.
	})

	socketioServer.OnDisconnect(m.nsp, func(c socketio.Conn, reason string) {
		if c == nil {
			log.Println("[ConnNil] OnDisconnect", reason)
			return
		}

		m.mutex.Lock() //
		gameRoomName, bExist := gameServer.CLeave(c.ID())
		if bExist {
			c.Leave(gameRoomName) //
		}
		m.mutex.Unlock() //
		m.currentCon--
	})
}
