package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("GRIZZLE_CONNECT_API_URL", "https://test-api.com")
	os.Setenv("GRIZZLE_CONNECT_API_USERNAME", "testuser")
	os.Setenv("GRIZZLE_CONNECT_API_PASSWORD", "testpass")
	os.Setenv("GRIZZLE_CONNECT_DEBUG", "true")

	config, influxConfig, err := LoadConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Nil(t, influxConfig)
	assert.Equal(t, "https://test-api.com", config.APIHost)
	assert.Equal(t, "testuser", config.Username)
	assert.Equal(t, "testpass", config.Password)
	assert.True(t, config.Debug)
}

func TestLoadConfig_MissingUsername(t *testing.T) {
	os.Setenv("GRIZZLE_CONNECT_API_URL", "https://test-api.com")
	os.Unsetenv("GRIZZLE_CONNECT_API_USERNAME")
	os.Setenv("GRIZZLE_CONNECT_API_PASSWORD", "testpass")
	os.Setenv("GRIZZLE_CONNECT_DEBUG", "true")

	config, influxConfig, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, config)
	assert.Nil(t, influxConfig)
}

func TestLoadConfig_MissingPassword(t *testing.T) {
	os.Setenv("GRIZZLE_CONNECT_API_URL", "https://test-api.com")
	os.Setenv("GRIZZLE_CONNECT_API_USERNAME", "testuser")
	os.Unsetenv("GRIZZLE_CONNECT_API_PASSWORD")
	os.Setenv("GRIZZLE_CONNECT_DEBUG", "true")

	config, influxConfig, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, config)
	assert.Nil(t, influxConfig)
}

func TestLoadConfig_MissingAPIHost(t *testing.T) {
	os.Unsetenv("GRIZZLE_CONNECT_API_URL")
	os.Setenv("GRIZZLE_CONNECT_API_USERNAME", "testuser")
	os.Setenv("GRIZZLE_CONNECT_API_PASSWORD", "testpass")
	os.Setenv("GRIZZLE_CONNECT_DEBUG", "true")

	config, _, err := LoadConfig()
	require.NoError(t, err)
	// Check that we used the default host
	assert.Equal(t, DefaultConnectApiHost, config.APIHost)
}

func TestLoadConfig_MissingDebug(t *testing.T) {
	os.Setenv("GRIZZLE_CONNECT_API_URL", "https://test-api.com")
	os.Setenv("GRIZZLE_CONNECT_API_USERNAME", "testuser")
	os.Setenv("GRIZZLE_CONNECT_API_PASSWORD", "testpass")
	os.Unsetenv("GRIZZLE_CONNECT_DEBUG")

	config, _, err := LoadConfig()
	require.NoError(t, err)
	assert.False(t, config.Debug)
}

func TestLoadInfluxConfig(t *testing.T) {
	os.Setenv("INFLUX_HOST", "http://test-influx.com")
	os.Setenv("INFLUX_TOKEN", "testtoken")
	os.Setenv("INFLUX_ORG", "testorg")
	os.Setenv("INFLUX_BUCKET", "testbucket")

	influxConfig, err := LoadInfluxConfig()
	require.NoError(t, err)
	assert.NotNil(t, influxConfig)
	assert.Equal(t, "http://test-influx.com", influxConfig.Host)
	assert.Equal(t, "testtoken", influxConfig.Token)
	assert.Equal(t, "testorg", influxConfig.Org)
	assert.Equal(t, "testbucket", influxConfig.Bucket)
}

func TestLoadInfluxConfig_MissingToken(t *testing.T) {
	os.Setenv("INFLUX_HOST", "http://test-influx.com")
	os.Unsetenv("INFLUX_TOKEN")
	os.Setenv("INFLUX_ORG", "testorg")
	os.Setenv("INFLUX_BUCKET", "testbucket")

	influxConfig, err := LoadInfluxConfig()
	require.Error(t, err)
	assert.Nil(t, influxConfig)
}

func TestLoadInfluxConfig_MissingOrg(t *testing.T) {
	os.Setenv("INFLUX_HOST", "http://test-influx.com")
	os.Setenv("INFLUX_TOKEN", "testtoken")
	os.Setenv("INFLUX_BUCKET", "testbucket")
	os.Unsetenv("INFLUX_ORG")

	influxConfig, err := LoadInfluxConfig()
	require.NoError(t, err)
	assert.Equal(t, DefaultInfluxOrg, influxConfig.Org)
}

func TestLoadInfluxConfig_MissingBucket(t *testing.T) {
	os.Setenv("INFLUX_HOST", "http://test-influx.com")
	os.Setenv("INFLUX_TOKEN", "testtoken")
	os.Setenv("INFLUX_ORG", "testorg")
	os.Unsetenv("INFLUX_BUCKET")

	influxConfig, err := LoadInfluxConfig()
	require.NoError(t, err)
	assert.Equal(t, DefaultInfluxBucket, influxConfig.Bucket)
}
