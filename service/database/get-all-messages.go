package database

import "fmt"

func (db *appdbimpl) GetAllMessages() ([]Message, error) {
	query := `
        SELECT Message_id, Sender_id, Recipient_id, MessageContent, Timestamp
        FROM Messages
        WHERE Recipient_id IS NOT NULL
        ORDER BY Timestamp ASC;
    `

	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllMessages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Message_id, &msg.Sender_id, &msg.Recipient_id, &msg.MessageContent, &msg.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("GetAllMessages: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
