package config

import (
	"flag"

	"github.com/caarlos0/env/v7"
)

// Config хранит значения параметров приложения и позволяет получить
// их через геттеры.
type Config struct {
	parameters *parameters
}

// Builder реализует методы для загрузки значений параметров.
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
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
}

const (
	defaultServerAddress = "localhost:8080"
	defaultBaseURL       = "http://localhost:8080"
)

// NewBuilder возвращает указатель на новый экземпляр Builder.
func NewBuilder() *Builder {
	return &Builder{
		parameters: &parameters{
			ServerAddress: defaultServerAddress,
			BaseURL:       defaultBaseURL,
		},
	}
}

// SetDefaultServerAddress устанавливает значение адреса сервера по умолчанию.
func (b *Builder) SetDefaultServerAddress(addr string) *Builder {
	b.parameters.ServerAddress = addr

	return b
}

// SetDefaultBaseURL устанавливает значение базового URL сокращенных ссылок по умолчанию.
func (b *Builder) SetDefaultBaseURL(url string) *Builder {
	b.parameters.BaseURL = url

	return b
}

// LoadEnv загружает значения переменных окружения.
func (b *Builder) LoadEnv() *Builder {
	b.err = env.Parse(b.parameters)

	return b
}

// LoadFlags загружает значения флагов командной строки.
func (b *Builder) LoadFlags() *Builder {
	flag.StringVar(&b.parameters.ServerAddress, "a", b.parameters.ServerAddress, "адрес запуска HTTP-сервера")
	flag.StringVar(&b.parameters.BaseURL, "b", b.parameters.BaseURL, "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&b.parameters.FileStoragePath, "f", "", "путь к файлу для хранения сокращенных URL")
	flag.StringVar(&b.parameters.DatabaseDSN, "d", "", "адрес подключения к PostgreSQL")
	flag.BoolVar(&b.parameters.EnableHTTPS, "s", false, "включает HTTPS в веб-сервере")
	flag.Parse()

	return b
}

// Build возвращает Config для чтения загруженных значений параметров.
func (b *Builder) Build() (*Config, error) {
	return &Config{b.parameters}, b.err
}

// ServerAddress возвращает значение адреса сервера.
func (c *Config) ServerAddress() string {
	return c.parameters.ServerAddress
}

// BaseURL возвращает значение базового URL сокращенных ссылок.
func (c *Config) BaseURL() string {
	return c.parameters.BaseURL
}

// FileStoragePath возвращает путь к файлу для хранения сокращенных URL.
func (c *Config) FileStoragePath() string {
	return c.parameters.FileStoragePath
}

// HMACKey возвращает значение ключа для создания HMAC подписи.
func (c *Config) HMACKey() string {
	return c.parameters.HMACKey
}

// DatabaseDSN возвращает строку подключения к PostgreSQL.
func (c *Config) DatabaseDSN() string {
	return c.parameters.DatabaseDSN
}

// EnableHTTPS возвращает значение флага включения HTTPS в веб-сервере.
func (c *Config) EnableHTTPS() bool {
	return c.parameters.EnableHTTPS
}
