package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) UpdateOrCreateConversation(sender int, recipient int, messageId int, timestamp time.Time, isRead bool, senderMsgId int) (int, error) {
	var conversationId int
	fmt.Println("isread", isRead)
	query := `
        INSERT INTO Conversations (Sender_id, Recipient_id, LastMessage_id, LastMessageTimestamp, LastMessageIsRead, LastMessageSenderId)
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(Sender_id, Recipient_id) 
        DO UPDATE SET 
            LastMessage_id = excluded.LastMessage_id,
            LastMessageTimestamp = excluded.LastMessageTimestamp,
            LastMessageIsRead = excluded.LastMessageIsRead,
            LastMessageSenderId = excluded.LastMessageSenderId
        RETURNING Conversation_id;
    `
	err := db.c.QueryRow(query, sender, recipient, messageId, timestamp, isRead, senderMsgId).Scan(&conversationId)
	if err != nil {
		return 0, fmt.Errorf("updateOrCreateConversation: %w", err)
	}
	return conversationId, nil
}
