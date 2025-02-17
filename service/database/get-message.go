package database

import "fmt"

func (db *appdbimpl) GetMessage(messageId int) (Message, error) {
	var msg Message

	// Recupera i dettagli principali del messaggio
	query := `
        SELECT Message_id, Sender_id, Conversation_id, MessageContent, Timestamp
        FROM Messages
        WHERE Message_id = ?;
    `
	err := db.c.QueryRow(query, messageId).Scan(
		&msg.Message_id,
		&msg.Sender_id,
		&msg.Conversation_id,
		&msg.MessageContent,
		&msg.Timestamp,
	)
	if err != nil {
		return msg, fmt.Errorf("GetMessage: error retrieving message: %w", err)
	}
	return msg, nil
}
