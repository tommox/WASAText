package database

func (db *appdbimpl) CheckUser(u User) (User, error) {
	// SELECT
	row := db.c.QueryRow("SELECT User_id, Nickname FROM Users WHERE User_id = ?;", u.User_id)
	var found User
	err := row.Scan(&found.User_id, &found.Nickname)
	if err != nil {
		return User{}, err // If "no rows" => sql.ErrNoRows
	}
	return found, nil
}
