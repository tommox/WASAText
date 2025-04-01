package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateMessage(senderId int, conversationId int, messageContent string, timestamp time.Time, isReply *int, isForward bool) (int, error) {
	query := `
        INSERT INTO Messages (Sender_id, Conversation_id, MessageContent, Timestamp, IsReply, IsForward) 
        VALUES (?, ?, ?, ?, ?, ?)
    `
	result, err := db.c.Exec(query, senderId, conversationId, messageContent, timestamp, isReply, isForward)
	if err != nil {
		return 0, fmt.Errorf("CreateMessage: %w", err)
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateMessage: failed to retrieve message ID: %w", err)
	}

	updateQuery := `
		UPDATE Conversations 
		SET LastMessage_id = ? 
		WHERE Conversation_id = ?
	`
	_, err = db.c.Exec(updateQuery, messageId, conversationId)
	if err != nil {
		return 0, fmt.Errorf("CreateMessage: failed to update lastMessageId: %w", err)
	}

	return int(messageId), nil
}
