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
	SearchNickname(string) (bool, error)
	UpdateUserPhoto(userID int, photoData []byte) error

	// Messages
	CreateMessage(senderId int, recipientId int, messageContent string, timestamp time.Time) (int, error)
	GetMessage(messageId int) (Message, error)

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
									    Sender_id 		 INTEGER NOT NULL,
										Recipient_id     INTEGER NOT NULL, 
										messageContent   VARCHAR(1000) NOT NULL,
										timestamp        DATETIME NOT NULL);`

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
									Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP, 
									FOREIGN KEY (Message_id) REFERENCES Messages (Message_id) ON DELETE CASCADE,
									FOREIGN KEY (User_id) REFERENCES Users (User_id) ON DELETE CASCADE);`
		_, err = db.Exec(reactions)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: Reactions %w", err)

		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
