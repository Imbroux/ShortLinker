package services

import "YandexLearnMiddle/internal/db"

func CreateUserTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        full_url TEXT NOT NULL,
        short_url TEXT NOT NULL
    );`
	
	_, err := db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
