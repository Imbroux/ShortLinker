package main

import (
	"YandexLearnMiddle/cmd/config"
	"YandexLearnMiddle/database"
	"YandexLearnMiddle/internal/handlers"
	"YandexLearnMiddle/internal/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `postgres`, `625325`, `shortlinks`)
	database.InitDB(dataSourceName)

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
	r.Get("/", handlers.HandleGet())
	r.Get("/ping", handlers.GetPing())
	log.Printf("Запуск сервера на %s", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, config.GzipMiddleware(r))
}
