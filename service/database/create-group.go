package database

import (
	"fmt"
	"time"
)

func (db *appdbimpl) CreateGroup(groupName string, creatorId int, createdAt time.Time) (int, error) {
	// Inserisce il nuovo gruppo nella tabella Groups
	query := `
        INSERT INTO Groups (Group_name, Creator_id, Created_at) 
        VALUES (?, ?, ?)
    `
	result, err := db.c.Exec(query, groupName, creatorId, createdAt)
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: errore nella creazione del gruppo: %w", err)
	}

	// Recupera l'ID del gruppo appena creato
	groupId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: errore nel recupero dell'ID del gruppo: %w", err)
	}

	// Inserisce nella tabella GroupConversations, includendo anche il nome del gruppo
	queryInsertGroupConversation := `
        INSERT INTO GroupConversations (Group_id, Sender_id, LastMessage_id, LastMessageTimestamp, GroupName) 
        VALUES (?, ?, 0, ?, ?)
    `
	_, err = db.c.Exec(queryInsertGroupConversation, groupId, creatorId, createdAt, groupName)
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: errore nella creazione della conversazione di gruppo: %w", err)
	}

	// Aggiunge il creatore come membro del gruppo con ruolo "admin"
	err = db.AddUserToGroup(int(groupId), creatorId, "admin")
	if err != nil {
		return 0, fmt.Errorf("CreateGroup: errore nell'aggiunta del creatore al gruppo: %w", err)
	}

	// Ritorna l'ID del gruppo creato
	return int(groupId), nil
}
