package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var flagParameters = &parameters{}

func TestBuilder_LoadEnv(t *testing.T) {
	var (
		serverAddress   = "localhost:8080"
		baseURL         = "http://localhost:8080"
		fileStoragePath = "/path"
		hmacKey         = "key"
		databaseDSN     = "dsn"
		enableHTTPS     = "true"
		builder         = &Builder{
			parameters: &parameters{},
		}
	)

	require.NoError(t, os.Setenv("SERVER_ADDRESS", serverAddress))
	require.NoError(t, os.Setenv("BASE_URL", baseURL))
	require.NoError(t, os.Setenv("FILE_STORAGE_PATH", fileStoragePath))
	require.NoError(t, os.Setenv("HMAC_KEY", hmacKey))
	require.NoError(t, os.Setenv("DATABASE_DSN", databaseDSN))
	require.NoError(t, os.Setenv("ENABLE_HTTPS", enableHTTPS))

	cfg, err := builder.LoadEnv().Build()
	require.NoError(t, err)
	assert.Equal(t, serverAddress, cfg.ServerAddress())
	assert.Equal(t, baseURL, cfg.BaseURL())
	assert.Equal(t, fileStoragePath, cfg.FileStoragePath())
	assert.Equal(t, hmacKey, cfg.HMACKey())
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN())
	assert.True(t, cfg.EnableHTTPS())
}

func TestBuilder_LoadFile(t *testing.T) {
	var (
		serverAddress   = "localhost:8080"
		baseURL         = "http://localhost:8080"
		fileStoragePath = "/path"
		hmacKey         = "key"
		databaseDSN     = "dsn"
	)
	f, err := os.CreateTemp("", "config.json")
	require.NoError(t, err)
	configParameters := &parameters{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		HMACKey:         hmacKey,
		DatabaseDSN:     databaseDSN,
		EnableHTTPS:     true,
	}
	p, err := json.Marshal(configParameters)
	require.NoError(t, err)
	_, err = f.Write(p)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	builder := &Builder{
		parameters: &parameters{},
		flags:      flagParameters,
		arguments: []string{
			"-c", f.Name(),
		},
	}
	builder.prepareFlags()
	cfg, err := builder.LoadFile().Build()
	require.NoError(t, err)
	assert.Equal(t, serverAddress, cfg.ServerAddress())
	assert.Equal(t, baseURL, cfg.BaseURL())
	assert.Equal(t, fileStoragePath, cfg.FileStoragePath())
	assert.Equal(t, hmacKey, cfg.HMACKey())
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN())
	assert.True(t, cfg.EnableHTTPS())

	_ = os.Remove(f.Name())
}

func TestBuilder_LoadFlags(t *testing.T) {
	var (
		serverAddress   = "localhost:8080"
		baseURL         = "http://localhost:8080"
		fileStoragePath = "/path"
		databaseDSN     = "dsn"
		builder         = &Builder{
			parameters: &parameters{},
			flags:      flagParameters,
			arguments: []string{
				"-a", serverAddress,
				"-b", baseURL,
				"-f", fileStoragePath,
				"-d", databaseDSN,
				"-s",
			},
		}
	)

	cfg, err := builder.LoadFlags().Build()
	require.NoError(t, err)
	assert.Equal(t, serverAddress, cfg.ServerAddress())
	assert.Equal(t, baseURL, cfg.BaseURL())
	assert.Equal(t, fileStoragePath, cfg.FileStoragePath())
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN())
	assert.True(t, cfg.EnableHTTPS())
}
