package database

func (db *appdbimpl) CheckExistingConversation(userId int, recipientId int) (int, error) {
	var conversationId int
	query := `
        SELECT Conversation_id FROM Conversations 
        WHERE (Sender_id = ? AND Recipient_id = ?) 
           OR (Sender_id = ? AND Recipient_id = ?) 
        LIMIT 1;
    `
	err := db.c.QueryRow(query, userId, recipientId, recipientId, userId).Scan(&conversationId)
	if err != nil {
		return 0, nil
	}
	return conversationId, nil
}
