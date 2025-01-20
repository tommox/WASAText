package api

import (
	"time"

	"github.com/tommox/WASAText/service/database"
)

// Nickname
type Nickname struct {
	Nickname string `json:"nickname"`
}

// Struttura User
type User struct {
	User_id  int    `json:"user_id"`
	Nickname string `json:"nickname"`
}

// Struttura Photo
type Complete_Photo struct {
	Photo_id  int       `json:"photo_id"`
	Owner     User      `json:"owner"`
	Timestamp time.Time `json:"timestamp"`
}

// Message rappresenta la struttura di un messaggio inviata al client
type Message struct {
	Message_id     int       `json:"message_id"`
	Sender_id      int       `json:"sender_id"`
	Recipient_id   int       `json:"recipient_id"`
	MessageContent string    `json:"message_content"`
	Timestamp      time.Time `json:"timestamp"`
}

func (u User) toDataBase() database.User {
	return database.User{
		User_id:  u.User_id,
		Nickname: u.Nickname,
	}
}

func toDatabaseMessage(dbMsg database.Message) Message {
	return Message{
		Message_id:     dbMsg.Message_id,
		Sender_id:      dbMsg.Sender_id,
		Recipient_id:   dbMsg.Recipient_id,
		MessageContent: dbMsg.MessageContent,
		Timestamp:      dbMsg.Timestamp,
	}
}
