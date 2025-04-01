package database

import "fmt"

// Funzione per aggiornare il campo "IsForward" in un messaggio esistente
func (db *appdbimpl) MarkIsForward(messageId int, isForward bool) error {
	query := `UPDATE Messages SET IsForward = ? WHERE Message_id = ?`
	_, err := db.c.Exec(query, isForward, messageId)
	if err != nil {
		return fmt.Errorf("MarkIsForward: %w", err)
	}
	return nil
}

// Funzione per aggiornare il campo "IsForward" in un messaggio di gruppo esistente
func (db *appdbimpl) MarkIsForwardGroup(messageId int, isForward bool) error {
	query := `UPDATE GroupMessages SET IsForward = ? WHERE GroupMessage_id = ?`
	_, err := db.c.Exec(query, isForward, messageId)
	if err != nil {
		return fmt.Errorf("MarkIsForwardGroup: %w", err)
	}
	return nil
}
