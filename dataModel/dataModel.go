package datamodel

//UserInfo is
type UserInfo struct {
	Name      string `json:"name"`
	Point     int    `json:"point"`
	WinStreak int    `json:"winstreak"`
}

//DataHandlerInterface is
type DataHandlerInterface interface {
	//create
	AddUser(sessionID string) *UserInfo
	//read
	GetUserInfo(sessionID string) *UserInfo
	//update
	UpdateUserName(sessionID, name string) bool
	UpdateUserPoint(sessionID, point, winstreak int) bool
	//delete
	RemoveUser(sessionID string) bool
	//close
	Close()
}

//NewDataHandler is
func NewDataHandler(dbConn string) DataHandlerInterface {
	return newPQHandler(dbConn)
}
