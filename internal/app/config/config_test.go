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
		serverAddress     = "localhost:8080"
		grpcServerAddress = "localhost:50051"
		baseURL           = "http://localhost:8080"
		fileStoragePath   = "/path"
		hmacKey           = "key"
		databaseDSN       = "dsn"
		enableHTTPS       = "true"
		trustedSubnet     = "192.168.0.0/24"
		builder           = &Builder{
			parameters: &parameters{},
		}
	)

	require.NoError(t, os.Setenv("SERVER_ADDRESS", serverAddress))
	require.NoError(t, os.Setenv("GRPC_SERVER_ADDRESS", grpcServerAddress))
	require.NoError(t, os.Setenv("BASE_URL", baseURL))
	require.NoError(t, os.Setenv("FILE_STORAGE_PATH", fileStoragePath))
	require.NoError(t, os.Setenv("HMAC_KEY", hmacKey))
	require.NoError(t, os.Setenv("DATABASE_DSN", databaseDSN))
	require.NoError(t, os.Setenv("ENABLE_HTTPS", enableHTTPS))
	require.NoError(t, os.Setenv("TRUSTED_SUBNET", trustedSubnet))

	cfg, err := builder.LoadEnv().Build()
	require.NoError(t, err)
	assert.Equal(t, serverAddress, cfg.ServerAddress())
	assert.Equal(t, grpcServerAddress, cfg.GRPCServerAddress())
	assert.Equal(t, baseURL, cfg.BaseURL())
	assert.Equal(t, fileStoragePath, cfg.FileStoragePath())
	assert.Equal(t, hmacKey, cfg.HMACKey())
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN())
	assert.True(t, cfg.EnableHTTPS())
	assert.Equal(t, trustedSubnet, cfg.TrustedSubnet())
}

func TestBuilder_LoadFile(t *testing.T) {
	var (
		serverAddress     = "localhost:8080"
		grpcServerAddress = "localhost:50051"
		baseURL           = "http://localhost:8080"
		fileStoragePath   = "/path"
		hmacKey           = "key"
		databaseDSN       = "dsn"
		trustedSubnet     = "192.168.0.0/24"
	)
	f, err := os.CreateTemp("", "config.json")
	require.NoError(t, err)
	configParameters := &parameters{
		ServerAddress:     serverAddress,
		GRPCServerAddress: grpcServerAddress,
		BaseURL:           baseURL,
		FileStoragePath:   fileStoragePath,
		HMACKey:           hmacKey,
		DatabaseDSN:       databaseDSN,
		EnableHTTPS:       true,
		TrustedSubnet:     trustedSubnet,
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
	assert.Equal(t, grpcServerAddress, cfg.GRPCServerAddress())
	assert.Equal(t, baseURL, cfg.BaseURL())
	assert.Equal(t, fileStoragePath, cfg.FileStoragePath())
	assert.Equal(t, hmacKey, cfg.HMACKey())
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN())
	assert.True(t, cfg.EnableHTTPS())
	assert.Equal(t, trustedSubnet, cfg.TrustedSubnet())

	_ = os.Remove(f.Name())
}

func TestBuilder_LoadFlags(t *testing.T) {
	var (
		serverAddress     = "localhost:8080"
		grpcServerAddress = "localhost:50051"
		baseURL           = "http://localhost:8080"
		fileStoragePath   = "/path"
		databaseDSN       = "dsn"
		trustedSubnet     = "192.168.0.0/24"
		builder           = &Builder{
			parameters: &parameters{},
			flags:      flagParameters,
			arguments: []string{
				"-a", serverAddress,
				"-g", grpcServerAddress,
				"-b", baseURL,
				"-f", fileStoragePath,
				"-d", databaseDSN,
				"-t", trustedSubnet,
				"-s",
			},
		}
	)

	cfg, err := builder.LoadFlags().Build()
	require.NoError(t, err)
	assert.Equal(t, serverAddress, cfg.ServerAddress())
	assert.Equal(t, grpcServerAddress, cfg.GRPCServerAddress())
	assert.Equal(t, baseURL, cfg.BaseURL())
	assert.Equal(t, fileStoragePath, cfg.FileStoragePath())
	assert.Equal(t, databaseDSN, cfg.DatabaseDSN())
	assert.True(t, cfg.EnableHTTPS())
	assert.Equal(t, trustedSubnet, cfg.TrustedSubnet())
}
