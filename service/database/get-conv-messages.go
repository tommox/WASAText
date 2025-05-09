package database

import "fmt"

// GetConversationMessages recupera tutti i messaggi (testo e immagini) associati a una conversazione tra due utenti
func (db *appdbimpl) GetConversationMessages(conversationId int) ([]Message, error) {
	query := `
		SELECT m.Message_id, m.Sender_id, m.Conversation_id, m.MessageContent, m.ImageData, m.Timestamp, m.IsRead, m.IsReply, m.IsForward
		FROM Messages m
		WHERE m.Conversation_id = ?
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
			&msg.Conversation_id,
			&msg.MessageContent,
			&msg.ImageData,
			&msg.Timestamp,
			&msg.IsRead,
			&msg.IsReply,
			&msg.IsForward,
		)
		if err != nil {
			return nil, fmt.Errorf("GetConversationMessages: %w", err)
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetConversationMessages: errore iterando le righe: %w", err)
	}
	return messages, nil
}
