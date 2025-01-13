package database

import "time"

// User structure
type User struct {
	User_id  int    `json:"User_id"`
	Nickname string `json:"nickname"`
}

// Struttura Photo
type Complete_Photo struct {
	Photo_id  int       `json:"photo_id"`
	Owner     User      `json:"owner"`
	Timestamp time.Time `json:"timestamp"`
}
