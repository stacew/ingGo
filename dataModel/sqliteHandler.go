package datamodel

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3" //위의 sql에 대한 go-sqlite3를 사용한다는 암시적 의미로 추가
)

type sqHandler struct {
	db    *sql.DB
	mutex *sync.RWMutex
}

func (m *sqHandler) genID() int {
	m.mutex.Lock() //mutex
	//1
	rows, err := m.db.Query("SELECT genID FROM UserIdGenerator") //+++
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	genID := -1
	rows.Scan(&genID) //+

	//2
	stmt, err := m.db.Prepare("UPDATE UserIdGenerator SET genID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err := stmt.Exec(genID + 1) //+++
	if err != nil {
		log.Fatal(err)
	}
	cnt, _ := rst.RowsAffected()
	if cnt != 1 {
		log.Fatal("[FATAL] Database genID() : update AffectedCount")
	}

	m.mutex.Unlock() //mutex
	return genID
}

//create ///////////////////////////////////////////////////////////

func (m *sqHandler) addUser(ePlatform eOAuthPlatform, sessionID string) *UserInfo {
	var stmt *sql.Stmt
	var err error
	//1
	if ePlatform == eGoogle {
		stmt, err = m.db.Prepare("INSERT INTO GoogleUserID (sessionID, userID) VALUES (?, ?)")
	} else if ePlatform == eFacebook {
		stmt, err = m.db.Prepare("INSERT INTO FacebookUserID (sessionID, userID) VALUES (?, ?)")
	} else {
		log.Fatal("[FATAL] Database AddUser() : ePlatform Err")
	}

	if err != nil {
		log.Fatal(err)
	}

	userID := m.genID()                      //+++
	rst, err := stmt.Exec(sessionID, userID) //+++ 프라이머리 키 : sessionID (중복 삽입 시, err 뜨는지 봐야 함. Get() 오류로 Add()하면 Fatal로 강제종료 위험할 수 있음.)
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ := rst.RowsAffected()
	if cnt != 1 {
		log.Fatal(err)
	}

	//2
	stmt, err = m.db.Prepare("INSERT INTO GameInfo (userID,name,point,winStreak,kill,resurrect) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	userInfo := &UserInfo{
		UserID:    userID,
		Name:      "noname",
		Point:     1000,
		WinStreak: 0,
		Kill:      0,
		Resurrect: 0,
	}

	//+++ 프라이머리 키 : userInfo.UserID
	rst, err = stmt.Exec(userInfo.UserID, userInfo.Name, userInfo.Point, userInfo.WinStreak, userInfo.Kill, userInfo.Resurrect) //+++
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ = rst.RowsAffected()
	if cnt != 1 {
		log.Fatal(err)
	}

	return userInfo
}

//read ///////////////////////////////////////////////////////////
func (m *sqHandler) GetUserInfo(ePlatform eOAuthPlatform, sessionID string) *UserInfo {
	var rows1 *sql.Rows
	var err error
	//1
	if ePlatform == eGoogle {
		rows1, err = m.db.Query("SELECT userID FROM GoogleUserID WHERE sessionID=?", sessionID) //p+++
	} else if ePlatform == eFacebook {
		rows1, err = m.db.Query("SELECT userID FROM FacebookUserID WHERE sessionID=?", sessionID) //p+++
	} else {
		log.Fatal("[FATAL] Database GetUserInfo() : ePlatform Err")
	}
	if err != nil {
		panic(err)
	}
	defer rows1.Close()
	userID := -1
	rows1.Scan(&userID)
	if userID == -1 {
		return nil
	}

	//2
	rows2, err := m.db.Query("SELECT (name,point,winStreak,kill,resurrect) FROM GameInfo WHERE userID=?", userID) //p+++
	if err != nil {
		panic(err)
	}
	defer rows2.Close()
	userInfo := &UserInfo{}
	rows2.Scan(&userInfo.UserID, &userInfo.Name, &userInfo.Point, &userInfo.WinStreak, &userInfo.Kill, &userInfo.Resurrect)
	return userInfo
}

//update ///////////////////////////////////////////////////////////
func (m *sqHandler) UpdateUserName(userID int, name string) bool {
	stmt, err := m.db.Prepare("UPDATE GameInfo SET name=? WHERE userID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err := stmt.Exec(name, userID)
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ := rst.RowsAffected()
	return cnt == 1
}
func (m *sqHandler) UpdateUserPoint(userID, point, winstreak, kill, resurrect int) bool {
	stmt, err := m.db.Prepare("UPDATE GameInfo SET (point=?,winStreak=?,kill=?,resurrect=?) WHERE userID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err := stmt.Exec(point, winstreak, kill, resurrect, userID)
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ := rst.RowsAffected()
	return cnt == 1
}

//delete ///////////////////////////////////////////////////////////
func (m *sqHandler) RemoveUser(userID int) bool {
	//1
	stmt, err := m.db.Prepare("DELETE FROM GoogleUserID WHERE userID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err := stmt.Exec(userID)
	if err != nil {
		log.Fatal(err)
	}
	cnt1, _ := rst.RowsAffected()

	stmt, err = m.db.Prepare("DELETE FROM FacebookUserID WHERE userID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err = stmt.Exec(userID)
	if err != nil {
		log.Fatal(err)
	}
	cnt2, _ := rst.RowsAffected()

	cnt := cnt1 + cnt2
	if cnt == 1 {
		log.Print("[PRINT] Database RemoveUser() : ePlatform Delete User Success")
	} else if cnt == 0 {
		log.Print("[PRINT] Database RemoveUser() : ePlatform Delete User Not Found")
	} else if cnt > 1 {
		log.Fatal("[FATAL] Database RemoveUser() : ePlatform Delete Count : ", cnt)
	}

	//2
	stmt, err = m.db.Prepare("DELETE FROM GameInfo WHERE userID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err = stmt.Exec(userID)
	if err != nil {
		log.Fatal(err)
	}
	cnt, _ = rst.RowsAffected()

	return cnt == 1
}

func newSQliteHandler(filePath string) DataHandlerInterface {
	database, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS UserIdGenerator (
			genID	INTEGER
		);

		CREATE TABLE IF NOT EXISTS GoogleUserID (
			userID		INTEGER PRIMARY KEY,
			sessionID	STRING
		);

		CREATE INDEX IF NOT EXISTS sessionIDIndexOnGoogleUserID ON GoogleUserID (
			sessionID ASC
		);

		CREATE TABLE IF NOT EXISTS FacebookUserID (
			userID		INTEGER PRIMARY KEY,
			sessionID	STRING
		);

		CREATE INDEX IF NOT EXISTS sessionIDIndexOnFacebookUserID ON FacebookUserID (
			sessionID ASC
		);

		CREATE TABLE IF NOT EXISTS GameInfo (
			userID		INTEGER PRIMARY KEY AUTOINCREMENT,
			name		TEXT,
			point		INTEGER,
			winStreak	INTEGER,
			kill	 	INTEGER,
			resurrect 	INTEGER
		);

		CREATE INDEX IF NOT EXISTS userIDIndexOnGameInfo ON GameInfo (
			userID ASC
		);`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	return &sqHandler{db: database, mutex: new(sync.RWMutex)}
}

func (m *sqHandler) Close() {
	m.db.Close()
}
