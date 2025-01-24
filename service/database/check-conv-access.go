package database

import "fmt"

func (db *appdbimpl) CheckConversationAccess(userId, conversationId int) (bool, bool, error) {
	// Controlla l'accesso alle conversazioni normali
	queryConversations := `
		SELECT EXISTS (
			SELECT 1 
			FROM Conversations
			WHERE Conversation_id = ?
			AND (
				Sender_id = ? 
				OR Recipient_id = ?
			)
		)
	`
	var hasAccess bool
	err := db.c.QueryRow(queryConversations, conversationId, userId, userId).Scan(&hasAccess)
	if err != nil {
		return false, false, fmt.Errorf("CheckConversationAccess (normal): %w", err)
	}
	if hasAccess {
		return true, false, nil // Accesso garantito a una conversazione normale
	}

	// Controlla l'accesso alle conversazioni di gruppo
	queryGroupConversations := `
		SELECT EXISTS (
			SELECT 1 
			FROM GroupConversations
			WHERE GroupConversation_id = ?
			AND Group_id IN (
				SELECT Group_id FROM GroupMembers WHERE User_id = ?
			)
		)
	`
	err = db.c.QueryRow(queryGroupConversations, conversationId, userId).Scan(&hasAccess)
	if err != nil {
		return false, false, fmt.Errorf("CheckConversationAccess (group): %w", err)
	}
	if hasAccess {
		return true, true, nil // Accesso garantito a una conversazione di gruppo
	}

	return false, false, nil // Nessun accesso trovato
}
