package gameroom

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

	nRating     int
	killCount   int
	reviveCount int
}

//RoomInfo is
type RoomInfo struct {
	socketioServer *socketio.Server
	nsp            string
	gameRoomName   string
	isStarted      bool

	roomCapacity int
	randPos      []xy

	playerMap        map[string]*playerInfo //key = c.ID()
	bShotTime        bool
	bAttackTeamBlack bool

	turnCount    int
	averageBlack int
	averageWhite int
}

//GetPlayerCount is
func (m *RoomInfo) GetPlayerCount() int {
	return len(m.playerMap)
}

//Join is
func (m *RoomInfo) Join(conID string) string {
	m.playerMap[conID] = &playerInfo{nRating: 10000000} //todo: 유저 레이팅 처리 필요
	return m.gameRoomName
}

//Leave is
func (m *RoomInfo) Leave(conID string) (int, string) {
	delete(m.playerMap, conID)

	nRemainCount := len(m.playerMap)
	return nRemainCount, m.gameRoomName
}

func (m *RoomInfo) startUp() {

	m.isStarted = true

	i := 0
	for _, player := range m.playerMap {
		player.pos.x = float64(m.randPos[i].x)
		player.pos.y = float64(m.randPos[i].y)
		player.nRadius = 30
		i++
		if i%2 == 0 {
			player.bwTeam = "b"
			m.averageBlack += player.nRating
		} else {
			player.bwTeam = "w"
			m.averageWhite += player.nRating
		}

		player.bLive = true
		player.shot.bSettedShot = false
		player.reviveCount = 0
		player.killCount = 0
	}

	m.averageBlack /= (m.roomCapacity / 2)
	m.averageWhite /= (m.roomCapacity / 2)
	m.turnCount = 0
}

/////////////////////////////////////////////

//BroadcastInfoRoom is
func (m *RoomInfo) BroadcastInfoRoom() {
	if m.isStarted { // 시작한 방은 나갈 때 인원 정보 안보여주기 위해서
		return
	}
	msg := ".i" + //.i : infoRoom
		strconv.Itoa(len(m.playerMap)) + "," +
		strconv.Itoa(m.roomCapacity) + ","

	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sDecoder", msg)
}
func (m *RoomInfo) broadcastStart() {
	msg := ".s" //.s: Startin
	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sDecoder", msg)
}

func (m *RoomInfo) broadcastOneShotStartEnd(start bool) {
	msg := ".o" //.o : oneShotStartEnd
	if start {
		msg = msg + "s,"
		m.broadcastPlaying() //shot 시작 시, attackteam 색 알려주기
	} else {
		msg = msg + "e,"
	}
	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sDecoder", msg)
}

func (m *RoomInfo) broadcastClientTimer(nTime int) {
	msg := ".t" + //.t : Timer
		strconv.Itoa(nTime) + ","

	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sDecoder", msg)
}
func (m *RoomInfo) broadcastPlaying() {
	msg := ".a" //.a : attackTeam
	if m.bAttackTeamBlack {
		msg = msg + "b,"
	} else {
		msg = msg + "w,"
	}

	for cID, playerInfo := range m.playerMap {
		live := "d"
		if playerInfo.bLive {
			live = "l"
		}

		msg = msg + ".p" + //.p : playing
			cID + "," +
			strconv.Itoa(int(playerInfo.pos.x)) + "," +
			strconv.Itoa(int(playerInfo.pos.y)) + "," +
			strconv.Itoa(playerInfo.nRadius) + "," +
			live + "," +
			playerInfo.bwTeam + ","
	}
	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sDecoder", msg)
}

func (m *RoomInfo) broadcastOver() {
	msg := ".x" //exit
	m.socketioServer.BroadcastToRoom(m.nsp, m.gameRoomName, "sDecoder", msg)
}

/////////////////////////////////////////////
func (m *RoomInfo) gameOverProcess() {
	bLiveBlack := false
	bLiveWhite := false
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}

		if playerInfo.bwTeam == "b" {
			bLiveBlack = true
		} else {
			bLiveWhite = true
		}
	}

	//todo
	//m.averageBlack
	//m.averageWhite
	//m.broadcastOver()
}

