package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) updateOrCreateGroupConversation(groupconversation_id int, groupId int, senderId int, messageId int, timestamp time.Time) error {
	query := `
        INSERT INTO GroupConversations (GroupConversation_id, Group_id, Sender_id, LastMessage_id, LastMessageTimestamp)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT(GroupConversation_id) 
        DO UPDATE SET 
            Sender_id = excluded.Sender_id,
            LastMessage_id = excluded.LastMessage_id,
            LastMessageTimestamp = excluded.LastMessageTimestamp;
    `
	_, err := db.c.Exec(query, groupconversation_id, groupId, senderId, messageId, timestamp)
	if err != nil {
		return fmt.Errorf("updateOrCreateGroupConversation: %w", err)
	}
	return nil
}
