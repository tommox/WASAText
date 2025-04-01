package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateGroupMessage(groupId int, senderId int, messageContent string, timestamp time.Time, isReply *int) (int, error) {
	query := `
        INSERT INTO GroupMessages (Group_id, Sender_id, MessageContent, Timestamp, IsReply) 
        VALUES (?, ?, ?, ?, ?)
    `
	result, err := db.c.Exec(query, groupId, senderId, messageContent, timestamp, isReply)
	if err != nil {
		return 0, fmt.Errorf("CreateGroupMessage: %w", err)
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateGroupMessage: failed to retrieve message ID: %w", err)
	}

	err = db.updateOrCreateGroupConversation(groupId, groupId, senderId, int(messageId), timestamp)
	if err != nil {
		return 0, fmt.Errorf("CreateGroupMessage: failed to update group conversation: %w", err)
	}

	return int(messageId), nil
}
