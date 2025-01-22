package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) IsGroupAdmin(groupId int, userId int) (bool, error) {
	var exists int
	query := `
        SELECT 1 
        FROM GroupMembers 
        WHERE Group_id = ? AND User_id = ? AND Role = 'admin';
    `
	err := db.c.QueryRow(query, groupId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("IsGroupAdmin: %w", err)
	}
	return exists == 1, nil
}
