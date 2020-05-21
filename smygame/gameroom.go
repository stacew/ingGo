package smygame

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

const (
	boardX = 1000
	boardY = 1000
	nFrame = 16 //25 = 40f, 20 = 50f, 16 = 60f
)

type xyfloat64 struct {
	x float64
	y float64
}
type xy struct {
	x int
	y int
}

type shotInfo struct {
	target      xy //0 ~ 1000
	bSettedShot bool
}
type playerInfo struct {
	bwTeam  string //b is black, w is white
	pos     xyfloat64
	shot    shotInfo
	bLive   bool
	speed   xyfloat64
	nRadius int
}

//NewGameRoomInfo is
func NewGameRoomInfo(socketioServer *socketio.Server, nsp string, nRoomCapacity int) *GameRoomInfo {
	const rowCount = 3
	const colCount = 3
	const maxPosCount = rowCount * colCount //보드 확장 가능성

	gameRoomName := uuid.New().String()
	if gameRoomName == "" {
		log.Println("[Check uuid.New()]")
		gameRoomName = uuid.New().String()
	}

	gameRoomInfo := &GameRoomInfo{
		socketioServer: socketioServer,
		nsp:            nsp,
		roomCapacity:   nRoomCapacity,
		playerMap:      make(map[string]*playerInfo),
		gameRoomName:   gameRoomName,
		randPos:        make([]xy, maxPosCount, maxPosCount),
	}

	//tucker random
	for i := 0; i < rowCount; i++ {
		for j := 0; j < colCount; j++ {
			gameRoomInfo.randPos[i*rowCount+j].x = 150 + i*350
			gameRoomInfo.randPos[i*rowCount+j].y = 150 + j*350
		}
	}

	nSeed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(nSeed)
	for i := 0; i < maxPosCount; i++ {
		n1 := r.Intn(maxPosCount)
		n2 := r.Intn(maxPosCount)
		gameRoomInfo.randPos[n1], gameRoomInfo.randPos[n2] = gameRoomInfo.randPos[n2], gameRoomInfo.randPos[n1]
	}
	//tucker random

	return gameRoomInfo
}

//GameRoomInfo is
type GameRoomInfo struct {
	socketioServer *socketio.Server
	nsp            string
	gameRoomName   string

	roomCapacity int
	randPos      []xy

	playerMap        map[string]*playerInfo //key = c.ID()
	bShotTime        bool
	bAttackTeamBlack bool
}

//GetPlayerCount is
func (m *GameRoomInfo) GetPlayerCount() int {
	return len(m.playerMap)
}

//Join is
func (m *GameRoomInfo) Join(conID string) string {
	m.playerMap[conID] = &playerInfo{}
	return m.gameRoomName
}

//Leave is
func (m *GameRoomInfo) Leave(conID string) (int, string) {
	delete(m.playerMap, conID)

	nRemainCount := len(m.playerMap)
	return nRemainCount, m.gameRoomName
}

func (m *GameRoomInfo) playerStartUp() {
	i := 0
	for _, player := range m.playerMap {
		(*player).pos.x = float64(m.randPos[i].x)
		(*player).pos.y = float64(m.randPos[i].y)
		(*player).nRadius = 30
		i++
		if i%2 == 0 {
			player.bwTeam = "b"
		} else {
			player.bwTeam = "w"
		}

		player.bLive = true
		player.shot.bSettedShot = false
	}
}

/////////////////////////////////////////////

//BroadcastInfoRoom is
func (m *GameRoomInfo) BroadcastInfoRoom() {
	msg := ".i" + //.i : infoRoom
		strconv.Itoa(len(m.playerMap)) + "," +
		strconv.Itoa(m.roomCapacity) + ","

	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sGame", msg)
}
func (m *GameRoomInfo) broadcastOneShotStartEnd(start bool) {
	msg := ".o" //.o : oneShotStartEnd
	if start {
		msg = msg + "s,"
	} else {
		msg = msg + "e,"
	}
	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sGame", msg)
}
func (m *GameRoomInfo) broadcastClientTimer(nTime int) {
	msg := ".t" + //.c : Timer
		strconv.Itoa(nTime) + ","

	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sGame", msg)
}
func (m *GameRoomInfo) broadcastPlaying(starting bool) {
	var msg string
	if starting {
		msg = ".s" //.s: Starting
	}
	msg = msg + ".a" //.a : attackTeam
	if m.bAttackTeamBlack {
		msg = msg + "b,"
	} else {
		msg = msg + "w,"
	}

	for cID, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			msg = msg + ".d" + //.d : die
				cID + ","
			continue
		}
		msg = msg + ".p" + //.p : playing
			cID + "," +
			strconv.Itoa(int(playerInfo.pos.x)) + "," +
			strconv.Itoa(int(playerInfo.pos.y)) + "," +
			strconv.Itoa(playerInfo.nRadius) + "," +
			playerInfo.bwTeam + ","
	}
	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sGame", msg)
}

/////////////////////////////////////////////

