/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Users
	CreateUser(User) error
	ChangeNickname(User, string) error

	FindUserId(User) (int, error)
	CheckUser(User) (User, error)
	CheckUserId(u User) (User, error)
	SearchNickname(string) (bool, error)
	UpdateUserPhoto(userID int, photoData []byte) error
	GetUsers() ([]User, error)
	GetUserPhoto(userID int) ([]byte, error)

	// Messages
	CreateMessage(senderId int, conversationId int, messageContent string, timestamp time.Time, isReply *int) (int, error)
	GetMessage(messageId int) (Message, error)
	DeleteMessage(messageId int) error
	GetMessageImage(messageId int) (int, []byte, string, bool, *int, error)
	GetReactionsForMessage(messageId int, isGroup bool) ([]Reaction, error)
	MarkIsReply(messageId int, replyMessageId int) error
	CreateImageMessage(senderId int, conversationId int, imageContent []byte, timestamp time.Time) (int, error)
	UpdateOrCreateConversation(sender int, recipient int, messageId int, timestamp time.Time, isRead bool, senderMsgId int) (int, error)
	// Reactions
	AddReaction(messageId int, userId int, emoji string, isGroup bool) error
	RemoveReaction(messageId int, userId int, isGroup bool) error

	// Groups
	CreateGroup(groupName string, creatorId int, createdAt time.Time) (int, error)
	AddUserToGroup(groupId int, userId int, role string) error
	RemoveUserFromGroup(groupId int, userId int) error
	PromoteToAdmin(groupId int, userId int) error
	DeleteGroupMessage(messageId int) error
	GetGroupByMessageId(messageId int) (Group, error)
	GetGroupMembers(groupId int) ([]GroupMember, error)
	IsGroupAdmin(groupId int, userId int) (bool, error)
	IsGroupMember(groupId int, userId int) (bool, error)
	DeleteGroup(groupId int) error
	GetGroupMessage(groupId, messageId int) (GroupMessage, error)
	ChangeGroupName(groupId int, newGroupName string) error
	GetGroupPhoto(groupId int) ([]byte, error)
	GetGroupMessageImage(messageId int) (int, []byte, string, bool, *int, error)
	UpdateGroupPhoto(groupId int, photoData []byte) error
	CreateGroupImageMessage(groupId int, senderId int, imageData []byte, timestamp time.Time) (int, error)
	CreateGroupMessage(groupId int, senderId int, messageContent string, timestamp time.Time, isReply *int) (int, error)
	updateOrCreateGroupConversation(groupconversation_id int, groupId int, senderId int, messageId int, timestamp time.Time) error
	DeleteAllMessagesFromUserInGroup(groupId int, userId int) error
	MarkIsReplyGroup(messageId int, replyMessageId int) error

	// Conversations
	CheckPrivateConversationAccess(userId, conversationId int) (bool, error)
	CheckGroupConversationAccess(userId, groupConversationId int) (bool, error)
	GetConversationMessages(conversationId int) ([]Message, error)
	GetGroupConversationMessages(groupConversationId int) ([]GroupMessage, error)
	GetUserConversations(userId int) (map[string]interface{}, error)
	CheckExistingConversation(userId int, recipientId int) (int, error)
	GetConversationIdByMessageId(messageId int) (int, error)
	DeleteConversation(conversationId int) error
	MarkConversationAsRead(conversationId int, userId int) error
	MarkGroupConversationAsRead(groupId int, userId int) error

	// Authorization
	CheckUserPermission(userId, messageId int) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		// Creating DB for Users if not existing
		users := `CREATE TABLE IF NOT EXISTS Users 
									   (User_id  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
										Nickname VARCHAR(16) NOT NULL UNIQUE,
										Photo    BLOB);`
		_, err = db.Exec(users)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: Users %w", err)
		}

		// Creating DB for Messages if not existing
		messages := `CREATE TABLE IF NOT EXISTS Messages
									   (Message_id       INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
									    Sender_id  		 INTEGER NOT NULL,
										Conversation_id  INTEGER NOT NULL, 
										ImageData		 BLOB,
										MessageContent   VARCHAR(1000),
										Timestamp        DATETIME NOT NULL,
										IsRead           BOOLEAN DEFAULT FALSE,
										IsReply 		 INTEGER DEFAULT NULL);`

		_, err = db.Exec(messages)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: Messages %w", err)

		}

		// Creating DB for Reactions if not existing
		reactions := `CREATE TABLE IF NOT EXISTS Reactions 
								   (Reaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
									Message_id INTEGER NOT NULL,            
									User_id INTEGER NOT NULL,                     
									Emoji TEXT NOT NULL,
									UNIQUE(Message_id, User_id),
									FOREIGN KEY (Message_id) REFERENCES Messages (Message_id) ON DELETE CASCADE,
									FOREIGN KEY (User_id) REFERENCES Users (User_id));`
		_, err = db.Exec(reactions)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: Reactions %w", err)
		}

		// Creating DB for Groups if not existing
		groups := `CREATE TABLE IF NOT EXISTS Groups 
								   (Group_id INTEGER PRIMARY KEY AUTOINCREMENT,
									Group_name TEXT NOT NULL,
									Creator_id INTEGER NOT NULL,
									Created_at DATETIME NOT NULL,
									Photo      BLOB,
									FOREIGN KEY (Creator_id) REFERENCES Users (User_id) ON DELETE CASCADE);`

		_, err = db.Exec(groups)
		if err != nil {
			return nil, fmt.Errorf("error creating Groups structure: %w", err)
		}

		// Creating DB for GroupMembers if not existing
		groupMembers := `CREATE TABLE IF NOT EXISTS GroupMembers
								   (GroupMember_id INTEGER PRIMARY KEY AUTOINCREMENT,
									Group_id INTEGER NOT NULL,
									User_id INTEGER NOT NULL,
									Role TEXT NOT NULL CHECK (Role IN ('member', 'admin')),
									UNIQUE(Group_id, User_id),
									FOREIGN KEY (Group_id) REFERENCES Groups (Group_id) ON DELETE CASCADE,
									FOREIGN KEY (User_id) REFERENCES Users (User_id));`

		_, err = db.Exec(groupMembers)
		if err != nil {
			return nil, fmt.Errorf("error creating GroupMembers structure: %w", err)
		}

		// Creating DB for GroupMessages if not existing
		groupMessages := `CREATE TABLE IF NOT EXISTS GroupMessages
								   (GroupMessage_id INTEGER PRIMARY KEY AUTOINCREMENT,
									Group_id INTEGER NOT NULL,
									Sender_id INTEGER NOT NULL,
									MessageContent VARCHAR(1000),
									ImageData BLOB,
									Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
									IsRead BOOLEAN DEFAULT FALSE,
									IsReply INTEGER DEFAULT NULL,
									FOREIGN KEY (Group_id) REFERENCES Groups (Group_id) ON DELETE CASCADE,
									FOREIGN KEY (Sender_id) REFERENCES Users (User_id) ON DELETE CASCADE);`
		_, err = db.Exec(groupMessages)
		if err != nil {
			return nil, fmt.Errorf("error creating GroupMessages structure: %w", err)
		}

		// Creating DB for Reactions if not existing
		groupReactions := `CREATE TABLE IF NOT EXISTS GroupReactions 
									(Reaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
									GroupMessage_id INTEGER NOT NULL,             
									User_id INTEGER NOT NULL,                     
									Emoji TEXT NOT NULL,
									UNIQUE(GroupMessage_id, User_id), 
									FOREIGN KEY (GroupMessage_id) REFERENCES GroupMessages (GroupMessage_id) ON DELETE CASCADE,  -- Correzione qui
									FOREIGN KEY (User_id) REFERENCES Users (User_id));`
		_, err = db.Exec(groupReactions)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: GroupReactions %w", err)
		}

		groupMessagesRead := `CREATE TABLE IF NOT EXISTS GroupMessagesRead (
									GroupMessage_id INTEGER NOT NULL,
									User_id INTEGER NOT NULL,
									PRIMARY KEY (GroupMessage_id, User_id),
									FOREIGN KEY (GroupMessage_id) REFERENCES GroupMessages(GroupMessage_id) ON DELETE CASCADE,
									FOREIGN KEY (User_id) REFERENCES Users(User_id) ON DELETE CASCADE
								);`
		_, err = db.Exec(groupMessagesRead)
		if err != nil {
			return nil, fmt.Errorf("error creating GroupMessagesRead structure: GroupMessagesRead %w", err)
		}

		// Creating DB for Conversations if not existing
		conversations := `CREATE TABLE IF NOT EXISTS Conversations
								   (Conversation_id INTEGER PRIMARY KEY AUTOINCREMENT,
								    Sender_id INTEGER NOT NULL,
									Recipient_id INTEGER NOT NULL,
									LastMessage_id INTEGER DEFAULT NULL,
									LastMessageTimestamp DATETIME DEFAULT NULL,
									LastMessageIsRead BOOLEAN DEFAULT FALSE,
									LastMessageSenderId INTEGER DEFAULT NULL,
									FOREIGN KEY (Sender_id) REFERENCES Users (User_id) ON DELETE CASCADE,
									FOREIGN KEY (Recipient_id) REFERENCES Users (User_id) ON DELETE CASCADE,
									FOREIGN KEY (LastMessage_id) REFERENCES Messages (Message_id) ON DELETE SET NULL,
									UNIQUE(Sender_id, Recipient_id));`
		_, err = db.Exec(conversations)
		if err != nil {
			return nil, fmt.Errorf("error creating Conversations structure: %w", err)
		}

		// Creating DB for GroupConversation if not existing
		groupConversations := `CREATE TABLE IF NOT EXISTS GroupConversations
								   (GroupConversation_id INTEGER PRIMARY KEY AUTOINCREMENT,
									Group_id INTEGER NOT NULL,
									GroupName VAR CHAR(1000),
									Sender_id INTEGER NOT NULL,
									LastMessage_id INTEGER DEFAULT NULL,
									LastMessageTimestamp DATETIME DEFAULT NULL,
									LastMessageIsRead BOOLEAN DEFAULT FALSE,
									LastMessageSenderId INTEGER DEFAULT NULL,
									FOREIGN KEY (Group_id) REFERENCES Groups (Group_id) ON DELETE CASCADE,
									FOREIGN KEY (Sender_id) REFERENCES Users (User_id) ON DELETE SET NULL,
									FOREIGN KEY (LastMessage_id) REFERENCES GroupMessages (GroupMessage_id) ON DELETE SET NULL,
									UNIQUE(Group_id));`
		_, err = db.Exec(groupConversations)
		if err != nil {
			return nil, fmt.Errorf("error creating GroupConversations structure: %w", err)
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
