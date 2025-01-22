package database

import "fmt"

func (db *appdbimpl) PromoteToAdmin(groupId int, userId int) error {
	query := `UPDATE GroupMembers SET Role = 'admin' WHERE Group_id = ? AND User_id = ?`
	result, err := db.c.Exec(query, groupId, userId)
	if err != nil {
		return fmt.Errorf("PromoteToAdmin: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("PromoteToAdmin: user not found in group or already admin")
	}
	return nil
}
