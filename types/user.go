package types

type UserInfo struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
}

type AddUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}
