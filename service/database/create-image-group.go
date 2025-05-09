package database

import (
	"fmt"
	"time"
)

// CreateGroupImageMessage salva un messaggio di gruppo con immagine nel database, inclusi il timestamp.
func (db *appdbimpl) CreateGroupImageMessage(groupId int, senderId int, imageData []byte, timestamp time.Time) (int, error) {
	messageContent := ""
	if len(imageData) > 0 {
		messageContent = ""
	}

	// Query per inserire un messaggio di gruppo con immagine (BLOB) e IsRead
	query := `
        INSERT INTO GroupMessages (Group_id, Sender_id, MessageContent, ImageData, Timestamp, IsRead, IsForward) 
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	// Eseguiamo la query con IsRead = FALSE di default
	result, err := db.c.Exec(query, groupId, senderId, messageContent, imageData, timestamp, false, false)
	if err != nil {
		return 0, fmt.Errorf("CreateGroupImageMessage: %w", err)
	}

	// Recupera l'ID dell'ultimo messaggio inserito
	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateGroupImageMessage: failed to retrieve message ID: %w", err)
	}

	// Aggiorna la conversazione di gruppo
	err = db.updateOrCreateGroupConversation(groupId, groupId, senderId, int(messageId), timestamp)
	if err != nil {
		return 0, fmt.Errorf("CreateGroupImageMessage: failed to update group conversation: %w", err)
	}

	return int(messageId), nil
}
