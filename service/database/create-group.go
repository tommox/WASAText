package database

import (
	"fmt"
)

func (db *appdbimpl) CreateGroup(name string, creatorId int) (int, error) {
	query := `INSERT INTO Groups (Group_name, Creator_id) VALUES (?, ?)`
	result, err := db.c.Exec(query, name, creatorId)
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: %w", err)
	}

	groupId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: could not retrieve group ID: %w", err)
	}

	// Aggiungi il creator come admin
	err = db.AddUserToGroup(int(groupId), creatorId, "admin")
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: failed to add creator as admin: %w", err)
	}

	return int(groupId), nil
}
