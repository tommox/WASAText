package database

import "fmt"

// GetGroupConversationMessages recupera tutti i messaggi (testo e immagini) associati a una conversazione di gruppo
func (db *appdbimpl) GetGroupConversationMessages(groupConversationId int) ([]GroupMessage, error) {
	query := `
		SELECT gm.GroupMessage_id, gm.Sender_id, gm.Group_id, gm.MessageContent, gm.ImageData, gm.Timestamp, gm.IsRead
		FROM GroupMessages gm
		WHERE gm.Group_id = (
			SELECT GroupConversation_id
			FROM GroupConversations
			WHERE GroupConversation_id = ?
		)
		ORDER BY gm.Timestamp ASC;
	`

	rows, err := db.c.Query(query, groupConversationId)
	if err != nil {
		return nil, fmt.Errorf("GetGroupConversationMessages: %w", err)
	}
	defer rows.Close()

	var groupMessages []GroupMessage
	for rows.Next() {
		var msg GroupMessage
		err := rows.Scan(
			&msg.GroupMessage_id,
			&msg.Sender_id,
			&msg.Group_id,
			&msg.MessageContent,
			&msg.ImageData,
			&msg.Timestamp,
			&msg.IsRead,
		)
		if err != nil {
			return nil, fmt.Errorf("GetGroupConversationMessages: %w", err)
		}
		groupMessages = append(groupMessages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetGroupConversationMessages: errore iterando le righe: %w", err)
	}
	return groupMessages, nil
}
