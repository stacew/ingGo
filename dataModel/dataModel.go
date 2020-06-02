package dm

import "stacew/teamgo/sign"

//UserInfo is
type UserInfo struct {
	UserID    int    `json:"userID"`
	Name      string `json:"name"`
	Point     int    `json:"point"`
	WinStreak int    `json:"winStreak"`
	Kill      int    `json:"kill"`
	Resurrect int    `json:"resurrect"`
}

//DataHandlerInterface is
type DataHandlerInterface interface {
	//read & create
	GetAndAddUserInfo(ePlatform sign.PlatformType, platformID string) *UserInfo
	//update
	UpdateUserName(userID int, name string) bool
	UpdateUserPoint(userID, point, winStreak, kill, resurrect int) bool
	//delete
	RemoveUser(userID int) bool
	//close
	Close()
}

//NewDataHandler is
func NewDataHandler(dbConn string) DataHandlerInterface {
	return newSQliteHandler(dbConn)
}
