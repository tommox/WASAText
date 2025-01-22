package database

import "fmt"

// DeleteGroup rimuove un gruppo e tutti i membri associati dal database.
func (db *appdbimpl) DeleteGroup(groupId int) error {
	// Elimina i membri del gruppo
	queryRemoveUsers := `DELETE FROM GroupMembers WHERE Group_id = ?`
	_, err := db.c.Exec(queryRemoveUsers, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: error deleting group members: %w", err)
	}

	// Ripristina il contatore AUTOINCREMENT
	queryReset := `DELETE FROM sqlite_sequence WHERE name = 'GroupMembers'`
	_, err = db.c.Exec(queryReset)
	if err != nil {
		return fmt.Errorf("ResetGroupMemberId: error resetting AUTOINCREMENT: %w", err)
	}

	// Elimina il gruppo
	query := `DELETE FROM Groups WHERE Group_id = ?`
	result, err := db.c.Exec(query, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: error deleting group: %w", err)
	}

	// Verifica se il gruppo è stato eliminato
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrGroupNotFound
	}

	return nil
}
