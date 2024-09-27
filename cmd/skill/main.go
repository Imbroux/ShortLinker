package main

import (
	"YandexLearnMiddle/internal/handler"
	"YandexLearnMiddle/postgresql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
	}

	// Проверка переменных окружения
	requiredEnvVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME", "SIGNING_KEY"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			logger.Fatal(fmt.Sprintf("Environment variable %s is required", envVar))
		}
	}

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	postgresql.InitDB(dataSourceName, logger)

	handler.Logger = logger // Передаём logger в handler
	run(logger)             // Передаем logger в run
}

func run(log *zap.Logger) {
	r := handler.InitRouters()

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error starting server", zap.Error(err))
	}
}
