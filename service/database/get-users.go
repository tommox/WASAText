package database

import "fmt"

func (db *appdbimpl) GetUsers() ([]User, error) {
	query := `
        SELECT u.User_id, u.Nickname 
        FROM Users u
        INNER JOIN Users a ON u.User_id = a.User_id`

	rows, err := db.c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.User_id, &user.Nickname)
		if err != nil {
			return nil, fmt.Errorf("GetUsers: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUsers: errore iterando le righe: %w", err)
	}
	return users, nil
}
