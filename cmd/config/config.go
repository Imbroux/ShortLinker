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
		addr         string
		baseURL      string
		flagLogLevel string
	)

	flag.StringVar(&addr, "a", ":8888", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "8000", "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	flag.Parse()
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
