package database

import "fmt"

func (db *appdbimpl) UpdateGroupPhoto(groupId int, photoData []byte) error {
	query := `UPDATE Groups SET Photo = ? WHERE Group_id = ?`
	result, err := db.c.Exec(query, photoData, groupId)
	if err != nil {
		return fmt.Errorf("UpdateGroupPhoto: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrGroupNotFound
	}

	return nil
}
