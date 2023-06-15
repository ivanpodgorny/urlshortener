package config

import (
	"encoding/json"
	"flag"
	"os"

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
	flags      *parameters
	err        error
	arguments  []string
}

type parameters struct {
	ServerAddress     string `env:"SERVER_ADDRESS" json:"server_address"`
	GRPCServerAddress string `env:"GRPC_SERVER_ADDRESS" json:"grpc_server_address"`
	BaseURL           string `env:"BASE_URL" json:"base_url"`
	FileStoragePath   string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	HMACKey           string `env:"HMAC_KEY" json:"hmac_key"`
	DatabaseDSN       string `env:"DATABASE_DSN" json:"database_dsn"`
	ConfigFile        string
	TrustedSubnet     string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	EnableHTTPS       bool   `env:"ENABLE_HTTPS" json:"enable_https"`
}

const (
	defaultServerAddress     = "localhost:8080"
	defaultGRPCServerAddress = "localhost:3200"
	defaultBaseURL           = "http://localhost:8080"
)

// NewBuilder возвращает указатель на новый экземпляр Builder.
func NewBuilder() *Builder {
	b := &Builder{
		arguments: os.Args[1:],
		parameters: &parameters{
			ServerAddress:     defaultServerAddress,
			GRPCServerAddress: defaultGRPCServerAddress,
			BaseURL:           defaultBaseURL,
		},
		flags: &parameters{},
	}
	b.prepareFlags()

	return b
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
	if err := env.Parse(b.parameters); err != nil {
		b.err = err
	}

	return b
}

// LoadFlags загружает значения флагов командной строки.
func (b *Builder) LoadFlags() *Builder {
	err := flag.CommandLine.Parse(b.arguments)
	if err != nil {
		b.err = err

		return b
	}

	if b.flags.ServerAddress != "" {
		b.parameters.ServerAddress = b.flags.ServerAddress
	}
	if b.flags.GRPCServerAddress != "" {
		b.parameters.GRPCServerAddress = b.flags.GRPCServerAddress
	}
	if b.flags.BaseURL != "" {
		b.parameters.BaseURL = b.flags.BaseURL
	}
	if b.flags.FileStoragePath != "" {
		b.parameters.FileStoragePath = b.flags.FileStoragePath
	}
	if b.flags.HMACKey != "" {
		b.parameters.HMACKey = b.flags.HMACKey
	}
	if b.flags.DatabaseDSN != "" {
		b.parameters.DatabaseDSN = b.flags.DatabaseDSN
	}
	if b.flags.EnableHTTPS {
		b.parameters.EnableHTTPS = b.flags.EnableHTTPS
	}
	if b.flags.TrustedSubnet != "" {
		b.parameters.TrustedSubnet = b.flags.TrustedSubnet
	}

	return b
}

// LoadFile загружает значения из файла конфигурации.
func (b *Builder) LoadFile() *Builder {
	if err := flag.CommandLine.Parse(b.arguments); err != nil {
		b.err = err

		return b
	}

	configFile := b.flags.ConfigFile
	configFileEnv := os.Getenv("CONFIG")
	if configFileEnv != "" {
		configFile = configFileEnv
	}

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			b.err = err

			return b
		}

		b.err = json.Unmarshal(data, &b.parameters)
	}

	return b
}

// Build возвращает Config для чтения загруженных значений параметров.
func (b *Builder) Build() (*Config, error) {
	return &Config{b.parameters}, b.err
}

func (b *Builder) prepareFlags() {
	flag.StringVar(&b.flags.ServerAddress, "a", b.parameters.ServerAddress, "адрес запуска HTTP-сервера")
	flag.StringVar(&b.flags.GRPCServerAddress, "g", b.parameters.GRPCServerAddress, "адрес запуска GRPC-сервера")
	flag.StringVar(&b.flags.BaseURL, "b", b.parameters.BaseURL, "базовый адрес результирующего сокращённого URL")
	flag.StringVar(&b.flags.FileStoragePath, "f", b.parameters.FileStoragePath, "путь к файлу для хранения сокращенных URL")
	flag.StringVar(&b.flags.DatabaseDSN, "d", b.parameters.DatabaseDSN, "адрес подключения к PostgreSQL")
	flag.BoolVar(&b.flags.EnableHTTPS, "s", b.parameters.EnableHTTPS, "включает HTTPS в веб-сервере")
	flag.StringVar(&b.flags.TrustedSubnet, "t", b.parameters.TrustedSubnet, "CIDR доверенной подсети")
	flag.StringVar(&b.flags.ConfigFile, "c", b.parameters.ConfigFile, "путь к конфигурационному файлу")
	flag.StringVar(&b.flags.ConfigFile, "config", b.parameters.ConfigFile, "путь к конфигурационному файлу")
}

// ServerAddress возвращает значение адреса HTTP-сервера.
func (c *Config) ServerAddress() string {
	return c.parameters.ServerAddress
}

// GRPCServerAddress возвращает значение адреса GRPC-сервера.
func (c *Config) GRPCServerAddress() string {
	return c.parameters.GRPCServerAddress
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

// TrustedSubnet возвращает CIDR доверенной подсети.
func (c *Config) TrustedSubnet() string {
	return c.parameters.TrustedSubnet
}
