package handlers

import (
	"YandexLearnMiddle/internal/db"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func IsValidUrl(urlStr string) bool {
	return strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")
}

func Shorting() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func saveToDatabase(data URLData) error {
	query := `
  INSERT INTO users (full_url, short_url) VALUES ($1, $2);
 `
	_, err := db.DB.Exec(query, data.OriginalURL, data.ShortURL)
	if err != nil {
		return err
	}
	return nil
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil || !IsValidUrl(req.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	shortUrl := Shorting()

	urlData := URLData{
		ShortURL:    shortUrl,
		OriginalURL: req.URL,
	}

	if err := saveToDatabase(urlData); err != nil {
		http.Error(w, "Unable to save data", http.StatusInternalServerError)
		return
	}

	res := struct {
		Result string `json:"result"`
	}{
		Result: "http://localhost:8080/" + shortUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func HandleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := chi.URLParam(r, "shortURL")

		if shortUrl == "" {
			log.Println("Short URL is empty")
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		query := `
   SELECT full_url FROM users WHERE short_url = $1;
  `
		var fullUrl string
		row := db.DB.QueryRow(query, shortUrl)
		err := row.Scan(&fullUrl)
		if err != nil {
			log.Println("Error querying database:", err)
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", fullUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(fullUrl))
	}
}

func GetPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.DB.Ping(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
