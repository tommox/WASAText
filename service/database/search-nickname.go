package database

func (db *appdbimpl) SearchNickname(nickname string) (bool, error) {

	var exist int
	err := db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE nickname = ?)", nickname).Scan(&exist)
	if exist == 1 {
		return true, err
	}
	return false, err
}
