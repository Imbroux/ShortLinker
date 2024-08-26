package services

import "YandexLearnMiddle/internal/db"

func CreateUrlsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        full_url TEXT UNIQUE NOT NULL,
        short_url TEXT NOT NULL
    );`

	_, err := db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
