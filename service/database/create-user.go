package database

// Metodo per aggiungere un utente nel DB
func (db *appdbimpl) CreateUser(u User) error {

	// Eseguiamo una INSERT nel DB
	_, err := db.c.Exec("INSERT INTO Users(Nickname) VALUES (?)", u.Nickname)

	// Gestione degli errori
	if err != nil {
		return err
	}

	return nil

}