//Start is goroutine
func (m *GameRoomInfo) Start() {

	m.bAttackTeamBlack = true
	m.playerStartUp()

	m.broadcastPlaying(true) //true = starting

	for {
		bGameOver := m.checkGameOver()
		if bGameOver {
			m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sOver") //client js bLive 상태면 win임.
			//게임 & 소켓 룸 제거는 disconnect, leave 호출되면 사라짐.
			break
		}

		//입력 받는 슬립 구간
		m.bShotTime = true
		m.broadcastOneShotStartEnd(true)
		nTimeCount := 5 //그냥 입력 시간 5초?
		for i := nTimeCount; i >= 0; i-- {
			m.broadcastClientTimer(i) //초당 메시지 한 번 씩 전송
			time.Sleep(time.Second)
		}
		m.broadcastOneShotStartEnd(false)
		m.bShotTime = false

		nPlayingTime := 1
		m.setShotInfo(nFrame, nPlayingTime) //사용자 shot 입력 -> player speed 세팅
		m.playing(nFrame, nPlayingTime)     //게임 진행 및 플레이

		m.growing() //캐릭터 성장

		m.bAttackTeamBlack = !m.bAttackTeamBlack //턴 변경
	}

}

func (m *GameRoomInfo) physicsCollision() {
	for _, p1 := range m.playerMap { //black
		if p1.bLive == false || p1.bwTeam == "w" {
			continue
		}
		for _, p2 := range m.playerMap { //white
			if p1 == p2 || p2.bwTeam == "b" || p2.bLive == false {
				continue
			}

			fDistace := math.Sqrt(
				math.Pow(float64(p2.pos.x-p1.pos.x), 2) +
					math.Pow(float64(p2.pos.y-p1.pos.y), 2))
			if fDistace > float64(p1.nRadius+p2.nRadius) { // 거리 조건 스킵
				continue
			}

			if m.bAttackTeamBlack {
				p2.bLive = false
			} else {
				p1.bLive = false
			}
		}
	}
}

func (m *GameRoomInfo) physicsMove() {
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}

		playerInfo.pos.x += playerInfo.speed.x
		playerInfo.pos.y += playerInfo.speed.y
	}
}

func (m *GameRoomInfo) playerSpeedInit() {
	for _, playerInfo := range m.playerMap {
		playerInfo.speed.x = 0
		playerInfo.speed.y = 0
	}

}

func (m *GameRoomInfo) playing(nFrame, nPlayingTime int) {
	loopCnt := nPlayingTime * 1000 / nFrame
	for i := 0; i < loopCnt; i++ {
		m.physicsMove()
		m.physicsCollision()
		m.broadcastPlaying(false)

		time.Sleep(time.Duration(nFrame) * time.Millisecond)
	}
	m.playerSpeedInit()
}
func (m *GameRoomInfo) growing() {
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}
		playerInfo.nRadius += 10
	}
	m.broadcastPlaying(false)
}
func (m *GameRoomInfo) setShotInfo(nFrame, nPlayingTime int) {
	loopCnt := nPlayingTime * 1000 / nFrame

	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			playerInfo.speed.x = 0 //혹시 모르니 초기화 해준다.
			playerInfo.speed.y = 0
			continue
		}

		if playerInfo.shot.bSettedShot == false {
			continue
		}
		playerInfo.shot.bSettedShot = false

		xDistance := float64(playerInfo.shot.target.x) - playerInfo.pos.x
		yDistance := float64(playerInfo.shot.target.y) - playerInfo.pos.y

		playerInfo.speed.x = xDistance / float64(loopCnt)
		playerInfo.speed.y = yDistance / float64(loopCnt)
	}
}

func (m *GameRoomInfo) checkGameOver() bool {
	bTeam := false
	wTeam := false
	nLiveCount := 0
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}
		nLiveCount++

		if playerInfo.bwTeam == "b" {
			bTeam = true
		} else {
			wTeam = true
		}
	}

	bGameOver := bTeam == false || wTeam == false
	return bGameOver
}

//CShot is
func (m *GameRoomInfo) CShot(conID string, msg string) {
	if m.bShotTime == false { //입력 시간 아님.
		return
	}

	player, ok := m.playerMap[conID]
	if !ok {
		return
	}
	if player.shot.bSettedShot == true {
		return
	}
	player.shot.bSettedShot = true

	nIndex := strings.IndexAny(msg, ",")
	if nIndex == -1 || len(msg) == nIndex { //사용자가 js바꿨을 경우.
		return
	}

	x := string(msg[:nIndex])
	y := string(msg[nIndex+1:])
	nX, err := strconv.Atoi(x)
	if err != nil {
		return
	}
	nY, err := strconv.Atoi(y)
	if err != nil {
		return
	}
	recoverPos(&nX)
	recoverPos(&nY)

	player.shot.target.x = nX
	player.shot.target.y = nY
}

func recoverPos(nPos *int) {
	if *nPos < 0 {
		*nPos = 0
	} else if *nPos > 1000 {
		*nPos = 1000
	}
}
