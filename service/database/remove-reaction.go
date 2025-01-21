package database

import "fmt"

// RemoveReaction rimuove una reazione specifica dal database.
func (db *appdbimpl) RemoveReaction(messageId int, userId int) error {
	query := `DELETE FROM Reactions WHERE Message_id = ? AND User_id = ?`
	result, err := db.c.Exec(query, messageId, userId)
	if err != nil {
		return fmt.Errorf("RemoveReaction: %w", err)
	}

	// Controlla se la reazione è stata eliminata
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrReactionNotFound
	}

	return nil
}
