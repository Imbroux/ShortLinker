package service

import (
	"YandexLearnMiddle/internal/store"
	"YandexLearnMiddle/postgresql"
	"database/sql"
	"fmt"
	"math/rand"
)

type LinkService interface {
	GenerateShortLink() string
	SaveLink(link store.Link) error
	GetOriginalLink(shortLink string, username string) (string, error)
	GetAllLinks(userID int) ([]store.Link, error)
}

type linkService struct{}

func NewLinkService() LinkService {
	return &linkService{}
}

func (s *linkService) GenerateShortLink() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortLink := make([]byte, 8)
	for i := range shortLink {
		shortLink[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortLink)
}

func (s *linkService) SaveLink(link store.Link) error {
	_, err := postgresql.DB.Exec("INSERT INTO links (original, shortened, user_id) VALUES ($1, $2, $3)", link.Original, link.Shortened, link.UserID)
	return err
}

func (s *linkService) GetOriginalLink(shortLink string, username string) (string, error) {
	var originalLink string
	err := postgresql.DB.QueryRow("SELECT original FROM links WHERE shortened = $1 AND user_id = (SELECT id FROM users WHERE username = $2)", shortLink, username).Scan(&originalLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("link not found")
		}
		return "", err
	}
	return originalLink, nil
}

func (s *linkService) GetAllLinks(userID int) ([]store.Link, error) {
	rows, err := postgresql.DB.Query("SELECT original, shortened FROM links WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []store.Link
	for rows.Next() {
		var link store.Link
		if err := rows.Scan(&link.Original, &link.Shortened); err != nil {
			return nil, err
		}
		link.UserID = userID
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
