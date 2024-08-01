package handlers

import (
	"YandexLearnMiddle/internal/maps"
	"io"
	"math/rand"
	"net/http"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var UrlData = maps.New()

func IsValidUrl(token string) bool {
	u, err := url.Parse(token)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func Shorting() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	urlStr := string(body)

	if IsValidUrl(urlStr) {
		shortUrl := "http://localhost:8080/" + Shorting()
		UrlData.Add(shortUrl, urlStr)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortUrl))
	} else {
		http.Error(w, "Invalid URL", http.StatusBadGateway)
		return
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path

	if value, exists := UrlData.Get("http://localhost:8080" + shortUrl); exists {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, _ = w.Write([]byte("Location: " + value))
	} else {
		http.Error(w, "Short URL not found", http.StatusNotFound)
	}
}
