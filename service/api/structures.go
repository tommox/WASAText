package api

import (
	"github.com/tommox/WASAText/service/database"
)

// Username
type UserName struct {
	userName string `json:"userName"`
}

// User structure
type User struct {
	userId   int    `json:"userId"`
	userName string `json:"userName"`
}

func (u User) toDataBase() database.User {
	return database.User{
		userId:   u.userId,
		userName: u.userName,
	}
}
