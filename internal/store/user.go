package store

// User представляет сущность user в базе данных.
type User struct {
	Username string `db:"username"` // Имя пользователя
	Password string `db:"password"` // Пароль (скрыт от JSON для безопасности)
}
