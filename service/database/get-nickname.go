package database

func (db *appdbimpl) GetNickname(user_id int) (string, error) {

	var user User
	err := db.c.QueryRow("SELECT * FROM Users WHERE user_id=?", user_id).Scan(&user.User_id, &user.Nickname)
	return user.Nickname, err
}
