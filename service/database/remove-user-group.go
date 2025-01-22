package database

import "fmt"

func (db *appdbimpl) RemoveUserFromGroup(groupId int, userId int) error {
	query := `DELETE FROM GroupMembers WHERE Group_id = ? AND User_id = ?`
	result, err := db.c.Exec(query, groupId, userId)
	if err != nil {
		return fmt.Errorf("RemoveUserFromGroup: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUserNotInGroup
	}
	return nil
}
