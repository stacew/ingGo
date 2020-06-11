package gameserver

import (
	"errors"
	gameroom "stacew/teamgo/gameRoom"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

const (
	nRoomCapacity = 4 //4, 6, 8
)

//ClientShot is Client Shot Input
func (m *MyGameServer) ClientShot(conID string, msg string) error {
	gameRoomInfo, ok := m.conGameRoomMap[conID]
	if !ok {
		return errors.New("[Check Error] No Room User Send CShot")
	}

	gameRoomInfo.ClientShot(conID, msg)
	return nil
}

//BroadCastJoinAndStart is
func (m *MyGameServer) BroadCastJoinAndStart(conID string) {

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	gameRoomInfo, _ := m.conGameRoomMap[conID]
	if gameRoomInfo.GetPlayerCount() < nRoomCapacity {
		gameRoomInfo.BroadcastInfoRoom(conID)
		return
	}

	go gameRoomInfo.Start()
}

func (m *MyGameServer) joinGameRoom(conID string, gameRoomInfo *gameroom.RoomInfo) string {
	roomName := gameRoomInfo.Join(conID)
	m.conGameRoomMap[conID] = gameRoomInfo
	return roomName
}

//ClientJoin is
func (m *MyGameServer) ClientJoin(conID string) (string, error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	//추후 레이팅 매칭 작업
	for gameRoomInfo, gameRoomName := range m.matchRoom {
		m.joinGameRoom(conID, gameRoomInfo)
		roomJoinCount := gameRoomInfo.GetPlayerCount()
		if roomJoinCount == nRoomCapacity {
			m.playRoom[gameRoomInfo] = gameRoomName
			delete(m.matchRoom, gameRoomInfo)
		}
		return gameRoomName, nil
	}

	//matchRoom이 없으면 새로 만든다.
	gameRoomInfo := gameroom.NewRoomInfo(m.socketioServer, m.nsp, nRoomCapacity, m.nMakedRoomCount)
	gameRoomName := m.joinGameRoom(conID, gameRoomInfo)
	m.matchRoom[gameRoomInfo] = gameRoomName
	return gameRoomName, nil
}

//ClientLeave is
func (m *MyGameServer) ClientLeave(conID string) (string, bool) {

	m.mutex.RLock()
	gameRoomInfo, ok := m.conGameRoomMap[conID]
	m.mutex.RUnlock()
	if !ok {
		return "", false //not joined game room
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	nRemainCount, roomName := gameRoomInfo.Leave(conID)
	if nRemainCount == 0 {
		delete(m.customRoom, gameRoomInfo)
		delete(m.freeRoom, gameRoomInfo)
		delete(m.matchRoom, gameRoomInfo)
		delete(m.playRoom, gameRoomInfo)
	} else {
		gameRoomInfo.BroadcastInfoRoom(conID)
	}

	delete(m.conGameRoomMap, conID)
	return roomName, true
}

//NewGameServer is
func NewGameServer(socketioServer *socketio.Server, nsp string) *MyGameServer {
	gameServer := &MyGameServer{
		socketioServer: socketioServer,
		nsp:            nsp,
		mutex:          new(sync.RWMutex),

		nMakedRoomCount: 0,
		matchRoom:       make(map[*gameroom.RoomInfo]string),
		freeRoom:        make(map[*gameroom.RoomInfo]string),
		customRoom:      make(map[*gameroom.RoomInfo]string),
		playRoom:        make(map[*gameroom.RoomInfo]string),
		conGameRoomMap:  make(map[string]*gameroom.RoomInfo)}

	return gameServer
}

//MyGameServer is
type MyGameServer struct {
	socketioServer *socketio.Server
	nsp            string
	mutex          *sync.RWMutex

	nMakedRoomCount int                           // use matchRoom, freeRoom GameRoomName
	matchRoom       map[*gameroom.RoomInfo]string // rating game (login essential)
	freeRoom        map[*gameroom.RoomInfo]string // free match room
	customRoom      map[*gameroom.RoomInfo]string // make and link
	playRoom        map[*gameroom.RoomInfo]string // playing room
	conGameRoomMap  map[string]*gameroom.RoomInfo
}
