package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) DeleteGroup(groupId int) error {
	// Elimina il gruppo dal database
	query := `DELETE FROM Groups WHERE Group_id = ?`
	result, err := db.c.Exec(query, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: %w", err)
	}

	// Controlla se è stato eliminato
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows // Gruppo non trovato
	}

	return nil
}
