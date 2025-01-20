package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func (db *appdbimpl) GetMessage(messageId int) (Message, error) {
	var msg Message
	var reactionsJSON sql.NullString // Variabile temporanea per la colonna JSON

	query := `
        SELECT Message_id, Sender_id, Recipient_id, MessageContent, Timestamp, Reactions
        FROM Messages
        WHERE Message_id = ?;
    `
	err := db.c.QueryRow(query, messageId).Scan(
		&msg.Message_id,
		&msg.Sender_id,
		&msg.Recipient_id,
		&msg.MessageContent,
		&msg.Timestamp,
		&reactionsJSON,
	)
	if err != nil {
		return msg, fmt.Errorf("GetMessage: error retrieving message: %w", err)
	}

	// Decodifica JSON delle reazioni, se esiste
	if reactionsJSON.Valid && reactionsJSON.String != "" {
		err = json.Unmarshal([]byte(reactionsJSON.String), &msg.Reactions)
		if err != nil {
			return msg, fmt.Errorf("GetMessage: error unmarshaling reactions: %w", err)
		}
	} else {
		msg.Reactions = make(map[string][]int) // Inizializza come mappa vuota
	}

	return msg, nil
}
