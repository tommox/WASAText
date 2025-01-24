package database

import (
	"fmt"
)

// CreateGroupMessage salva un messaggio di gruppo nel database.
func (db *appdbimpl) CreateGroupMessage(groupId int, senderId int, messageContent string) (int, error) {
	// Query per inserire un messaggio nella tabella GroupMessages
	query := `
        INSERT INTO GroupMessages (Group_id, Sender_id, MessageContent) 
        VALUES (?, ?, ?)
    `
	// Esegui la query
	result, err := db.c.Exec(query, groupId, senderId, messageContent)
	if err != nil {
		return 0, fmt.Errorf("CreateGroupMessage: %w", err)
	}

	// Recupera l'ID dell'ultimo messaggio inserito
	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateGroupMessage: failed to retrieve message ID: %w", err)
	}

	return int(messageId), nil
}
