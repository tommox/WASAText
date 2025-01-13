package api

import (
	"time"

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

// Struttura Photo
type Complete_Photo struct {
	Photo_id  int       `json:"photo_id"`
	Owner     User      `json:"owner"`
	Timestamp time.Time `json:"timestamp"`
}

func (u User) toDataBase() database.User {
	return database.User{
		User_id:  u.User_id,
		Nickname: u.Nickname,
	}
}
