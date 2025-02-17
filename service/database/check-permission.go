package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) CheckUserPermission(userId, messageId int) (bool, error) {
	var exists int
	query := `
        SELECT 1 
        FROM Messages 
        WHERE Message_id = ? AND (Sender_id = ?);`
	err := db.c.QueryRow(query, messageId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists == 1, nil
}
