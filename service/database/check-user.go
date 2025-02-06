package database

func (db *appdbimpl) CheckUser(u User) (User, error) {
	// SELECT
	row := db.c.QueryRow("SELECT User_id, Nickname FROM Users WHERE Nickname = ?;", u.Nickname)
	var found User
	err := row.Scan(&found.User_id, &found.Nickname)
	if err != nil {
		return User{}, err
	}
	return found, nil
}
