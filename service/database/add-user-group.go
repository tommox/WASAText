package database

import "fmt"

func (db *appdbimpl) AddUserToGroup(groupId int, userId int) error {
	query := `
        INSERT INTO GroupMembers (Group_id, User_id)
        VALUES (?, ?);
    `
	_, err := db.c.Exec(query, groupId, userId)
	if IsUniqueConstraintError(err) {
		return ErrUserAlreadyInGroup
	}
	if err != nil {
		return fmt.Errorf("AddUserToGroup: %w", err)
	}
	return nil
}
