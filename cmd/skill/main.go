package main

import (
	"YandexLearnMiddle/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", handlers.HandlePost)
	r.Get("/*", handlers.HandleGet)

	return http.ListenAndServe(`:8080`, r)
}
