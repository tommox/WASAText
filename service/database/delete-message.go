package database

import "fmt"

// DeleteMessage rimuove un messaggio dal database, insieme alle reazioni associate e aggiorna la conversazione corrispondente.
func (db *appdbimpl) DeleteMessage(messageId int) error {
	// Elimina tutte le reazioni associate al messaggio
	deleteReactionsQuery := `DELETE FROM Reactions WHERE Message_id = ?`
	_, err := db.c.Exec(deleteReactionsQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteMessage: errore nell'eliminare le reazioni associate: %w", err)
	}

	// Trova la conversazione a cui appartiene il messaggio
	var conversationId int
	getConversationQuery := `SELECT Conversation_id FROM Messages WHERE Message_id = ?`
	err = db.c.QueryRow(getConversationQuery, messageId).Scan(&conversationId)
	if err != nil {
		return fmt.Errorf("DeleteMessage: errore nel trovare la conversazione del messaggio: %w", err)
	}

	// Elimina il messaggio specifico
	deleteMessageQuery := `DELETE FROM Messages WHERE Message_id = ?`
	result, err := db.c.Exec(deleteMessageQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteMessage: %w", err)
	}

	// Controlla se il messaggio è stato eliminato
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	// Controlla se il messaggio eliminato era il last_message_id della conversazione
	var lastMessageId int
	getLastMessageQuery := `SELECT LastMessage_id FROM Conversations WHERE Conversation_id = ?`
	err = db.c.QueryRow(getLastMessageQuery, conversationId).Scan(&lastMessageId)
	if err != nil {
		return fmt.Errorf("DeleteMessage: errore nel controllare last_message_id della conversazione: %w", err)
	}

	if lastMessageId == messageId {
		// Trova il nuovo ultimo messaggio (il più recente dopo quello eliminato)
		var newLastMessageId *int // Usiamo un puntatore perché potrebbe non esserci un nuovo messaggio
		findNewLastMessageQuery := `
			SELECT Message_id FROM Messages 
			WHERE Conversation_id = ? 
			ORDER BY Timestamp DESC 
			LIMIT 1`
		err = db.c.QueryRow(findNewLastMessageQuery, conversationId).Scan(&newLastMessageId)
		if err != nil {
			// Se non ci sono più messaggi nella conversazione, impostiamo last_message_id a NULL
			newLastMessageId = nil
		}

		// Aggiorna la conversazione con il nuovo last_message_id
		updateConversationQuery := `UPDATE Conversations SET LastMessage_id = ? WHERE Conversation_id = ?`
		_, err = db.c.Exec(updateConversationQuery, newLastMessageId, conversationId)
		if err != nil {
			return fmt.Errorf("DeleteMessage: errore nell'aggiornare last_message_id della conversazione: %w", err)
		}
	}

	return nil
}
