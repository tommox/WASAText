package database

func (db *appdbimpl) GetMessageImage(messageId int) (int, []byte, string, bool, *int, bool, error) {
	row := db.c.QueryRow("SELECT Sender_id, ImageData, Timestamp, IsRead, IsReply, IsForward FROM Messages WHERE Message_id = ?", messageId)

	var senderId int
	var imageData []byte
	var timestamp string
	var isRead bool
	var isReply *int
	var isForward bool

	err := row.Scan(&senderId, &imageData, &timestamp, &isRead, &isReply, &isForward)
	if err != nil {
		return 0, nil, "", false, nil, false, err
	}

	return senderId, imageData, timestamp, isRead, isReply, isForward, nil
}

// Aggiorna la funzione GetGroupMessageImage per includere Sender_id, IsRead e altre informazioni
func (db *appdbimpl) GetGroupMessageImage(messageId int) (int, []byte, string, bool, *int, bool, error) {
	row := db.c.QueryRow("SELECT Sender_id, ImageData, Timestamp, IsRead, IsReply, IsForward FROM GroupMessages WHERE GroupMessage_id = ?", messageId)

	var senderId int
	var imageData []byte
	var timestamp string
	var isRead bool
	var isReply *int
	var isForward bool

	err := row.Scan(&senderId, &imageData, &timestamp, &isRead, &isReply, &isForward)
	if err != nil {
		return 0, nil, "", false, nil, false, err
	}

	return senderId, imageData, timestamp, isRead, isReply, isForward, nil
}
