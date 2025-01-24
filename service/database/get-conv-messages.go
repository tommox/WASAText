package database

import "fmt"

func (db *appdbimpl) GetConversationMessages(conversationId int) ([]Message, error) {
	query := `
		SELECT Message_id, Sender_id, Recipient_id, MessageContent, Timestamp
		FROM Messages
		WHERE Conversation_id = ?
		ORDER BY Timestamp ASC;
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
