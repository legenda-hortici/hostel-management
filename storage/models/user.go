package models

import (
	"database/sql"
)

type User struct {
	Number       int
	ID           int
	Username     string
	Role         string
	Password     string
	Email        string
	Institute    sql.NullString
	Room_id      int
	RoomNumber   int
	HostelNumber int
}

type ResidentUpdateRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Institute string `json:"institute"`
	Role      string `json:"role"`
	Password  string `json:"password"`
}
