package database

import "fmt"

func (db *appdbimpl) DeleteGroupMessage(messageId int) error {
	deleteReactionsQuery := `DELETE FROM Reactions WHERE Message_id = ?`
	_, err := db.c.Exec(deleteReactionsQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nell'eliminare le reazioni associate: %w", err)
	}

	var groupId int
	getGroupIdQuery := `SELECT Group_id FROM GroupMessages WHERE GroupMessage_id = ?`
	err = db.c.QueryRow(getGroupIdQuery, messageId).Scan(&groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nel trovare il group_id del messaggio: %w", err)
	}

	deleteMessageQuery := `DELETE FROM GroupMessages WHERE GroupMessage_id = ?`
	result, err := db.c.Exec(deleteMessageQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nell'eliminare il messaggio: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	var lastMessageId int
	getLastMessageQuery := `SELECT LastMessage_id FROM GroupConversations WHERE Group_id = ?`
	err = db.c.QueryRow(getLastMessageQuery, groupId).Scan(&lastMessageId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nel controllare last_message_id della conversazione: %w", err)
	}

	if lastMessageId == messageId {
		var newLastMessageId *int
		var newLastTimestamp *string
		findNewLastMessageQuery := `
			SELECT GroupMessage_id, Timestamp 
			FROM GroupMessages 
			WHERE Group_id = ? 
			ORDER BY Timestamp DESC 
			LIMIT 1`
		err = db.c.QueryRow(findNewLastMessageQuery, groupId).Scan(&newLastMessageId, &newLastTimestamp)
		if err != nil {
			newLastMessageId = nil
			newLastTimestamp = nil
		}

		updateQuery := `UPDATE GroupConversations SET LastMessage_id = ?, LastMessageTimestamp = ? WHERE Group_id = ?`
		_, err = db.c.Exec(updateQuery, newLastMessageId, newLastTimestamp, groupId)
		if err != nil {
			return fmt.Errorf("DeleteGroupMessage: errore nell'aggiornare last message della conversazione: %w", err)
		}
	}

	return nil
}
