package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) MarkConversationAsRead(conversationId int, userId int) error {
	// Aggiorniamo i messaggi ricevuti dall'utente come letti
	_, err := db.c.Exec(`
		UPDATE Messages
		SET IsRead = TRUE
		WHERE Conversation_id = ? AND Sender_id != ?;
	`, conversationId, userId)
	if err != nil {
		return fmt.Errorf("MarkConversationAsRead: %w", err)
	}

	// Controlliamo se l'ultimo messaggio è stato inviato da qualcun altro
	var senderId int
	err = db.c.QueryRow(`
		SELECT LastMessageSenderId FROM Conversations WHERE Conversation_id = ?;
	`, conversationId).Scan(&senderId)
	if err != nil {
		return fmt.Errorf("MarkConversationAsRead - get sender: %w", err)
	}

	// Solo se l'utente è il destinatario, aggiorniamo il flag di lettura
	if senderId != userId {
		_, err = db.c.Exec(`
			UPDATE Conversations
			SET LastMessageIsRead = TRUE
			WHERE Conversation_id = ?;
		`, conversationId)
		if err != nil {
			return fmt.Errorf("MarkConversationAsRead - update conversation flag: %w", err)
		}
	}

	return nil
}

func (db *appdbimpl) MarkGroupConversationAsRead(groupId int, userId int) error {
	// 1. Recuperiamo tutti i messaggi del gruppo non ancora letti da questo utente
	rows, err := db.c.Query(`
		SELECT gm.GroupMessage_id, gm.Sender_id
		FROM GroupMessages gm
		LEFT JOIN GroupMessagesRead gmr ON gm.GroupMessage_id = gmr.GroupMessage_id AND gmr.User_id = ?
		WHERE gm.Group_id = ? AND gmr.User_id IS NULL;
	`, userId, groupId)
	if err != nil {
		return fmt.Errorf("MarkGroupConversationAsRead - select unread messages: %w", err)
	}
	defer rows.Close()

	var messagesToMark []struct {
		MessageId int
		SenderId  int
	}
	for rows.Next() {
		var m struct {
			MessageId int
			SenderId  int
		}
		err := rows.Scan(&m.MessageId, &m.SenderId)
		if err != nil {
			return fmt.Errorf("MarkGroupConversationAsRead - scan: %w", err)
		}
		// Evitiamo di registrare la lettura se l'utente è il mittente
		if m.SenderId != userId {
			messagesToMark = append(messagesToMark, m)
		}
	}

	// 2. Per ogni messaggio, aggiungiamo la lettura nella tabella di tracciamento
	for _, m := range messagesToMark {
		_, err := db.c.Exec(`
			INSERT OR IGNORE INTO GroupMessagesRead (GroupMessage_id, User_id) VALUES (?, ?);
		`, m.MessageId, userId)
		if err != nil {
			return fmt.Errorf("MarkGroupConversationAsRead - insert read: %w", err)
		}
	}

	// 3. Verifica se l’ultimo messaggio del gruppo è stato letto da tutti (tranne il mittente)
	var lastMessageId, senderId sql.NullInt64
	err = db.c.QueryRow(`
		SELECT LastMessage_id, LastMessageSenderId FROM GroupConversations WHERE Group_id = ?;
	`, groupId).Scan(&lastMessageId, &senderId)
	if err != nil {
		return fmt.Errorf("MarkGroupConversationAsRead - get last message: %w", err)
	}
	if !lastMessageId.Valid || !senderId.Valid {
		// Nessun messaggio ancora nel gruppo
		return nil
	}

	var totalMembers int
	err = db.c.QueryRow(`
		SELECT COUNT(*) FROM GroupMembers WHERE Group_id = ? AND User_id != ?;
	`, groupId, senderId.Int64).Scan(&totalMembers)
	if err != nil {
		return fmt.Errorf("MarkGroupConversationAsRead - count members: %w", err)
	}

	var readCount int
	err = db.c.QueryRow(`
		SELECT COUNT(*) FROM GroupMessagesRead WHERE GroupMessage_id = ?;
	`, lastMessageId.Int64).Scan(&readCount)
	if err != nil {
		return fmt.Errorf("MarkGroupConversationAsRead - count reads: %w", err)
	}

	if readCount >= totalMembers {
		_, err := db.c.Exec(`
			UPDATE GroupConversations SET LastMessageIsRead = TRUE WHERE Group_id = ?;
		`, groupId)
		if err != nil {
			return fmt.Errorf("MarkGroupConversationAsRead - update conversation: %w", err)
		}
	}

	// 4. Se tutti hanno letto un messaggio, aggiorniamo anche GroupMessages.IsRead
	for _, m := range messagesToMark {
		var readers int
		err := db.c.QueryRow(`
			SELECT COUNT(*) FROM GroupMessagesRead WHERE GroupMessage_id = ?;
		`, m.MessageId).Scan(&readers)
		if err != nil {
			return fmt.Errorf("MarkGroupConversationAsRead - count readers: %w", err)
		}

		var total int
		err = db.c.QueryRow(`
			SELECT COUNT(*) FROM GroupMembers WHERE Group_id = ? AND User_id != ?;
		`, groupId, m.SenderId).Scan(&total)
		if err != nil {
			return fmt.Errorf("MarkGroupConversationAsRead - count members for IsRead: %w", err)
		}

		if readers >= total {
			_, err := db.c.Exec(`
				UPDATE GroupMessages SET IsRead = TRUE WHERE GroupMessage_id = ?;
			`, m.MessageId)
			if err != nil {
				return fmt.Errorf("MarkGroupConversationAsRead - update GroupMessages IsRead: %w", err)
			}
		}
	}

	return nil
}
