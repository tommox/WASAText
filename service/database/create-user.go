package database

// Method to add a user in DB
func (db *appdbimpl) CreateUser(u User) error {

	// INSERT into DB
	_, err := db.c.Exec("INSERT INTO Users(Nickname) VALUES (?)", u.Nickname)

	// Error management
	if err != nil {
		return err
	}

	return nil

}
