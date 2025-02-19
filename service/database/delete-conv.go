package database

import "fmt"

// DeleteConversation rimuove una conversazione e i messaggi associati.
func (db *appdbimpl) DeleteConversation(conversationId int) error {
	// Eliminare tutti i messaggi associati alla conversazione
	deleteMessagesQuery := `DELETE FROM Messages WHERE Conversation_id = ?`
	_, err := db.c.Exec(deleteMessagesQuery, conversationId)
	if err != nil {
		return fmt.Errorf("DeleteConversation: error deleting associated messages: %w", err)
	}

	// Eliminare la conversazione
	deleteConversationQuery := `DELETE FROM Conversations WHERE Conversation_id = ?`
	result, err := db.c.Exec(deleteConversationQuery, conversationId)
	if err != nil {
		return fmt.Errorf("DeleteConversation: %w", err)
	}

	// Controlla se la conversazione è stata effettivamente eliminata
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}
