package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadEnv_Success(t *testing.T) {
	assert.Equal(t, 0, viper.GetInt("TEST_INT"))
	assert.Equal(t, "", viper.GetString("TEST_STRING"))
	assert.Equal(t, false, viper.GetBool("TEST_BOOL"))

	AddPaths(".")
	err := LoadEnv("testing")
	require.NoError(t, err)

	assert.Equal(t, 123, viper.GetInt("TEST_INT"))
	assert.Equal(t, "TEST", viper.GetString("TEST_STRING"))
	assert.Equal(t, true, viper.GetBool("TEST_BOOL"))
}

func TestLoadEnv_InvalidPath(t *testing.T) {
	err := LoadEnv("config_test")
	require.Error(t, err)
}

func TestLoadEnv_InvalidName(t *testing.T) {
	err := LoadEnv("")
	require.Error(t, err)
}
