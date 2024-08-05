package main

import (
	"YandexLearnMiddle/cmd/config"
	"YandexLearnMiddle/internal/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфигурации: %v", err)
	}

	if err := run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func run(cfg *config.Config) error {
	fmt.Println("Запуск сервера на", cfg.Addr)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", handlers.HandlePost)
	r.Get("/*", handlers.HandleGet(cfg))

	return http.ListenAndServe(cfg.Addr, r)
}
