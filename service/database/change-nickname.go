package database

import "fmt"

func (db *appdbimpl) ChangeNickname(u User, newName string) error {
	// Semplice UPDATE nella tabella `Users`.
	// Assumiamo che la colonna autoincrement si chiami `userId` e quella del nickname `nickname`.

	query := `UPDATE Users SET Nickname = ? WHERE User_id = ?;`
	_, err := db.c.Exec(query, newName, u.User_id)
	if err != nil {
		return fmt.Errorf("ChangeNickname: %w", err)
	}

	return nil
}
