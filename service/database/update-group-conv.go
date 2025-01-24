package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) updateOrCreateGroupConversation(groupId int, senderId int, messageId int, timestamp time.Time) error {
	query := `
        INSERT INTO GroupConversations (Group_id, Sender_id, LastMessage_id, LastMessageTimestamp)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(Group_id) 
        DO UPDATE SET 
            Sender_id = excluded.Sender_id,
            LastMessage_id = excluded.LastMessage_id,
            LastMessageTimestamp = excluded.LastMessageTimestamp;
    `
	_, err := db.c.Exec(query, groupId, senderId, messageId, timestamp)
	if err != nil {
		return fmt.Errorf("updateOrCreateGroupConversation: %w", err)
	}
	return nil
}
