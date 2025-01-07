package database

import (
	"strconv"
)

func (db *appdbimpl) FindUserId(u User) (int, error) {
	var id string
	err := db.c.QueryRow("SELECT User_id FROM Users WHERE Nickname=?", u.Nickname).Scan(&id)
	if err != nil {
		return -1, err
	}
	int_id, err := strconv.Atoi(id)
	if err != nil {
		return -1, err
	}

	return int_id, nil
}
