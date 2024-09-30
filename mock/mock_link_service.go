package mock

import (
	"YandexLearnMiddle/internal/store"
)

// MockLinkService - мок-реализация интерфейса LinkService
type MockLinkService struct {
	ShortenLinkFunc     func() string
	SaveLinkFunc        func(link store.Link) error
	GetOriginalLinkFunc func(shortLink, username string) (string, error)
	GetAllLinksFunc     func(userID int) ([]store.Link, error)
	DeleteLinksFunc     func(userID int, shortLinks []string) error
}

func (m *MockLinkService) GenerateShortLink() string {
	return m.ShortenLinkFunc()
}

func (m *MockLinkService) SaveLink(link store.Link) error {
	return m.SaveLinkFunc(link)
}

func (m *MockLinkService) GetOriginalLink(shortLink, username string) (string, error) {
	return m.GetOriginalLinkFunc(shortLink, username)
}

func (m *MockLinkService) GetAllLinks(userID int) ([]store.Link, error) {
	return m.GetAllLinksFunc(userID)
}

func (m *MockLinkService) DeleteLinks(userID int, shortLinks []string) error {
	return m.DeleteLinksFunc(userID, shortLinks)
}
