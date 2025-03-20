package database

import "fmt"

func (db *appdbimpl) GetGroupByMessageId(messageId int) (Group, error) {
	var group Group
	query := `
		SELECT g.Group_id, g.Group_name, g.Creator_id, g.Created_at
		FROM Groups g
		JOIN GroupMessages gm ON g.Group_id = gm.Group_id
		WHERE gm.GroupMessage_id = ?`
	err := db.c.QueryRow(query, messageId).Scan(&group.Group_id, &group.Group_name, &group.Creator_id, &group.Created_at)
	if err != nil {
		return group, fmt.Errorf("GetGroupByMessageId: %w", err)
	}
	return group, nil
}
