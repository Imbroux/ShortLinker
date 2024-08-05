package config

import (
	"flag"
)

type Config struct {
	Addr    string
	BaseURL string
}

func NewConfig() (*Config, error) {
	var addr string
	var baseURL string

	flag.StringVar(&addr, "a", ":8888", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "8000", "Базовый адрес результирующего сокращённого URL")

	flag.Parse()

	return &Config{
		Addr:    addr,
		BaseURL: baseURL,
	}, nil
}
