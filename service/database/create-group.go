package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateGroup(name string, creatorId int, createdAt time.Time) (int, error) {
	query := `INSERT INTO Groups (Group_name, Creator_id, Created_at) VALUES (?, ?, ?)`
	result, err := db.c.Exec(query, name, creatorId, createdAt)
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
