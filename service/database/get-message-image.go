package database

func (db *appdbimpl) GetMessageImage(messageId int) (int, []byte, string, bool, *int, error) {
	row := db.c.QueryRow("SELECT Sender_id, ImageData, Timestamp, IsRead, IsReply FROM Messages WHERE Message_id = ?", messageId)

	var senderId int
	var imageData []byte
	var timestamp string
	var isRead bool
	var isReply *int

	err := row.Scan(&senderId, &imageData, &timestamp, &isRead, &isReply)
	if err != nil {
		return 0, nil, "", false, nil, err
	}

	return senderId, imageData, timestamp, isRead, isReply, nil
}

// Aggiorna la funzione GetGroupMessageImage per includere Sender_id, IsRead e altre informazioni
func (db *appdbimpl) GetGroupMessageImage(messageId int) (int, []byte, string, bool, *int, error) {
	row := db.c.QueryRow("SELECT Sender_id, ImageData, Timestamp, IsRead, IsReply FROM GroupMessages WHERE GroupMessage_id = ?", messageId)

	var senderId int
	var imageData []byte
	var timestamp string
	var isRead bool
	var isReply *int

	err := row.Scan(&senderId, &imageData, &timestamp, &isRead, &isReply)
	if err != nil {
		return 0, nil, "", false, nil, err
	}

	return senderId, imageData, timestamp, isRead, isReply, nil
}
