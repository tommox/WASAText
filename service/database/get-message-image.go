package database

func (db *appdbimpl) GetMessageImage(messageId int) ([]byte, string, error) {
	row := db.c.QueryRow("SELECT ImageData, Timestamp FROM Messages WHERE Message_id = ?", messageId)

	var imageData []byte
	var timestamp string

	err := row.Scan(&imageData, &timestamp)
	if err != nil {
		return nil, "", err
	}

	return imageData, timestamp, nil
}

func (db *appdbimpl) GetGroupMessageImage(messageId int) ([]byte, string, error) {
	row := db.c.QueryRow("SELECT ImageData, Timestamp FROM GroupMessages WHERE GroupMessage_id = ?", messageId)

	var imageData []byte
	var timestamp string

	err := row.Scan(&imageData, &timestamp)
	if err != nil {
		return nil, "", err
	}

	return imageData, timestamp, nil
}
