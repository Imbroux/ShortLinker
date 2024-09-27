package store

type Link struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
	UserID    int    `json:"user_id"`
}
