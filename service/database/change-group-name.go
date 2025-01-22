package database

import "fmt"

func (db *appdbimpl) ChangeGroupName(groupId int, newGroupName string) error {
	query := `UPDATE Groups SET Group_name = ? WHERE Group_id = ?`
	result, err := db.c.Exec(query, newGroupName, groupId)
	if err != nil {
		return fmt.Errorf("ChangeGroupName: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrGroupNotFound
	}

	return nil
}
