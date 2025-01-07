package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) CheckUser(u User) (User, error) {
	// Esegui la SELECT
	row := db.c.QueryRow("SELECT User_id, Nickname FROM Users WHERE Nickname = ?;", u.Nickname)
	var found User
	err := row.Scan(&found.User_id, &found.Nickname)
	if err != nil {
		return User{}, err // Se "no rows" => sql.ErrNoRows
	}
	return found, nil
}
