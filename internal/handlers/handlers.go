package handlers

import (
	"YandexLearnMiddle/database"
	"YandexLearnMiddle/internal/maps"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var (
	filePath = "urls.json"
	UrlData  = maps.New()
)

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()) // Простой UUID на основе времени
}
func saveToFile(data URLData) error {
	var existingData []URLData

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&existingData); err != nil && err.Error() != "EOF" {
		return err
	}

	existingData = append(existingData, data)
	file.Seek(0, 0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingData); err != nil {
		return err
	}

	return nil
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
	fullShortUrl := "/" + shortUrl

	urlData := URLData{
		UUID:        generateUUID(),
		ShortURL:    fullShortUrl,
		OriginalURL: req.URL,
	}

	if err := saveToFile(urlData); err != nil {
		http.Error(w, "Unable to save data", http.StatusInternalServerError)
		return
	}

	res := struct {
		Result string `json:"result"`
	}{
		Result: "http://localhost:8080" + fullShortUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
func HandleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := r.URL.Path
		if shortUrl == "/" {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		if value, exists := UrlData.Get(shortUrl); exists {
			w.Header().Set("Location", value)
			_, _ = w.Write([]byte(value))

			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "Short URL not found", http.StatusNotFound)
		}
	}
}
func GetPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := database.DB.Ping(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
