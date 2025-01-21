package database

import "fmt"

// DeleteMessage rimuove un messaggio dal database, insieme alle reazioni associate.
func (db *appdbimpl) DeleteMessage(messageId int) error {
	// Elimina tutte le reazioni associate al messaggio
	deleteReactionsQuery := `DELETE FROM Reactions WHERE Message_id = ?`
	_, err := db.c.Exec(deleteReactionsQuery, messageId)
	if err != nil {
		return fmt.Errorf("DeleteMessage: error deleting associated reactions: %w", err)
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

	return nil
}
