package database

func (db *appdbimpl) GetMessageImage(messageId int) ([]byte, error) {
	row := db.c.QueryRow("SELECT ImageData FROM Messages WHERE Message_id = ?", messageId)
	var imageData []byte
	err := row.Scan(&imageData)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

func (db *appdbimpl) GetMessageGroupImage(groupId int) ([]byte, error) {
	row := db.c.QueryRow("SELECT ImageData FROM GroupMessages WHERE Group_id = ?", groupId)
	var imageData []byte
	err := row.Scan(&imageData)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
