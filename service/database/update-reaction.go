package database

import (
	"fmt"
	"strings"
)

func (db *appdbimpl) UpdateMessageReaction(messageId int, emoji string, add bool) error {
	// Recupera il campo Reactions del messaggio
	var reactions string
	err := db.c.QueryRow(`SELECT Reactions FROM Messages WHERE Message_id = ?`, messageId).Scan(&reactions)
	if err != nil {
		return fmt.Errorf("UpdateMessageReaction: error retrieving reactions: %w", err)
	}

	// Aggiungi o rimuovi l'emoji
	if add {
		reactions += emoji
	} else {
		reactions = strings.ReplaceAll(reactions, emoji, "")
	}

	// Aggiorna il campo nel database
	_, err = db.c.Exec(`UPDATE Messages SET Reactions = ? WHERE Message_id = ?`, reactions, messageId)
	if err != nil {
		return fmt.Errorf("UpdateMessageReaction: error updating reactions: %w", err)
	}

	return nil
}
