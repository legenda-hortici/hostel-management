package models

type User struct {
	Number       int
	ID           int
	Username     string
	Surname      string
	Role         string
	Password     string
	Email        string
	Institute    string
	SettlingDate string
	Avatar       string
	Room_id      int
	RoomNumber   int
	HostelNumber int
}

type UserRequest struct {
	Username  string `json:"username"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Institute string `json:"institute"`
	Role      string `json:"role"`
	Password  string `json:"password"`
}
