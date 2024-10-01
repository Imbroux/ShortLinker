package mock

import (
	"YandexLearnMiddle/internal/store"
)

// MockLinkService - мок-реализация интерфейса LinkService для целей тестирования.
type MockLinkService struct {
	// Функция для генерации сокращённой ссылки.
	ShortenLinkFunc func() string
	// Функция для сохранения ссылки.
	SaveLinkFunc func(link store.Link) error
	// Функция для получения оригинальной ссылки по сокращённой.
	GetOriginalLinkFunc func(shortLink, username string) (string, error)
	// Функция для получения всех ссылок пользователя.
	GetAllLinksFunc func(userID int) ([]store.Link, error)
	// Функция для удаления ссылок.
	DeleteLinksFunc func(userID int, shortLinks []string) error
}

// GenerateShortLink вызывает мок-функцию для генерации сокращённой ссылки.
func (m *MockLinkService) GenerateShortLink() string {
	return m.ShortenLinkFunc()
}

// SaveLink вызывает мок-функцию для сохранения ссылки.
func (m *MockLinkService) SaveLink(link store.Link) error {
	return m.SaveLinkFunc(link)
}

// GetOriginalLink вызывает мок-функцию для получения оригинальной ссылки по сокращённой.
func (m *MockLinkService) GetOriginalLink(shortLink, username string) (string, error) {
	return m.GetOriginalLinkFunc(shortLink, username)
}

// GetAllLinks вызывает мок-функцию для получения всех ссылок пользователя.
func (m *MockLinkService) GetAllLinks(userID int) ([]store.Link, error) {
	return m.GetAllLinksFunc(userID)
}

// DeleteLinks вызывает мок-функцию для удаления ссылок.
func (m *MockLinkService) DeleteLinks(userID int, shortLinks []string) error {
	return m.DeleteLinksFunc(userID, shortLinks)
}
