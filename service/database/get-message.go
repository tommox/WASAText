package database

import "fmt"

func (db *appdbimpl) GetMessage(messageId int) (Message, error) {
	var msg Message

	// Recupera i dettagli principali del messaggio
	query := `
        SELECT Message_id, Sender_id, Conversation_id, MessageContent, ImageData, Timestamp, IsRead, IsReply, IsForward
        FROM Messages
        WHERE Message_id = ?;
    `
	err := db.c.QueryRow(query, messageId).Scan(
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
		return msg, fmt.Errorf("GetMessage: error retrieving message: %w", err)
	}
	return msg, nil
}

// GetGroupMessage recupera un messaggio di gruppo specifico
func (db *appdbimpl) GetGroupMessage(groupId, messageId int) (GroupMessage, error) {
	var msg GroupMessage

	// Recupera i dettagli principali del messaggio di gruppo
	query := `
        SELECT GroupMessage_id, Sender_id, Group_id, MessageContent, ImageData, Timestamp, IsRead, IsReply, IsForward
        FROM GroupMessages
        WHERE Group_id = ? AND GroupMessage_id = ?;
    `
	err := db.c.QueryRow(query, groupId, messageId).Scan(
		&msg.GroupMessage_id,
		&msg.Sender_id,
		&msg.Group_id,
		&msg.MessageContent,
		&msg.ImageData,
		&msg.Timestamp,
		&msg.IsRead,
		&msg.IsReply,
		&msg.IsForward,
	)
	if err != nil {
		return msg, fmt.Errorf("GetGroupMessage: error retrieving group message: %w", err)
	}
	return msg, nil
}
