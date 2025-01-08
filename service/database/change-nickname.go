package database

import "fmt"

func (db *appdbimpl) ChangeNickname(u User, newName string) error {
	// UPDATE in table `Users`.
	// Let's assum that the colomn autoincrement is called `userId` and the nickname one is called `nickname`.

	query := `UPDATE Users SET Nickname = ? WHERE User_id = ?;`
	_, err := db.c.Exec(query, newName, u.User_id)
	if err != nil {
		return fmt.Errorf("ChangeNickname: %w", err)
	}

	return nil
}
