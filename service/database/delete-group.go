package database

import "fmt"

// DeleteGroup rimuove un gruppo e tutti i dati correlati dal database.
func (db *appdbimpl) DeleteGroup(groupId int) error {
	// Elimina prima tutti i messaggi del gruppo
	queryDeleteMessages := `DELETE FROM GroupMessages WHERE Group_id = ?`
	_, err := db.c.Exec(queryDeleteMessages, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: error deleting group messages: %w", err)
	}

	// Elimina la conversazione del gruppo
	queryDeleteConversation := `DELETE FROM GroupConversations WHERE Group_id = ?`
	_, err = db.c.Exec(queryDeleteConversation, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: error deleting group conversation: %w", err)
	}

	// Elimina i membri del gruppo
	queryDeleteMembers := `DELETE FROM GroupMembers WHERE Group_id = ?`
	_, err = db.c.Exec(queryDeleteMembers, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: error deleting group members: %w", err)
	}

	// Infine, elimina il gruppo stesso
	queryDeleteGroup := `DELETE FROM Groups WHERE Group_id = ?`
	result, err := db.c.Exec(queryDeleteGroup, groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroup: error deleting group: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrGroupNotFound
	}

	return nil
}
