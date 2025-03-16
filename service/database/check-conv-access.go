package database

import "fmt"

// Controlla l'accesso a una conversazione privata
func (db *appdbimpl) CheckPrivateConversationAccess(userId, conversationId int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM Conversations
			WHERE Conversation_id = ?
			AND (Sender_id = ? OR Recipient_id = ?)
		)
	`
	var hasAccess bool
	err := db.c.QueryRow(query, conversationId, userId, userId).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("CheckPrivateConversationAccess: %w", err)
	}
	return hasAccess, nil
}

// Controlla l'accesso a una conversazione di gruppo
func (db *appdbimpl) CheckGroupConversationAccess(userId, groupConversationId int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM GroupConversations
			WHERE GroupConversation_id = ?
			AND Group_id IN (
				SELECT Group_id FROM GroupMembers WHERE User_id = ?
			)
		)
	`
	var hasAccess bool
	err := db.c.QueryRow(query, groupConversationId, userId).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("CheckGroupConversationAccess: %w", err)
	}
	return hasAccess, nil
}
