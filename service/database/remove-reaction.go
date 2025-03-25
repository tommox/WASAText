package database

import "fmt"

func (db *appdbimpl) RemoveReaction(messageId int, userId int, isGroup bool) error {
	var query string
	if isGroup {
		query = `DELETE FROM GroupReactions WHERE GroupMessage_id = ? AND User_id = ?`
	} else {
		query = `DELETE FROM Reactions WHERE Message_id = ? AND User_id = ?`
	}
	result, err := db.c.Exec(query, messageId, userId)
	if err != nil {
		return fmt.Errorf("RemoveReaction: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrReactionNotFound
	}

	return nil
}
