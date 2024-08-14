package handlers

import (
	"YandexLearnMiddle/internal/maps"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var UrlData = maps.New()

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
	UrlData.Add(fullShortUrl, req.URL)
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
