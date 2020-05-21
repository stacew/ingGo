package ssocketmaker

import (
	"errors"
	"gogoma/teamgoing/smygame"
	"log"
	"strconv"

	socketio "github.com/googollee/go-socket.io"
)

//NewGameServerAndMakeSocketHandler is
func (m *SocketServerInfo) NewGameServerAndMakeSocketHandler(socketioServer *socketio.Server) {
	nsp := "/"
	m.conCount = 0
	m.currentCon = 0
	m.nsp = nsp
	gameServer := smygame.NewGameServer(socketioServer, m.nsp)
	m.makeBaseSocket(socketioServer, gameServer)
	m.makeRoomSocket(socketioServer, gameServer)
	m.makeCustomSocket(socketioServer, gameServer)
}

//SocketServerInfo is
type SocketServerInfo struct {
	nsp        string
	conCount   int
	currentCon int
}

// Undevelop is 서버 정보 확인용으로 추후 작업, SocketServerInfo
func (m *SocketServerInfo) Undevelop() {

}

func (m *SocketServerInfo) makeCustomSocket(socketioServer *socketio.Server, gameServer *smygame.MyGameServer) {
	socketioServer.OnEvent(m.nsp, "cShot", func(c socketio.Conn, msg string) {
		gameServer.CShot(c.ID(), msg)
	})

	socketioServer.OnEvent(m.nsp, "cCurrentCon", func(s socketio.Conn) string {
		//expand msg.
		return strconv.Itoa(m.currentCon)
	})
}

func (m *SocketServerInfo) makeRoomSocket(socketioServer *socketio.Server, gameServer *smygame.MyGameServer) {
	socketioServer.OnEvent(m.nsp, "cJoin", func(c socketio.Conn) {
		if c == nil {
			log.Println("[ConnNil] cJoin")
			return
		}

		roomName, err := gameServer.CJoin(c.ID())
		if err != nil {
			log.Println("[Check Error] Join")
			return
		}
		c.Join(roomName)

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

func (m *SocketServerInfo) makeBaseSocket(socketioServer *socketio.Server, gameServer *smygame.MyGameServer) {
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

		gameRoomName, bExist := gameServer.CLeave(c.ID())
		if bExist {
			c.Leave(gameRoomName) //
		}
		m.currentCon--
	})
}
