package database

import (
	"database/sql"
	"fmt"
)

// IsGroupMember verifica se un utente è membro di un determinato gruppo.
func (db *appdbimpl) IsGroupMember(groupId int, userId int) (bool, error) {
	// Query per verificare l'appartenenza al gruppo
	query := `
        SELECT 1 
        FROM GroupMembers 
        WHERE Group_id = ? AND User_id = ?;
    `
	var exists int
	err := db.c.QueryRow(query, groupId, userId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			// L'utente non è membro del gruppo
			return false, nil
		}
		// Altro errore durante l'esecuzione della query
		return false, fmt.Errorf("IsGroupMember: %w", err)
	}

	// L'utente è membro del gruppo
	return true, nil
}
