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
	Message_id      int       `json:"message_id"`
	Sender_id       int       `json:"sender_id"`
	Conversation_id int       `json:"conversation_id"`
	MessageContent  string    `json:"message_content"`
	ImageData       []byte    `json:"image_data"`
	Timestamp       time.Time `json:"timestamp"`
}

// GroupMessage rappresenta la struttura di un messaggio inviata a un gruppo
type GroupMessage struct {
	GroupMessage_id int       `json:"message_id"`
	Sender_id       int       `json:"sender_id"`
	Group_id        int       `json:"group_id"`
	MessageContent  string    `json:"message_content"`
	ImageData       []byte    `json:"image_data"`
	Timestamp       time.Time `json:"timestamp"`
}

type Reaction struct {
	Emoji   string `json:"emoji"`
	User_id int    `json:"user_id"`
}

// Group rappresenta un gruppo nel sistema.
type Group struct {
	Group_id   int       `json:"group_id"`
	Group_name string    `json:"group_name"`
	Creator_id int       `json:"creator_id"`
	Created_at time.Time `json:"created_at"`
}

// GroupMember rappresenta un membro di un gruppo.
type GroupMember struct {
	GroupMember_id int    `json:"group_member_id"`
	Group_id       int    `json:"group_id"`
	User_id        int    `json:"user_id"`
	Role           string `json:"role"`
	Nickname       string `json:"nickname"`
}

// Conversation rappresenta una conversazione di un utente con un altro.
type Conversation struct {
	Conversation_id      int       `json:"conversation_id"`
	Sender_id            int       `json:"sender_id"`
	Recipient_id         int       `json:"recipient_id"`
	LastMessage_id       int       `json:"last_message_id"`
	LastMessageTimestamp time.Time `json:"last_message_timestamp"`
}

// Conversation rappresenta una conversazione di un utente in un gruppo.
type GroupConversation struct {
	GroupConversation_id int       `json:"group_conversation_id"`
	Group_id             int       `json:"group_id"`
	GroupName            string    `json:"group_name"`
	Sender_id            int       `json:"sender_id"`
	LastMessage_id       int       `json:"last_message_id"`
	LastMessageTimestamp time.Time `json:"last_message_timestamp"`
}
