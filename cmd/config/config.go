package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr         string
	BaseURL      string
	FlagLogLevel string
}

func NewConfig() (*Config, error) {
	var (
		addr            string
		baseURL         string
		flagLogLevel    string
		fileStoragePath string
	)
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "File path to store URL data")
	flag.StringVar(&addr, "a", ":8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "8000", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	flag.Parse()
	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		addr = envRunAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}

	return &Config{
		Addr:         addr,
		BaseURL:      baseURL,
		FlagLogLevel: flagLogLevel,
	}, nil
}
