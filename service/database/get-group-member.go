package database

import "fmt"

func (db *appdbimpl) GetGroupMembers(groupId int) ([]GroupMember, error) {
	query := `SELECT GroupMember_id, Group_id, User_id, Role FROM GroupMembers WHERE Group_id = ?`
	rows, err := db.c.Query(query, groupId)
	if err != nil {
		return nil, fmt.Errorf("GetGroupMembers: %w", err)
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var member GroupMember
		err := rows.Scan(&member.GroupMember_id, &member.Group_id, &member.User_id, &member.Role)
		if err != nil {
			return nil, fmt.Errorf("GetGroupMembers: %w", err)
		}
		members = append(members, member)
	}

	return members, nil
}
