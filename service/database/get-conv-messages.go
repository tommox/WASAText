package database

import "fmt"

// GetConversationMessages recupera tutti i messaggi associati a una conversazione tra due utenti
func (db *appdbimpl) GetConversationMessages(conversationId int) ([]Message, error) {
	query := `
		SELECT m.Message_id, m.Sender_id, m.Recipient_id, m.MessageContent, m.Timestamp
		FROM Messages m
		JOIN Conversations c ON 
			((m.Sender_id = c.Sender_id AND m.Recipient_id = c.Recipient_id)
			OR (m.Sender_id = c.Recipient_id AND m.Recipient_id = c.Sender_id))
		WHERE c.Conversation_id = ?
		ORDER BY m.Timestamp ASC;
	`

	rows, err := db.c.Query(query, conversationId)
	if err != nil {
		return nil, fmt.Errorf("GetConversationMessages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.Message_id,
			&msg.Sender_id,
			&msg.Recipient_id,
			&msg.MessageContent,
			&msg.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("GetConversationMessages: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
