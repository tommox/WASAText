package database

import "fmt"

func (db *appdbimpl) GetConversationIdByMessageId(messageId int) (int, error) {
	var conversationId int
	query := `SELECT Conversation_id FROM Messages WHERE Message_id = ?`
	err := db.c.QueryRow(query, messageId).Scan(&conversationId)
	if err != nil {
		return 0, fmt.Errorf("GetConversationIdByMessageId: %w", err)
	}
	return conversationId, nil
}
