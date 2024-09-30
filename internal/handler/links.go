package handler

import (
	"YandexLearnMiddle/internal/service"
	"YandexLearnMiddle/internal/store"
	"YandexLearnMiddle/postgresql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LinkHandler interface {
	ShortenLink(w http.ResponseWriter, r *http.Request)
	GetOriginalLink(w http.ResponseWriter, r *http.Request)
	GetAllLinks(w http.ResponseWriter, r *http.Request)
	DeleteLinks(w http.ResponseWriter, r *http.Request)
}

type linkHandler struct {
	linkService service.LinkService
}

func NewLinkHandler(service service.LinkService) LinkHandler {
	return &linkHandler{linkService: service}
}

func (h *linkHandler) ShortenLink(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var link store.Link
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	shortLink := h.linkService.GenerateShortLink()

	userID, err := getUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	link.Shortened = shortLink
	link.UserID = userID

	if err := h.linkService.SaveLink(link); err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (h *linkHandler) GetOriginalLink(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	shortLink := r.URL.Query().Get("shortened")
	if shortLink == "" {
		http.Error(w, "Missing shortened link parameter", http.StatusBadRequest)
		return
	}

	originalLink, err := h.linkService.GetOriginalLink(shortLink, username)
	if err != nil {
		if err.Error() == "link not found" {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"original": originalLink,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *linkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := getUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	links, err := h.linkService.GetAllLinks(userID)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *linkHandler) DeleteLinks(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok || username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var shortLinks []string
	if err := json.NewDecoder(r.Body).Decode(&shortLinks); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDByUsername(username)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	go func() {
		if err := h.linkService.MarkLinksAsDeleted(shortLinks, userID); err != nil {
			log.Printf("Failed to delete links: %v", err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

func getUserIDByUsername(username string) (int, error) {
	var userID int
	err := postgresql.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}
	return userID, nil
}
