package database

import "time"

// User structure
type User struct {
	User_id  int    `json:"User_id"`
	Nickname string `json:"Nickname"`
}

// Struttura Photo
type Complete_Photo struct {
	Photo_id  int       `json:"photo_id"`
	Owner     User      `json:"owner"`
	Timestamp time.Time `json:"timestamp"`
}

// Message structure
type Message struct {
	Message_id     int       `json:"message_id"`
	Sender_id      int       `json:"sender_id"`
	Recipient_id   int       `json:"recipient_id"`
	MessageContent string    `json:"message_content"`
	Timestamp      time.Time `json:"timestamp"`
}
