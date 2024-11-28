package api

// Username
type userName struct {
	userName string `json:"userName"`
}

// User structure
type User struct {
	userId   int    `json:"userId"`
	userName string `json:"userName"`
}
