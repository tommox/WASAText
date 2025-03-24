package database

import "fmt"

func (db *appdbimpl) DeleteGroupMessage(messageId int) error {
	// 1. Elimina tutte le reazioni associate al messaggio
	deleteReactionsQuery := `DELETE FROM Reactions WHERE Message_id = ?`
	_, err := db.c.Exec(deleteReactionsQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nell'eliminare le reazioni associate: %w", err)
	}

	// 2. Trova il group_id a cui appartiene il messaggio
	var groupId int
	getGroupIdQuery := `SELECT Group_id FROM GroupMessages WHERE GroupMessage_id = ?`
	err = db.c.QueryRow(getGroupIdQuery, messageId).Scan(&groupId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nel trovare il group_id del messaggio: %w", err)
	}

	// 3. Elimina il messaggio
	deleteMessageQuery := `DELETE FROM GroupMessages WHERE GroupMessage_id = ?`
	result, err := db.c.Exec(deleteMessageQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nell'eliminare il messaggio: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	// 4. Ottieni il last_message_id della conversazione di gruppo
	var lastMessageId int
	getLastMessageQuery := `SELECT LastMessage_id FROM GroupConversations WHERE Group_id = ?`
	err = db.c.QueryRow(getLastMessageQuery, groupId).Scan(&lastMessageId)
	if err != nil {
		return fmt.Errorf("DeleteGroupMessage: errore nel controllare last_message_id della conversazione: %w", err)
	}

	if lastMessageId == messageId {
		// 5. Trova il nuovo ultimo messaggio (il più recente dopo quello eliminato)
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
			// Se non ci sono più messaggi, resettiamo i campi
			newLastMessageId = nil
			newLastTimestamp = nil
		}

		// 6. Aggiorna la conversazione di gruppo
		updateQuery := `UPDATE GroupConversations SET LastMessage_id = ?, LastMessageTimestamp = ? WHERE Group_id = ?`
		_, err = db.c.Exec(updateQuery, newLastMessageId, newLastTimestamp, groupId)
		if err != nil {
			return fmt.Errorf("DeleteGroupMessage: errore nell'aggiornare last message della conversazione: %w", err)
		}
	}

	return nil
}
