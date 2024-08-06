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
		shortUrl := Shorting()
		fullShortUrl := "/" + shortUrl
		UrlData.Add(fullShortUrl, urlStr)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("http://localhost:8888" + fullShortUrl))
	} else {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
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
