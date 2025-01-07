package api

import (
	"github.com/tommox/WASAText/service/database"
)

// Nickname
type Nickname struct {
	Nickname string `json:"nickname"`
}

// Struttura User
type User struct {
	User_id  int    `json:"user_id"`
	Nickname string `json:"nickname"`
}

func (u User) toDataBase() database.User {
	return database.User{
		User_id:  u.User_id,
		Nickname: u.Nickname,
	}
}
