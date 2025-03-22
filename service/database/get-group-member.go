package database

import "fmt"

func (db *appdbimpl) GetGroupMembers(groupId int) ([]GroupMember, error) {
	query := `
        SELECT gm.GroupMember_id, gm.Group_id, gm.User_id, gm.Role, u.Nickname
        FROM GroupMembers gm
        INNER JOIN Users u ON gm.User_id = u.User_id
        WHERE gm.Group_id = ?`

	rows, err := db.c.Query(query, groupId)
	if err != nil {
		return nil, fmt.Errorf("GetGroupMembers: %w", err)
	}
	defer rows.Close()

	var members []GroupMember
	for rows.Next() {
		var member GroupMember
		err := rows.Scan(&member.GroupMember_id, &member.Group_id, &member.User_id, &member.Role, &member.Nickname)
		if err != nil {
			return nil, fmt.Errorf("GetGroupMembers scan: %w", err)
		}
		members = append(members, member)
	}

	return members, nil
}
