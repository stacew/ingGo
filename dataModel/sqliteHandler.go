package datamodel

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" //위의 sql에 대한 go-sqlite3를 사용한다는 암시적 의미로 추가
)

type sqHandler struct {
	db *sql.DB
}

//create ///////////////////////////////////////////////////////////
func (m *sqHandler) addUser(ePlatform eOAuthPlatform, sessionID string) *UserInfo {
	//1 Insert GameUserInfo
	stmt, err := m.db.Prepare("INSERT INTO GameUserInfo (name,point,winStreak,kill,resurrect) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	userInfo := &UserInfo{
		Name:      "noname",
		Point:     1000,
		WinStreak: 0,
		Kill:      0,
		Resurrect: 0,
	}

	rst, err := stmt.Exec(userInfo.Name, userInfo.Point, userInfo.WinStreak, userInfo.Kill, userInfo.Resurrect) //+++
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ := rst.RowsAffected()
	if cnt != 1 {
		log.Fatal(err)
	}

	//2 Select UserID ///////
	userID, err := rst.LastInsertId() //+++
	if err != nil {
		log.Fatal(err)
	}

	//3 Insert ePlatform
	if ePlatform == eGoogle {
		stmt, err = m.db.Prepare("INSERT INTO GoogleUserID (userID, sessionID) VALUES (?, ?)")
	} else if ePlatform == eFacebook {
		stmt, err = m.db.Prepare("INSERT INTO FacebookUserID (userID, sessionID) VALUES (?, ?)")
	}

	if err != nil {
		log.Fatal(err)
	}

	rst, err = stmt.Exec(userID, sessionID) //+++
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
func (m *sqHandler) GetAndAddUserInfo(ePlatform eOAuthPlatform, sessionID string) *UserInfo {
	if ePlatform >= eNone {
		log.Fatal("[FATAL] Database AddUser() : ePlatform Err")
	}

	var rows1 *sql.Rows
	var err error
	//1
	if ePlatform == eGoogle {
		rows1, err = m.db.Query("SELECT (userID) FROM GoogleUserID WHERE (sessionID=?)", sessionID) //+++
	} else if ePlatform == eFacebook {
		rows1, err = m.db.Query("SELECT (userID) FROM FacebookUserID WHERE (sessionID=?)", sessionID) //+++
	}
	if err != nil {
		panic(err)
	}
	defer rows1.Close()

	userID := -1
	rows1.Scan(&userID)
	if userID == -1 {
		return m.addUser(ePlatform, sessionID) //+++++++++++
	}

	//2
	rows2, err := m.db.Query("SELECT (name,point,winStreak,kill,resurrect) FROM GameUserInfo WHERE (userID=?)", userID) //p+++
	if err != nil {
		panic(err)
	}
	defer rows2.Close()

	userInfo := &UserInfo{UserID: userID} //+
	rows2.Scan(&userInfo.Name, &userInfo.Point, &userInfo.WinStreak, &userInfo.Kill, &userInfo.Resurrect)
	return userInfo
}

//update ///////////////////////////////////////////////////////////
func (m *sqHandler) UpdateUserName(userID int, name string) bool {
	stmt, err := m.db.Prepare("UPDATE GameUserInfo SET (name=?) WHERE (userID=?)")
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
func (m *sqHandler) UpdateUserPoint(userID, point, winStreak, kill, resurrect int) bool {
	stmt, err := m.db.Prepare("UPDATE GameUserInfo SET (point=?,winStreak=?,kill=?,resurrect=?) WHERE (userID=?)")
	if err != nil {
		log.Fatal(err)
	}
	rst, err := stmt.Exec(point, winStreak, kill, resurrect, userID)
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
	//
	stmt, err = m.db.Prepare("DELETE FROM FacebookUserID WHERE userID=?")
	if err != nil {
		log.Fatal(err)
	}
	rst, err = stmt.Exec(userID)
	if err != nil {
		log.Fatal(err)
	}
	cnt2, _ := rst.RowsAffected()
	//
	cnt := cnt1 + cnt2
	//
	if cnt == 1 {
		log.Print("[PRINT] Database RemoveUser() : ePlatform Delete User Success")
	} else if cnt == 0 {
		log.Print("[PRINT] Database RemoveUser() : ePlatform Delete User Not Found")
	} else if cnt > 1 {
		log.Fatal("[FATAL] Database RemoveUser() : ePlatform Delete Count : ", cnt)
	}

	//2
	stmt, err = m.db.Prepare("DELETE FROM GameUserInfo WHERE userID=?")
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
		`CREATE TABLE IF NOT EXISTS GameUserInfo (
			userID		INTEGER PRIMARY KEY AUTOINCREMENT,
			name		TEXT,
			point		INTEGER,
			winStreak	INTEGER,
			kill	 	INTEGER,
			resurrect 	INTEGER
		);
		CREATE INDEX IF NOT EXISTS userIDIndexOnGameUserInfo ON GameUserInfo (
			userID ASC
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
		);`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	return &sqHandler{db: database}
}

func (m *sqHandler) Close() {
	m.db.Close()
}
