package datamodel

import (
	"database/sql"
	_ "github.com/lib/pq" //위 sql에 대한 "postgres" 사용한다는 암시적 의미로 추가
	"log"
)

type pqHandler struct {
	db *sql.DB
}

func (s *pqHandler) AddUser(sessionID, name string) *UserInfo {
	stmt, err := s.db.Prepare("INSERT INTO oauth (sessionID, name) VALUES ($1, $2)")
	if err != nil {
		log.Fatalln(err)
	}
	rst, err := stmt.Exec(sessionID, name)
	if err != nil {
		log.Fatalln(err)
	}
	cnt, _ := rst.RowsAffected()
	if cnt == 0 {
		log.Fatalln(err)
	}

	userInfo := &UserInfo{Name: name, Point: 1000}
	return userInfo
}
func (s *pqHandler) GetUserInfo(sessionID string) *UserInfo {
	rows, err := s.db.Query("SELECT name, point FROM oauth WHERE sessionID=$1", sessionID)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	userInfo := &UserInfo{}
	rows.Scan(&userInfo.Name, &userInfo.Point)
	return userInfo
}

func (s *pqHandler) UpdateUserName(sessionID, name string) bool {
	stmt, err := s.db.Prepare("UPDATE oauth SET name=$1 WHERE sessionID=$2")
	if err != nil {
		log.Fatalln(err)
	}
	rst, err := stmt.Exec(name, sessionID)
	if err != nil {
		log.Fatalln(err)
	}

	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *pqHandler) UpdateUserPoint(sessionID, point int) bool {
	stmt, err := s.db.Prepare("UPDATE oauth SET point=$1 WHERE sessionID=$2")
	if err != nil {
		log.Fatalln(err)
	}
	rst, err := stmt.Exec(point, sessionID)
	if err != nil {
		log.Fatalln(err)
	}

	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *pqHandler) RemoveUser(sessionID string) bool {
	stmt, err := s.db.Prepare("DELETE FROM oauth WHERE sessionID=$1")
	if err != nil {
		log.Fatalln(err)
	}
	rst, err := stmt.Exec(sessionID)
	if err != nil {
		log.Fatalln(err)
	}

	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func newPQHandler(dbConn string) DataHandlerInterface {
	database, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalln(err)
	}

	//  id 필요없어 보여서... 빼봄..
	// statement, err := database.Prepare(
	// 	`CREATE TABLE IF NOT EXISTS oauth (
	// 		id 	SERIAL PRIMARY KEY,
	// 		sessionID	VARCHAR(256),
	// 		name		TEXT,
	// 		point		INTEGER,
	// 	);`)
	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS oauth (
			sessionID	VARCHAR(256) PRIMARY KEY,
			name		TEXT,
			point		INTEGER,
		);`)

	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	statement, err = database.Prepare(
		`CREATE INDEX IF NOT EXISTS sessionIDIndexOnOauth ON oauth (
			sessionID ASC
		);`)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	return &pqHandler{db: database}
}

func (s *pqHandler) Close() {
	s.db.Close()
}
