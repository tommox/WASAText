package database

import "fmt"

// Funzione per aggiornare il campo "IsReply" in un messaggio esistente
func (db *appdbimpl) MarkIsReply(messageId int, replyMessageId int) error {
	query := `UPDATE Messages SET IsReply = ? WHERE Message_id = ?`
	_, err := db.c.Exec(query, replyMessageId, messageId)
	if err != nil {
		return fmt.Errorf("MarkConversationAsRead: %w", err)
	}
	return nil
}

// Funzione per aggiornare il campo "IsReply" in un messaggio esistente
func (db *appdbimpl) MarkIsReplyGroup(messageId int, replyMessageId int) error {
	query := `UPDATE GroupMessages SET IsReply = ? WHERE GroupMessage_id = ?`
	_, err := db.c.Exec(query, replyMessageId, messageId)
	if err != nil {
		return fmt.Errorf("MarkConversationAsRead: %w", err)
	}
	return nil
}
