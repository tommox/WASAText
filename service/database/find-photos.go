package database

func (db *appdbimpl) FindPhotos(user User) ([]Complete_Photo, error) {
	queryRes, err := db.c.Query("SELECT photo_id, timestamp FROM Photo WHERE owner=? ORDER BY timestamp DESC", user.User_id)

	// ! Controllo degli errori
	if err != nil {
		queryRes.Close()
		return nil, err
	}
	if queryRes.Err() != nil {
		return nil, queryRes.Err()
	}

	var photos []Complete_Photo

	for queryRes.Next() {

		// ! Prendiamo le info dal DB delle Photo
		var photo Complete_Photo
		photo.Owner = user
		nickname, err := db.FindUser(user)
		if err != nil {
			return nil, err
		}
		photo.Owner.Nickname = nickname
		err = queryRes.Scan(&photo.Photo_id, &photo.Timestamp) // Prendiamo i risultati ottenuti e li inseriamo nel singolo utente
		if err != nil {
			return nil, err
		}
	}

	return photos, nil
}
