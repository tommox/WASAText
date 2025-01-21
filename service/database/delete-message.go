package database

import "fmt"

// DeleteMessage rimuove un messaggio dal database e resetta il contatore degli ID.
func (db *appdbimpl) DeleteMessage(messageId int) error {
	// Elimina il messaggio specifico
	query := `DELETE FROM Messages WHERE Message_id = ?`
	result, err := db.c.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("DeleteMessage: %w", err)
	}

	// Controlla se il messaggio è stato eliminato
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	// Resetta l'incremento automatico degli ID per la tabella Messages
	_, err = db.c.Exec(`DELETE FROM sqlite_sequence WHERE name='Messages';`)
	if err != nil {
		return fmt.Errorf("DeleteMessage: error resetting auto-increment: %w", err)
	}

	return nil
}
