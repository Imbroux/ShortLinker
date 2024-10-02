package main

import (
	"YandexLearnMiddle/internal/handler"
	"YandexLearnMiddle/postgresql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	err := godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
	}

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

	go func() {
		logger.Info("Starting pprof server on :6060")
		logErr := http.ListenAndServe("localhost:6060", nil)
		if logErr != nil {
			logger.Fatal("Error starting pprof server", zap.Error(logErr))
		}
	}()

	handler.Logger = logger
	run(logger)
}

func run(log *zap.Logger) {
	s := handler.NewServer()

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", s.Router); err != nil {
		log.Fatal("Error starting server", zap.Error(err))
	}
}
