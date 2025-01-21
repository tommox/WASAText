package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) AddReaction(messageId int, userId int, emoji string) error {
	// Controlla se esiste già una reazione per questo utente e messaggio
	queryCheck := `SELECT Emoji FROM Reactions WHERE Message_id = ? AND User_id = ?`
	var existingEmoji string
	err := db.c.QueryRow(queryCheck, messageId, userId).Scan(&existingEmoji)

	if err == nil {
		// Reazione trovata: aggiorna l'emoji
		queryUpdate := `UPDATE Reactions SET Emoji = ? WHERE Message_id = ? AND User_id = ?`
		_, err := db.c.Exec(queryUpdate, emoji, messageId, userId)
		if err != nil {
			return fmt.Errorf("AddReaction: error updating reaction: %w", err)
		}
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("AddReaction: error checking existing reaction: %w", err)
	}

	// Nessuna reazione esistente: aggiungi una nuova reazione
	queryInsert := `INSERT INTO Reactions (Message_id, User_id, Emoji) VALUES (?, ?, ?)`
	_, err = db.c.Exec(queryInsert, messageId, userId, emoji)
	if err != nil {
		return fmt.Errorf("AddReaction: error inserting new reaction: %w", err)
	}

	return nil
}
