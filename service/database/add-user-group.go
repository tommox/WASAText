package database

import "fmt"

func (db *appdbimpl) AddUserToGroup(groupId int, userId int, role string) error {
	query := `
        INSERT INTO GroupMembers (Group_id, User_id, Role)
        VALUES (?, ?, ?);
    `
	_, err := db.c.Exec(query, groupId, userId, role)
	if IsUniqueConstraintError(err) {
		return ErrUserAlreadyInGroup
	}
	if err != nil {
		return fmt.Errorf("AddUserToGroup: %w", err)
	}
	return nil
}
