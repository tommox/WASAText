package database

func (db *appdbimpl) FindUser(user User) (string, error) {
	var nickname string
	err := db.c.QueryRow("SELECT nickname FROM Users WHERE user_id=?", user.User_id).Scan(&nickname)
	if err != nil {
		return "", err
	}
	return nickname, nil
}
