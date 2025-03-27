package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateImageMessage(senderId int, conversationId int, imageData []byte, timestamp time.Time) (int, error) {

	messageContent := ""
	if len(imageData) > 0 {
		messageContent = ""
	}

	query := `
    INSERT INTO Messages (Sender_id, Conversation_id, ImageData, MessageContent, Timestamp) 
    VALUES (?, ?, ?, ?, ?)
`
	result, err := db.c.Exec(query, senderId, conversationId, imageData, messageContent, timestamp)
	if err != nil {
		return 0, fmt.Errorf("CreateImageMessage: %w", err)
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateImageMessage: failed to retrieve message ID: %w", err)
	}

	updateQuery := `
		UPDATE Conversations 
		SET LastMessage_id = ? 
		WHERE Conversation_id = ?
	`
	_, err = db.c.Exec(updateQuery, messageId, conversationId)
	if err != nil {
		return 0, fmt.Errorf("CreateImageMessage: failed to update lastMessageId: %w", err)
	}

	return int(messageId), nil
}
