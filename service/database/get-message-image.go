package database

// Aggiorna la funzione GetMessageImage per includere IsRead
// Aggiorna la funzione GetMessageImage per includere Sender_id, IsRead e altre informazioni
func (db *appdbimpl) GetMessageImage(messageId int) (int, []byte, string, bool, error) {
	row := db.c.QueryRow("SELECT Sender_id, ImageData, Timestamp, IsRead FROM Messages WHERE Message_id = ?", messageId)

	var senderId int
	var imageData []byte
	var timestamp string
	var isRead bool

	err := row.Scan(&senderId, &imageData, &timestamp, &isRead)
	if err != nil {
		return 0, nil, "", false, err
	}

	return senderId, imageData, timestamp, isRead, nil
}

// Aggiorna la funzione GetGroupMessageImage per includere IsRead
// Aggiorna la funzione GetGroupMessageImage per includere Sender_id, IsRead e altre informazioni
func (db *appdbimpl) GetGroupMessageImage(messageId int) (int, []byte, string, bool, error) {
	row := db.c.QueryRow("SELECT Sender_id, ImageData, Timestamp, IsRead FROM GroupMessages WHERE GroupMessage_id = ?", messageId)

	var senderId int
	var imageData []byte
	var timestamp string
	var isRead bool

	err := row.Scan(&senderId, &imageData, &timestamp, &isRead)
	if err != nil {
		return 0, nil, "", false, err
	}

	return senderId, imageData, timestamp, isRead, nil
}
