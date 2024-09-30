package service

import (
	"YandexLearnMiddle/internal/store"
	"YandexLearnMiddle/postgresql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"math/rand"
)

type LinkService interface {
	GenerateShortLink() string
	SaveLink(link store.Link) error
	GetOriginalLink(shortLink string, username string) (string, error)
	GetAllLinks(userID int) ([]store.Link, error)
	MarkLinksAsDeleted(shortLinks []string, userID int) error
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

func (s *linkService) GetAllLinks(userID int) ([]store.Link, error) {
	rows, err := postgresql.DB.Query("SELECT original, shortened, is_deleted FROM links WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []store.Link
	for rows.Next() {
		var link store.Link
		if err := rows.Scan(&link.Original, &link.Shortened, &link.DeletedFlag); err != nil {
			return nil, err
		}
		link.UserID = userID
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (s *linkService) GetOriginalLink(shortLink string, username string) (string, error) {
	var originalLink string
	var deletedFlag bool
	err := postgresql.DB.QueryRow(
		"SELECT original, is_deleted FROM links WHERE shortened = $1 AND user_id = (SELECT id FROM users WHERE username = $2)",
		shortLink, username,
	).Scan(&originalLink, &deletedFlag)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("link not found")
		}
		return "", err
	}

	if deletedFlag {
		return "", fmt.Errorf("link deleted")
	}

	return originalLink, nil
}

func (s *linkService) MarkLinksAsDeleted(shortLinks []string, userID int) error {
	query := "UPDATE links SET is_deleted = true WHERE shortened = ANY($1) AND user_id = $2"
	_, err := postgresql.DB.Exec(query, pq.Array(shortLinks), userID)
	return err
}
