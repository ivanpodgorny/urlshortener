package main

import (
	"flag"
	"github.com/caarlos0/env/v7"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	HMACKey         string `env:"HMAC_KEY"`
}

const (
	DefaultServerAddress = "localhost:8080"
	DefaultBaseURL       = "http://localhost:8080"
)

func LoadConfig() (Config, error) {
	cfg := Config{}

	flag.StringVar(&cfg.ServerAddress, "a", DefaultServerAddress, "адрес запуска HTTP-сервера")
	flag.StringVar(&cfg.BaseURL, "b", DefaultBaseURL, "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "путь к файлу для хранения сокращенных URL")
	flag.Parse()

	err := env.Parse(&cfg)

	return cfg, err
}
