package database

import "fmt"

func (db *appdbimpl) GetGroupById(groupId int) (Group, error) {
	var group Group
	query := `SELECT Group_id, Group_name, Creator_id, Created_at FROM Groups WHERE Group_id = ?`
	err := db.c.QueryRow(query, groupId).Scan(&group.Group_id, &group.Group_name, &group.Creator_id, &group.Created_at)
	if err != nil {
		return group, fmt.Errorf("GetGroupById: %w", err)
	}
	return group, nil
}
