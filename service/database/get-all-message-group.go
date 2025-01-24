package database

import "fmt"

func (db *appdbimpl) GetAllGroupMessages() ([]GroupMessage, error) {
	query := `
        SELECT GroupMessage_id, Group_id, Sender_id, MessageContent, Timestamp
        FROM GroupMessages
        ORDER BY Timestamp ASC;
    `

	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllGroupMessages: %w", err)
	}
	defer rows.Close()

	var messages []GroupMessage
	for rows.Next() {
		var msg GroupMessage
		err := rows.Scan(&msg.GroupMessage_id, &msg.Group_id, &msg.Sender_id, &msg.MessageContent, &msg.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("GetAllGroupMessages: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
