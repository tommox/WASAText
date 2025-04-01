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
	Message_id      int       `json:"message_id"`
	Sender_id       int       `json:"sender_id"`
	Conversation_id int       `json:"conversation_id"`
	MessageContent  string    `json:"message_content"`
	ImageData       []byte    `json:"image_data"`
	Timestamp       time.Time `json:"timestamp"`
	IsRead          bool      `json:"isRead"`
	IsReply         *int      `json:"isReply,omitempty"`
}

// GroupMessage rappresenta la struttura di un messaggio inviata a un gruppo
type GroupMessage struct {
	GroupMessage_id int       `json:"message_id"`
	Sender_id       int       `json:"sender_id"`
	Group_id        int       `json:"group_id"`
	MessageContent  string    `json:"message_content"`
	ImageData       []byte    `json:"image_data"`
	Timestamp       time.Time `json:"timestamp"`
	IsRead          bool      `json:"isRead"`
	IsReply         *int      `json:"isReply,omitempty"`
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

type Reaction struct {
	Emoji   string `json:"emoji"`
	User_id int    `json:"user_id"`
}

// Conversation rappresenta una conversazione di un utente con un altro.
type Conversation struct {
	Conversation_id      int       `json:"conversation_id"`
	Sender_id            int       `json:"sender_id"`
	Recipient_id         int       `json:"recipient_id"`
	LastMessage_id       int       `json:"last_message_id"`
	LastMessageTimestamp time.Time `json:"last_message_timestamp"`
	LastMessageIsRead    bool      `json:"last_message_isRead"`
	LastMessageSenderId  int       `json:"last_message_sender_id"`
}

// Conversation rappresenta una conversazione di un utente in un gruppo.
type GroupConversation struct {
	GroupConversation_id int       `json:"groupconversation_id"`
	Group_id             int       `json:"group_id"`
	Sender_id            int       `json:"sender_id"`
	LastMessage_id       int       `json:"last_message_id"`
	LastMessageTimestamp time.Time `json:"last_message_timestamp"`
	LastMessageIsRead    bool      `json:"last_message_isRead"`
	LastMessageSenderId  int       `json:"last_message_sender_id"`
}

func (u User) toDataBase() database.User {
	return database.User{
		User_id:  u.User_id,
		Nickname: u.Nickname,
	}
}

func toDatabaseMessage(dbMsg database.Message) Message {
	return Message{
		Message_id:      dbMsg.Message_id,
		Sender_id:       dbMsg.Sender_id,
		Conversation_id: dbMsg.Conversation_id,
		MessageContent:  dbMsg.MessageContent,
		Timestamp:       dbMsg.Timestamp,
		IsRead:          dbMsg.IsRead,
		IsReply:         dbMsg.IsReply,
	}
}
