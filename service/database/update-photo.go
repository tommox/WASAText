package database

func (db *appdbimpl) UpdateUserPhoto(userID int, photoData []byte) error {
	_, err := db.c.Exec(`
        UPDATE Users 
           SET photo = ?
         WHERE user_id = ?;`,
		photoData, userID)
	return err
}
