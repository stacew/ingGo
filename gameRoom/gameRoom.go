package gameroom

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

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
	//for game
	bwTeam byte //'b' is black, 'w' is white
	bLive  bool
	pos    xyfloat64
	shot   shotInfo
	speed  xyfloat64
	//for point
	gameUserID   int
	nRating      int
	liveTurn     int
	killCount    int
	reviveCount  int
	fResultPoint float64
}

//RoomInfo is
type RoomInfo struct {
	//for service 1
	socketioServer *socketio.Server
	nsp            string
	roomName       string
	//for service 2
	roomCapacity int
	playerMap    map[string]*playerInfo //key : conID
	isStarted    bool
	//for game
	nRadius          int
	bShotTime        bool
	bAttackTeamBlack bool
	//for point
	fAverageBlack float64
	fAverageWhite float64
	turnCount     int
}

//GetPlayerCount is
func (m *RoomInfo) GetPlayerCount() int {
	return len(m.playerMap)
}

//Join is
func (m *RoomInfo) Join(conID string) string {
	m.playerMap[conID] = &playerInfo{
		nRating:    10000000, //todo: 유저 레이팅 처리 필요
		gameUserID: 12313213} //todo: 유저 아이디 필요
	return m.roomName
}

//Leave is
func (m *RoomInfo) Leave(conID string) (int, string) {
	delete(m.playerMap, conID)

	nRemainCount := len(m.playerMap)
	return nRemainCount, m.roomName
}

func (m *RoomInfo) startUp() {

	//random Pos
	const rowCount = 3
	const colCount = 3
	const maxPosCount = 9
	randPos := make([]xy, maxPosCount)
	for i := 0; i < rowCount; i++ {
		for j := 0; j < colCount; j++ {
			randPos[i*rowCount+j].x = 150 + i*350
			randPos[i*rowCount+j].y = 150 + j*350
		}
	}
	nSeed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(nSeed)
	for i := 0; i < maxPosCount; i++ {
		n1 := r.Intn(maxPosCount)
		n2 := r.Intn(maxPosCount)
		randPos[n1], randPos[n2] = randPos[n2], randPos[n1]
	}
	//random Pos

	i := 0
	teamRatingBlack := 0
	teamRatingWhite := 0
	for _, player := range m.playerMap {

		player.pos.x = float64(randPos[i].x)
		player.pos.y = float64(randPos[i].y)
		//player.speed.x = 0
		//player.speed.y = 0
		i++
		if i%2 == 0 {
			player.bwTeam = 'b'
			teamRatingBlack += player.nRating
		} else {
			player.bwTeam = 'w'
			teamRatingWhite += player.nRating
		}

		player.bLive = true
		player.shot.bSettedShot = false
		player.liveTurn = 1 // warn div 0
		player.killCount = 0
		player.reviveCount = 0
		player.fResultPoint = 0
	}

	//for game
	m.nRadius = 20
	m.bShotTime = false
	m.bAttackTeamBlack = true
	//for point
	m.fAverageBlack = float64(teamRatingBlack / (m.roomCapacity / 2))
	m.fAverageWhite = float64(teamRatingWhite / (m.roomCapacity / 2))
	m.turnCount = 1 // warn div 0
}

//BroadcastInfoRoom is
func (m *RoomInfo) BroadcastInfoRoom(conID string) {
	msg := ""
	if m.isStarted { // 시작한 방은 나갈 때 인원 정보 안보여주기 위해서
		msg = msg + ".l" + //.l : leave
			conID + ","
	} else {
		msg = msg + ".i" + //.i : infoRoom
			strconv.Itoa(len(m.playerMap)) + "," +
			strconv.Itoa(m.roomCapacity) + "," +
			conID + ","
	}

	m.socketioServer.BroadcastToRoom(m.nsp, m.roomName, "sDecoder", msg)
}

func (m *RoomInfo) broadcastStart() {
	msg := ".s" //.s: StartGame
	m.socketioServer.BroadcastToRoom(m.nsp, m.roomName, "sDecoder", msg)

}

func (m *RoomInfo) broadcastOneShotStartEnd(start bool) {
	msg := ".o" //.o : oneShotStartEnd
	if start {
		msg = msg + "s,"
		m.broadcastPlaying() //shot 시작 시, attackteam 색 알려주기
	} else {
		msg = msg + "e,"
	}
	m.socketioServer.BroadcastToRoom(m.nsp, m.roomName, "sDecoder", msg)
}

func (m *RoomInfo) broadcastClientTimer(nTime int) {
	msg := ".t" + //.t : Timer
		strconv.Itoa(nTime) + ","

	m.socketioServer.BroadcastToRoom(m.nsp, m.roomName, "sDecoder", msg)
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
			live + "," +
			string(playerInfo.bwTeam) + ","
	}
	m.socketioServer.BroadcastToRoom(m.nsp, m.roomName, "sDecoder", msg)
}

func (m *RoomInfo) broadcastOver() {
	msg := ".x" //exit
	m.socketioServer.BroadcastToRoom(m.nsp, m.roomName, "sDecoder", msg)
}

