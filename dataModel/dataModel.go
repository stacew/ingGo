package datamodel

//UserInfo is
type UserInfo struct {
	UserID    int    `json:"userID"`
	Name      string `json:"name"`
	Point     int    `json:"point"`
	WinStreak int    `json:"winStreak"`
	Kill      int    `json:"kill"`
	Resurrect int    `json:"resurrect"`
}

type eOAuthPlatform int

const (
	eGoogle eOAuthPlatform = iota
	eFacebook
)

//DataHandlerInterface is
type DataHandlerInterface interface {
	//create
	addUser(ePlatform eOAuthPlatform, sessionID string) *UserInfo
	//read
	GetUserInfo(ePlatform eOAuthPlatform, sessionID string) *UserInfo
	//update
	UpdateUserName(userID int, name string) bool
	UpdateUserPoint(userID, point, winstreak, kill, resurrect int) bool
	//delete
	RemoveUser(userID int) bool
	//close
	Close()
}

//NewDataHandler is
func NewDataHandler(dbConn string) DataHandlerInterface {
	return newSQliteHandler(dbConn)
}
