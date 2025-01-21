package database

import (
	"fmt"
)

func (db *appdbimpl) RemoveReaction(messageId int, userId int) error {
	query := `DELETE FROM Reactions WHERE Message_id = ? AND User_id = ?`
	result, err := db.c.Exec(query, messageId, userId)
	if err != nil {
		return fmt.Errorf("RemoveReaction: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrReactionNotFound
	}

	// Resetta il contatore se tutte le reazioni sono state eliminate
	_, err = db.c.Exec(`DELETE FROM sqlite_sequence WHERE name='Reactions'`)
	if err != nil {
		return fmt.Errorf("RemoveReaction: error resetting Reaction_id: %w", err)
	}

	return nil
}