/////////////////////////////////////////////
func (m *RoomInfo) calcPoint(bwTeam byte, bWin bool) {
	teamTotal := float64(10 * m.roomCapacity / 2)

	diffAverage := 0.0
	if bwTeam == 'b' {
		diffAverage = m.fAverageWhite - m.fAverageBlack
	} else {
		diffAverage = m.fAverageBlack - m.fAverageBlack
	}

	if bWin && diffAverage > 0 {
		teamTotal += diffAverage
	} else if bWin == false && diffAverage < 0 {
		teamTotal += diffAverage
	}

	totalKill := 0
	totalRevive := 0
	totalContribute := 0.0
	for _, p := range m.playerMap {
		if p.bwTeam != bwTeam {
			continue
		}
		totalKill += p.killCount
		totalRevive += p.reviveCount
		p.fResultPoint = float64((p.killCount + p.reviveCount + 1) / p.liveTurn)
		totalContribute += p.fResultPoint
	}

	for _, p := range m.playerMap {
		if p.bwTeam != bwTeam {
			continue
		}

		fUserContribute := p.fResultPoint / totalContribute
		p.fResultPoint = teamTotal * fUserContribute
		if bWin == false {
			p.fResultPoint = -p.fResultPoint
		}
	}
}

func (m *RoomInfo) gameOverProcess() {
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive == false {
			continue
		}

		if playerInfo.bwTeam == 'b' {
			m.calcPoint('b', true)
			m.calcPoint('w', false)
			break
		} else {
			m.calcPoint('b', true)
			m.calcPoint('w', false)
			break
		}
	}

	//todo
	//m.broadcastOver()
}

//Start is goroutine
func (m *RoomInfo) Start() {
	m.isStarted = true
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

		m.growingAndLiveTurn()
		m.bAttackTeamBlack = !m.bAttackTeamBlack //턴 변경
	}

}
func (m *RoomInfo) physicsCollisionReviveTeam(bwTeam byte) {
	bwReviveMap := make(map[string]string)
	//find revive
	for outCon, outInfo := range m.playerMap {
		if outInfo.bwTeam == bwTeam || outInfo.bLive == false {
			continue
		}
		for inCon, inInfo := range m.playerMap {
			if inInfo.bwTeam == bwTeam || inInfo.bLive {
				continue
			} else if outInfo == inInfo {
				continue
			}

			fDistace := math.Sqrt(math.Pow(float64(inInfo.pos.x-outInfo.pos.x), 2) +
				math.Pow(float64(inInfo.pos.y-outInfo.pos.y), 2))

			if fDistace > float64(m.nRadius+m.nRadius) { // 거리 조건 스킵
				continue
			}

			bwReviveMap[outCon] = inCon
		}
	}
	//revive
	for outCon, inCon := range bwReviveMap {
		m.playerMap[outCon].reviveCount++
		m.playerMap[inCon].bLive = true
	}
}

func (m *RoomInfo) physicsCollisionRevive() {
	m.physicsCollisionReviveTeam('w')
	m.physicsCollisionReviveTeam('b')
}
func (m *RoomInfo) physicsCollisionKill() {
	for _, p1 := range m.playerMap { //black
		if p1.bLive == false || p1.bwTeam == 'w' {
			continue
		}
		for _, p2 := range m.playerMap { //white
			if p2.bLive == false || p2.bwTeam == 'b' {
				continue
			} else if p1 == p2 {
				continue
			}

			fDistace := math.Sqrt(math.Pow(float64(p2.pos.x-p1.pos.x), 2) + math.Pow(float64(p2.pos.y-p1.pos.y), 2))
			if fDistace > float64(m.nRadius+m.nRadius) { // 거리 조건 스킵
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
func (m *RoomInfo) growingAndLiveTurn() {
	for _, playerInfo := range m.playerMap {
		if playerInfo.bLive {
			playerInfo.liveTurn++
		}

		if m.nRadius < 100 {
			m.nRadius += 10
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

		if playerInfo.bwTeam == 'b' {
			bLiveBlack = true
		} else {
			bLiveWhite = true
		}
	}

	bGameOver := bLiveBlack == false || bLiveWhite == false
	return bGameOver
}

//ClientShot is
func (m *RoomInfo) ClientShot(conID string, msg string) {
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
func NewRoomInfo(socketioServer *socketio.Server, nsp string, nRoomCapacity, nMakedRoomCount int) *RoomInfo {
	gameRoom := &RoomInfo{
		socketioServer: socketioServer,
		nsp:            nsp,
		roomName:       strconv.Itoa(nMakedRoomCount),

		roomCapacity: nRoomCapacity,
		playerMap:    make(map[string]*playerInfo),

		fAverageBlack: 0,
		fAverageWhite: 0,

		isStarted:        false,
		bShotTime:        false,
		bAttackTeamBlack: true,
		turnCount:        1, //warn div 0
	}
	return gameRoom
}
