package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) AddReaction(messageId int, userId int, emoji string, isGroup bool) error {
	var queryCheck string
	if isGroup {
		queryCheck = `SELECT Emoji FROM GroupReactions WHERE GroupMessage_id = ? AND User_id = ?`
	} else {
		queryCheck = `SELECT Emoji FROM Reactions WHERE Message_id = ? AND User_id = ?`
	}

	var existingEmoji string
	err := db.c.QueryRow(queryCheck, messageId, userId).Scan(&existingEmoji)

	if err == nil {
		var queryUpdate string
		if isGroup {
			queryUpdate = `UPDATE GroupReactions SET Emoji = ? WHERE GroupMessage_id = ? AND User_id = ?`
		} else {
			queryUpdate = `UPDATE Reactions SET Emoji = ? WHERE Message_id = ? AND User_id = ?`
		}
		_, err := db.c.Exec(queryUpdate, emoji, messageId, userId)
		if err != nil {
			return fmt.Errorf("AddReaction: error updating reaction: %w", err)
		}
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("AddReaction: error checking existing reaction: %w", err)
	}

	var queryInsert string
	if isGroup {
		queryInsert = `INSERT INTO GroupReactions (GroupMessage_id, User_id, Emoji) VALUES (?, ?, ?)`
	} else {
		queryInsert = `INSERT INTO Reactions (Message_id, User_id, Emoji) VALUES (?, ?, ?)`
	}
	_, err = db.c.Exec(queryInsert, messageId, userId, emoji)
	if err != nil {
		return fmt.Errorf("AddReaction: error inserting new reaction: %w", err)
	}

	return nil
}
