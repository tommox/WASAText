package database

import "fmt"

func (db *appdbimpl) CheckConversationAccess(userId, conversationId int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM Conversations
			WHERE Conversation_id = ?
			AND (
				Sender_id = ? 
				OR Recipient_id = ?
				OR Group_id IN (
					SELECT Group_id FROM GroupMembers WHERE User_id = ?
				)
			)
		)
	`

	var hasAccess bool
	err := db.c.QueryRow(query, conversationId, userId, userId, userId).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("CheckConversationAccess: %w", err)
	}

	return hasAccess, nil
}
