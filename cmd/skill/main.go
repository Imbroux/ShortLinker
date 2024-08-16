package main

import (
	"YandexLearnMiddle/cmd/config"
	"YandexLearnMiddle/internal/handlers"
	"YandexLearnMiddle/internal/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфигурации: %v", err)
	}

	sugar, err := logger.InitLogger()
	if err != nil {
		log.Fatalf("Ошибка при инициализации логгера: %v", err)
	}
	defer sugar.Sync()

	logger.Sugar = sugar

	if err := run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func run(cfg *config.Config) error {
	r := chi.NewRouter()
	r.Use(logger.WithLogging)

	r.Post("/api/shorten", handlers.HandlePost)
	r.Get("/*", handlers.HandleGet())

	log.Printf("Запуск сервера на %s", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, config.GzipMiddleware(r))
}
