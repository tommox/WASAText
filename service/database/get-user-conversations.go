package database

import (
	"database/sql"
	"fmt"
)

// GetUserConversations recupera sia le conversazioni private che quelle di gruppo di un utente.
func (db *appdbimpl) GetUserConversations(userId int) (map[string]interface{}, error) {
	// Recupera le conversazioni private tra utenti
	queryUserConversations := `
		SELECT Conversation_id, Sender_id, Recipient_id, LastMessage_id, LastMessageTimestamp
		FROM Conversations
		WHERE Sender_id = ? OR Recipient_id = ?
		ORDER BY LastMessageTimestamp DESC;
	`

	rows, err := db.c.Query(queryUserConversations, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("GetUserConversations: errore nel recuperare le conversazioni utente: %w", err)
	}
	defer rows.Close()

	// Array per le conversazioni private
	var userConversations []Conversation
	for rows.Next() {
		var conv Conversation
		var lastMessageID sql.NullInt64
		var lastMessageTimestamp sql.NullTime

		err := rows.Scan(
			&conv.Conversation_id,
			&conv.Sender_id,
			&conv.Recipient_id,
			&lastMessageID,
			&lastMessageTimestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("GetUserConversations: errore nella scansione delle conversazioni utente: %w", err)
		}

		// Gestiamo il valore NULL per LastMessage_id
		if lastMessageID.Valid {
			conv.LastMessage_id = int(lastMessageID.Int64)
		} else {
			conv.LastMessage_id = 0
		}

		// Gestiamo il valore NULL per LastMessageTimestamp
		if lastMessageTimestamp.Valid {
			conv.LastMessageTimestamp = lastMessageTimestamp.Time
		} else {
			conv.LastMessageTimestamp = sql.NullTime{}.Time
		}

		userConversations = append(userConversations, conv)
	}

	// Recupera le conversazioni di gruppo in cui l'utente è membro
	queryGroupConversations := `
		SELECT gc.GroupConversation_id, gc.Group_id, g.Group_name, gc.LastMessage_id, gc.LastMessageTimestamp
		FROM GroupConversations gc
		INNER JOIN GroupMembers gm ON gc.Group_id = gm.Group_id
		INNER JOIN Groups g ON gc.Group_id = g.Group_id
		WHERE gm.User_id = ?
		ORDER BY gc.LastMessageTimestamp DESC;
	`

	groupRows, err := db.c.Query(queryGroupConversations, userId)
	if err != nil {
		return nil, fmt.Errorf("GetUserConversations: errore nel recuperare le conversazioni di gruppo: %w", err)
	}
	defer groupRows.Close()

	// Array per le conversazioni di gruppo
	var groupConversations []GroupConversation
	for groupRows.Next() {
		var groupConv GroupConversation
		var lastMessageID sql.NullInt64
		var lastMessageTimestamp sql.NullTime
		var groupName string

		err := groupRows.Scan(
			&groupConv.GroupConversation_id,
			&groupConv.Group_id,
			&groupName,
			&lastMessageID,
			&lastMessageTimestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("GetUserConversations: errore nella scansione delle conversazioni di gruppo: %w", err)
		}

		// Gestiamo `NULL` per LastMessage_id
		if lastMessageID.Valid {
			groupConv.LastMessage_id = int(lastMessageID.Int64)
		} else {
			groupConv.LastMessage_id = 0
		}

		// Gestiamo `NULL` per LastMessageTimestamp
		if lastMessageTimestamp.Valid {
			groupConv.LastMessageTimestamp = lastMessageTimestamp.Time
		} else {
			groupConv.LastMessageTimestamp = sql.NullTime{}.Time
		}
		groupConv.GroupName = groupName
		groupConversations = append(groupConversations, groupConv)
	}

	// Creiamo una mappa per organizzare meglio i dati
	result := map[string]interface{}{
		"private_conversations": userConversations,
		"group_conversations":   groupConversations,
	}

	return result, nil
}
