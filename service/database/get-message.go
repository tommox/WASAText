package database

func (db *appdbimpl) GetMessage(messageId int) (Message, error) {
	var msg Message
	query := `
        SELECT Message_id, Sender_id, Recipient_id, MessageContent, Timestamp, Reactions
        FROM Messages
        WHERE Message_id = ?;
    `
	err := db.c.QueryRow(query, messageId).Scan(
		&msg.Message_id,
		&msg.Sender_id,
		&msg.Recipient_id,
		&msg.MessageContent,
		&msg.Timestamp,
		&msg.Reactions,
	)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
