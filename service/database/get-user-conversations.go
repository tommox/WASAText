package database

import "fmt"

func (db *appdbimpl) GetUserConversations(userId int) ([]interface{}, error) {
	// Recupera le conversazioni normali (messaggi tra utenti)
	queryUserConversations := `
		SELECT Conversation_id, Sender_id, Recipient_id, LastMessage_id, LastMessageTimestamp
		FROM Conversations
		WHERE Sender_id = ? OR Recipient_id = ?
		ORDER BY LastMessageTimestamp DESC;
	`

	rows, err := db.c.Query(queryUserConversations, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("GetUserConversations: error fetching user conversations: %w", err)
	}
	defer rows.Close()

	// Array per le conversazioni normali (messaggi tra utenti)
	var userConversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(
			&conv.Conversation_id,
			&conv.Sender_id,
			&conv.Recipient_id,
			&conv.LastMessage_id,
			&conv.LastMessageTimestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("GetUserConversations: error scanning user conversations: %w", err)
		}
		userConversations = append(userConversations, conv)
	}

	// Recupera le conversazioni di gruppo (messaggi nei gruppi)
	queryGroupConversations := `
		SELECT GroupConversation_id, Group_id, LastMessage_id, LastMessageTimestamp
		FROM GroupConversations
		WHERE Group_id IN (
			SELECT Group_id FROM GroupMembers WHERE User_id = ?
		)
		ORDER BY LastMessageTimestamp DESC;
	`

	groupRows, err := db.c.Query(queryGroupConversations, userId)
	if err != nil {
		return nil, fmt.Errorf("GetUserConversations: error fetching group conversations: %w", err)
	}
	defer groupRows.Close()

	// Array per le conversazioni di gruppo (messaggi nei gruppi)
	var groupConversations []GroupConversation
	for groupRows.Next() {
		var groupConv GroupConversation
		err := groupRows.Scan(
			&groupConv.GroupConversation_id,
			&groupConv.Group_id,
			&groupConv.LastMessage_id,
			&groupConv.LastMessageTimestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("GetUserConversations: error scanning group conversations: %w", err)
		}
		groupConversations = append(groupConversations, groupConv)
	}

	// Combina i risultati in un unico slice
	var allConversations []interface{}
	// Aggiungi le conversazioni normali
	for _, conv := range userConversations {
		allConversations = append(allConversations, conv)
	}
	// Aggiungi le conversazioni di gruppo
	for _, groupConv := range groupConversations {
		allConversations = append(allConversations, groupConv)
	}

	return allConversations, nil
}
