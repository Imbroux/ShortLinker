package store

// Link представляет сущность ссылки в базе данных.
type Link struct {
	ID          int    `db:"id"`         // Уникальный идентификатор ссылки.
	Original    string `db:"original"`   // Исходная ссылка, которую нужно сократить.
	Shortened   string `db:"shortened"`  // Сокращённая версия исходной ссылки.
	UserID      int    `db:"user_id"`    // Идентификатор пользователя, который создал ссылку.
	DeletedFlag bool   `db:"is_deleted"` // Флаг, указывающий на то, была ли ссылка удалена.
}
