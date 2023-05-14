package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

type Config interface {
	ServerAddress() string
	BaseURL() string
	FileStoragePath() string
	HMACKey() string
	DatabaseDSN() string
}

type Builder struct {
	parameters *parameters
	err        error
}

type parameters struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	HMACKey         string `env:"HMAC_KEY"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

const (
	defaultServerAddress = "localhost:8080"
	defaultBaseURL       = "http://localhost:8080"
)

func NewBuilder() *Builder {
	return &Builder{
		parameters: &parameters{
			ServerAddress: defaultServerAddress,
			BaseURL:       defaultBaseURL,
		},
	}
}

func (b *Builder) SetDefaultServerAddress(addr string) *Builder {
	b.parameters.ServerAddress = addr

	return b
}

func (b *Builder) SetDefaultBaseURL(url string) *Builder {
	b.parameters.BaseURL = url

	return b
}

func (b *Builder) LoadEnv() *Builder {
	b.err = env.Parse(b.parameters)

	return b
}

func (b *Builder) LoadFlags() *Builder {
	flag.StringVar(&b.parameters.ServerAddress, "a", b.parameters.ServerAddress, "адрес запуска HTTP-сервера")
	flag.StringVar(&b.parameters.BaseURL, "b", b.parameters.BaseURL, "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&b.parameters.FileStoragePath, "f", "", "путь к файлу для хранения сокращенных URL")
	flag.StringVar(&b.parameters.DatabaseDSN, "d", "", "адрес подключения к PostgreSQL")
	flag.Parse()

	return b
}

func (b *Builder) Build() (Config, error) {
	return b, b.err
}

func (b *Builder) ServerAddress() string {
	return b.parameters.ServerAddress
}

func (b *Builder) BaseURL() string {
	return b.parameters.BaseURL
}

func (b *Builder) FileStoragePath() string {
	return b.parameters.FileStoragePath
}

func (b *Builder) HMACKey() string {
	return b.parameters.HMACKey
}

func (b *Builder) DatabaseDSN() string {
	return b.parameters.DatabaseDSN
}
