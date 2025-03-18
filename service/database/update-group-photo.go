package database

func (db *appdbimpl) UpdateGroupPhoto(groupId int, photoData []byte) error {
	_, err := db.c.Exec(`
        UPDATE Groups 
           SET Photo = ?
         WHERE Group_id = ?;`,
		photoData, groupId)
	return err
}
