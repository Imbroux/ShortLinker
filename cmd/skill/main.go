package main

import (
	"YandexLearnMiddle/internal/handlers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Webhook))
}
