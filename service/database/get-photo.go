package database

func (db *appdbimpl) GetUserPhoto(userID int) ([]byte, error) {
	var photoData []byte

	err := db.c.QueryRow(`SELECT photo FROM Users WHERE user_id = ?`, userID).Scan(&photoData)
	if err != nil {
		return nil, err
	}
	return photoData, nil
}

func (db *appdbimpl) GetGroupPhoto(groupId int) ([]byte, error) {
	var photoData []byte

	err := db.c.QueryRow(`SELECT photo FROM Groups WHERE Group_id = ?`, groupId).Scan(&photoData)
	if err != nil {
		return nil, err
	}
	return photoData, nil
}
