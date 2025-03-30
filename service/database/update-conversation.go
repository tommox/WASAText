package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) UpdateOrCreateConversation(sender int, recipient int, messageId int, timestamp time.Time) (int, error) {
	var conversationId int
	query := `
        INSERT INTO Conversations (Sender_id, Recipient_id, LastMessage_id, LastMessageTimestamp, LastMessageIsRead, LastMessageSenderId)
        VALUES (?, ?, ?, ?, FALSE, ?)
        ON CONFLICT(Sender_id, Recipient_id) 
        DO UPDATE SET 
            LastMessage_id = excluded.LastMessage_id,
            LastMessageTimestamp = excluded.LastMessageTimestamp,
            LastMessageIsRead = FALSE,
            LastMessageSenderId = excluded.LastMessageSenderId
        RETURNING Conversation_id;
    `
	err := db.c.QueryRow(query, sender, recipient, messageId, timestamp, sender).Scan(&conversationId)
	if err != nil {
		return 0, fmt.Errorf("updateOrCreateConversation: %w", err)
	}
	return conversationId, nil
}
