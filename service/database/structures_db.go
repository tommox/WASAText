package database

// Username
type UserId struct {
	userId string `json:"userId"`
}

// User structure
type User struct {
	userId   int    `json:"userId"`
	userName string `json:"userName"`
}