//Start is goroutine
func (m *RoomInfo) Start() {

	m.startUp()
	m.broadcastStart()

	for {
		m.turnCount++

		if m.checkGameOver() {
			m.gameOverProcess()
			//게임 & 소켓 룸 제거는 disconnect, leave 호출되면 사라짐.
			break
		}

		//입력 받는 슬립 구간
		m.bShotTime = true
		m.broadcastOneShotStartEnd(true)
		nTimeCount := 5 //그냥 입력 시간 5초?
		for i := nTimeCount; i >= 0; i-- {
			m.broadcastClientTimer(i)
			time.Sleep(time.Second) //초당 메시지 한 번 씩 전송
		}
		m.broadcastOneShotStartEnd(false)
		m.bShotTime = false

		nPlayingTime := 1 //이동 시간 1초에 그리기
		m.setShotInfo(nFrame, nPlayingTime)
		m.playing(nFrame, nPlayingTime)

		m.growing()
		m.bAttackTeamBlack = !m.bAttackTeamBlack //턴 변경
	}

}
func (m *RoomInfo) physicsCollisionReviveTeam(bwTeam string) {
	bwReviveMap := make(map[string]string)
	//find revive
	for bwOutName, bwOut := range m.playerMap {
		if bwOut.bwTeam == bwTeam || bwOut.bLive == false {
			continue
		}
		for bwInName, bwIn := range m.playerMap {
			if bwIn.bwTeam == bwTeam || bwIn.bLive {
				continue
			} else if bwOut == bwIn {
				continue
			}

			fDistace := math.Sqrt(math.Pow(float64(bwIn.pos.x-bwOut.pos.x), 2) + math.Pow(float64(bwIn.pos.y-bwOut.pos.y), 2))
			if fDistace > float64(bwOut.nRadius+bwIn.nRadius) { // 거리 조건 스킵
				continue
			}

			bwReviveMap[bwOutName] = bwInName
		}
	}
	//revive
	for bwOutName, bwInName := range bwReviveMap {
		m.playerMap[bwOutName].reviveCount++
		m.playerMap[bwInName].bLive = true
	}
}

func (m *RoomInfo) physicsCollisionRevive() {
	m.physicsCollisionReviveTeam("w")
	m.physicsCollisionReviveTeam("b")
}
func (m *RoomInfo) physicsCollisionKill() {
	for _, p1 := range m.playerMap { //black
		if p1.bLive == false || p1.bwTeam == "w" {
			continue
		}
		for _, p2 := range m.playerMap { //white
			if p2.bLive == false || p2.bwTeam == "b" {
				continue
			} else if p1 == p2 {
				continue
			}

			fDistace := math.Sqrt(math.Pow(float64(p2.pos.x-p1.pos.x), 2) + math.Pow(float64(p2.pos.y-p1.pos.y), 2))
			if fDistace > float64(p1.nRadius+p2.nRadius) { // 거리 조건 스킵
				continue
			}

			if m.bAttackTeamBlack {
				p1.killCount++
				p2.bLive = false
			} else {
				p2.killCount++
				p1.bLive = false
			}
		}
	}
}

func (m *RoomInfo) physicsMove() {
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}

		playerInfo.pos.x += playerInfo.speed.x
		playerInfo.pos.y += playerInfo.speed.y
	}
}

func (m *RoomInfo) playerSpeedInit() {
	for _, playerInfo := range m.playerMap {
		playerInfo.speed.x = 0
		playerInfo.speed.y = 0
	}

}

//게임 진행 및 플레이
func (m *RoomInfo) playing(nFrame, nPlayingTime int) {
	loopCnt := nPlayingTime * 1000 / nFrame
	for i := 0; i < loopCnt; i++ {
		m.physicsMove()
		m.physicsCollisionRevive()
		m.physicsCollisionKill()

		m.broadcastPlaying()
		time.Sleep(time.Duration(nFrame) * time.Millisecond)
	}
	m.playerSpeedInit()
}

//캐릭터 성장
func (m *RoomInfo) growing() {
	for _, playerInfo := range m.playerMap {
		if playerInfo.nRadius < 100 {
			playerInfo.nRadius += 10
		}
	}
	m.broadcastPlaying()
}

//사용자 shot 입력 -> player speed 세팅
func (m *RoomInfo) setShotInfo(nFrame, nPlayingTime int) {
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

func (m *RoomInfo) checkGameOver() bool {
	bLiveBlack := false
	bLiveWhite := false
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}

		if playerInfo.bwTeam == "b" {
			bLiveBlack = true
		} else {
			bLiveWhite = true
		}
	}

	bGameOver := bLiveBlack == false || bLiveWhite == false
	return bGameOver
}

//CShot is
func (m *RoomInfo) CShot(conID string, msg string) {
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

//NewRoomInfo is
func NewRoomInfo(socketioServer *socketio.Server, nsp string, nRoomCapacity int) *RoomInfo {
	const rowCount = 3
	const colCount = 3
	const maxPosCount = rowCount * colCount //보드 확장 가능성

	gameRoomName := uuid.New().String()
	if gameRoomName == "" {
		log.Println("[Check uuid.New()]")
		gameRoomName = uuid.New().String()
	}

	gameRoom := &RoomInfo{
		socketioServer:   socketioServer,
		nsp:              nsp,
		gameRoomName:     gameRoomName,
		isStarted:        false,
		roomCapacity:     nRoomCapacity,
		randPos:          make([]xy, maxPosCount, maxPosCount),
		playerMap:        make(map[string]*playerInfo),
		bShotTime:        false,
		bAttackTeamBlack: true,
	}

	//tucker random
	for i := 0; i < rowCount; i++ {
		for j := 0; j < colCount; j++ {
			gameRoom.randPos[i*rowCount+j].x = 150 + i*350
			gameRoom.randPos[i*rowCount+j].y = 150 + j*350
		}
	}

	nSeed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(nSeed)
	for i := 0; i < maxPosCount; i++ {
		n1 := r.Intn(maxPosCount)
		n2 := r.Intn(maxPosCount)
		gameRoom.randPos[n1], gameRoom.randPos[n2] = gameRoom.randPos[n2], gameRoom.randPos[n1]
	}
	//tucker random

	return gameRoom
}
