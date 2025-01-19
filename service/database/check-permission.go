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
        WHERE Message_id = ? AND (Sender_id = ? OR Recipient_id = ?);`
	err := db.c.QueryRow(query, messageId, userId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Nessun risultato trovato
		}
		return false, err
	}
	return exists == 1, nil
}
