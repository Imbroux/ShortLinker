package store

type Link struct {
	ID          int    `db:"id"`
	Original    string `db:"original"`
	Shortened   string `db:"shortened"`
	UserID      int    `db:"user_id"`
	DeletedFlag bool   `db:"is_deleted"`
}
