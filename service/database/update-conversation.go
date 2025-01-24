package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) updateOrCreateConversation(sender int, recipient int, messageId int, timestamp time.Time) error {
	query := `
        INSERT INTO Conversations (Sender_id, Recipient_id, LastMessage_id, LastMessageTimestamp)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(Sender_id, Recipient_id) 
        DO UPDATE SET 
            LastMessage_id = excluded.LastMessage_id,
            LastMessageTimestamp = excluded.LastMessageTimestamp;
    `
	_, err := db.c.Exec(query, sender, recipient, messageId, timestamp)
	if err != nil {
		return fmt.Errorf("updateOrCreateConversation: %w", err)
	}
	return nil
}
