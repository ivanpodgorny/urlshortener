package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	names := []string{"printf", "structtag"}
	confFileNames := []string{"shadow", "sortslice"}
	c, err := NewBuilder().SetDefaultAnalyzersNames(names...).LoadFile().Build()
	assert.NoError(t, err)
	assert.Equal(t, names, c.AnalyzersNames())

	f := prepareConfFile(t, &parameters{AnalyzersNames: confFileNames})
	defer func(name string) {
		require.NoError(t, os.Remove(name))
	}(f.Name())

	c, err = NewBuilder().SetDefaultAnalyzersNames(names...).LoadFile().Build()
	assert.NoError(t, err)
	assert.Equal(t, confFileNames, c.AnalyzersNames())
}

func prepareConfFile(t *testing.T, p *parameters) *os.File {
	pwd, err := getPWD()
	require.NoError(t, err)
	f, err := os.Create(filepath.Join(pwd, configFileName))
	require.NoError(t, err)
	defer func(f *os.File) {
		require.NoError(t, f.Close())
	}(f)

	data, err := json.Marshal(p)
	require.NoError(t, err)

	_, err = f.Write(data)
	require.NoError(t, err)

	return f
}
