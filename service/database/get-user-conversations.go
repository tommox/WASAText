package database

import (
	"database/sql"
	"fmt"
)

// GetUserConversations recupera sia le conversazioni private che quelle di gruppo di un utente.
func (db *appdbimpl) GetUserConversations(userId int) ([]interface{}, error) {
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
			conv.LastMessage_id = 0 // Oppure un valore che indica "nessun messaggio"
		}

		// Gestiamo il valore NULL per LastMessageTimestamp
		if lastMessageTimestamp.Valid {
			conv.LastMessageTimestamp = lastMessageTimestamp.Time
		} else {
			conv.LastMessageTimestamp = sql.NullTime{}.Time // Impostiamo un timestamp di default
		}

		userConversations = append(userConversations, conv)
	}

	// Recupera le conversazioni di gruppo
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
		return nil, fmt.Errorf("GetUserConversations: errore nel recuperare le conversazioni di gruppo: %w", err)
	}
	defer groupRows.Close()

	// Array per le conversazioni di gruppo
	var groupConversations []GroupConversation
	for groupRows.Next() {
		var groupConv GroupConversation
		var lastMessageID sql.NullInt64
		var lastMessageTimestamp sql.NullTime

		err := groupRows.Scan(
			&groupConv.GroupConversation_id,
			&groupConv.Group_id,
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

		groupConversations = append(groupConversations, groupConv)
	}

	// Creiamo un array misto di conversazioni private e di gruppo
	allConversations := make([]interface{}, 0, len(userConversations)+len(groupConversations))

	// Aggiungiamo le conversazioni private
	for _, conv := range userConversations {
		allConversations = append(allConversations, conv)
	}

	// Aggiungiamo le conversazioni di gruppo
	for _, groupConv := range groupConversations {
		allConversations = append(allConversations, groupConv)
	}

	// Se non ci sono conversazioni, restituiamo un array vuoto
	return allConversations, nil
}
