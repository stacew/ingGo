package smygame

import (
	"errors"

	socketio "github.com/googollee/go-socket.io"
)

const (
	nRoomCapacity = 4 //2, 4, 6, 8
)

//CShot is Client Shot Input
func (m *MyGameServer) CShot(conID string, msg string) error {
	gameRoomInfo, ok := m.conGameRoomMap[conID]
	if !ok {
		return errors.New("[Check Error] No Room User Send CShot")
	}

	gameRoomInfo.CShot(conID, msg)
	return nil
}

//BroadCastJoinAndStart is
func (m *MyGameServer) BroadCastJoinAndStart(conID string) {
	gameRoomInfo, _ := m.conGameRoomMap[conID]
	if gameRoomInfo.GetPlayerCount() < nRoomCapacity {
		gameRoomInfo.BroadcastInfoRoom()
		return
	}

	go gameRoomInfo.Start()
}

func (m *MyGameServer) joinProcess(conID string, gameRoomInfo *GameRoomInfo) string {
	roomName := gameRoomInfo.Join(conID)
	m.conGameRoomMap[conID] = gameRoomInfo
	return roomName
}

//CJoin is
func (m *MyGameServer) CJoin(conID string) (string, error) {
	//추후 레이팅 매칭 작업
	for gameRoomInfo, gameRoomName := range m.matchRoom {
		m.joinProcess(conID, gameRoomInfo)
		roomJoinCount := gameRoomInfo.GetPlayerCount()
		if roomJoinCount == nRoomCapacity {
			m.playRoom[gameRoomInfo] = gameRoomName
			delete(m.matchRoom, gameRoomInfo)
		}
		return gameRoomName, nil
	}

	//matchRoom이 없으면 새로 만든다.
	gameRoomInfo := NewGameRoomInfo(m.socketioServer, m.nsp, nRoomCapacity)
	gameRoomName := m.joinProcess(conID, gameRoomInfo)
	m.matchRoom[gameRoomInfo] = gameRoomName

	return gameRoomName, nil
}

//CLeave is
func (m *MyGameServer) CLeave(conID string) (string, bool) {
	gameRoomInfo, ok := m.conGameRoomMap[conID]
	if !ok {
		return "", false //not joined game room
	}

	nRemainCount, roomName := gameRoomInfo.Leave(conID)
	if nRemainCount == 0 {
		delete(m.matchRoom, gameRoomInfo)
		delete(m.playRoom, gameRoomInfo)
	} else {
		gameRoomInfo.BroadcastInfoRoom()
	}

	delete(m.conGameRoomMap, conID)
	return roomName, true
}

//NewGameServer is
func NewGameServer(socketioServer *socketio.Server, nsp string) *MyGameServer {
	gameserver := &MyGameServer{
		socketioServer: socketioServer,
		nsp:            nsp,

		matchRoom:      make(map[*GameRoomInfo]string),
		playRoom:       make(map[*GameRoomInfo]string),
		conGameRoomMap: make(map[string]*GameRoomInfo)}

	return gameserver
}

//MyGameServer is
type MyGameServer struct {
	socketioServer *socketio.Server
	nsp            string

	matchRoom      map[*GameRoomInfo]string
	playRoom       map[*GameRoomInfo]string
	conGameRoomMap map[string]*GameRoomInfo
}
