package database

import "fmt"

// GetReactionsForMessage recupera tutte le reazioni per un messaggio specifico.
func (db *appdbimpl) GetReactionsForMessage(messageId int, isGroup bool) ([]Reaction, error) {
	var query string
	if isGroup {
		query = `SELECT Emoji, User_id FROM GroupReactions WHERE GroupMessage_id = ?`
	} else {
		query = `SELECT Emoji, User_id FROM Reactions WHERE Message_id = ?`
	}

	rows, err := db.c.Query(query, messageId)
	if err != nil {
		return nil, fmt.Errorf("GetReactionsForMessage: %w", err)
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var reaction Reaction
		if err := rows.Scan(&reaction.Emoji, &reaction.User_id); err != nil {
			return nil, fmt.Errorf("GetReactionsForMessage: error scanning row: %w", err)
		}
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetReactionsForMessage: error iterating rows: %w", err)
	}
	return reactions, nil
}
