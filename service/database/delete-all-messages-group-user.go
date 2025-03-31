package database

import "fmt"

// DeleteAllMessagesFromUserInGroup elimina tutti i messaggi di un utente in un gruppo
func (db *appdbimpl) DeleteAllMessagesFromUserInGroup(groupId int, userId int) error {
	query := `SELECT GroupMessage_id FROM GroupMessages WHERE Group_id = ? AND Sender_id = ?`
	rows, err := db.c.Query(query, groupId, userId)
	if err != nil {
		return fmt.Errorf("DeleteAllMessagesFromUserInGroup: errore nella query: %w", err)
	}
	defer rows.Close()

	var messageIds []int
	for rows.Next() {
		var messageId int
		if err := rows.Scan(&messageId); err != nil {
			return fmt.Errorf("DeleteAllMessagesFromUserInGroup: errore nello scan: %w", err)
		}
		messageIds = append(messageIds, messageId)
	}

	for _, msgId := range messageIds {
		err := db.DeleteGroupMessage(msgId)
		if err != nil {
			return fmt.Errorf("DeleteAllMessagesFromUserInGroup: errore nel cancellare messaggio %d: %w", msgId, err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("DeleteAllMessagesFromUserInGroup: errore iterando le righe: %w", err)
	}
	return nil
}
