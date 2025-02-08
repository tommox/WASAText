package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateMessage(senderId int, conversationId int, messageContent string, timestamp time.Time) (int, error) {
	query := `
        INSERT INTO Messages (Sender_id, Conversation_id, MessageContent, Timestamp) 
        VALUES (?, ?, ?, ?)
    `
	result, err := db.c.Exec(query, senderId, conversationId, messageContent, timestamp)
	if err != nil {
		return 0, fmt.Errorf("CreateMessage: %w", err)
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateMessage: failed to retrieve message ID: %w", err)
	}
	return int(messageId), nil
}
