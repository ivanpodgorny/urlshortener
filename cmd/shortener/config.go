package main

import (
	"github.com/caarlos0/env/v7"
)

type Config struct {
	ServerAddress   string  `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string  `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)

	return cfg, err
}
