package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("GRIZZLE_CONNECT_API_URL", "https://test-api.com")
	os.Setenv("GRIZZLE_CONNECT_API_USERNAME", "testuser")
	os.Setenv("GRIZZLE_CONNECT_API_PASSWORD", "testpass")
	os.Setenv("GRIZZLE_CONNECT_DEBUG", "true")

	config, influxConfig, err := LoadConfig()
	assert.NoError(t, err)
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
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Nil(t, influxConfig)
}

func TestLoadConfig_MissingPassword(t *testing.T) {
	os.Setenv("GRIZZLE_CONNECT_API_URL", "https://test-api.com")
	os.Setenv("GRIZZLE_CONNECT_API_USERNAME", "testuser")
	os.Unsetenv("GRIZZLE_CONNECT_API_PASSWORD")
	os.Setenv("GRIZZLE_CONNECT_DEBUG", "true")

	config, influxConfig, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Nil(t, influxConfig)
}

func TestLoadInfluxConfig(t *testing.T) {
	os.Setenv("INFLUX_HOST", "http://test-influx.com")
	os.Setenv("INFLUX_TOKEN", "testtoken")
	os.Setenv("INFLUX_ORG", "testorg")
	os.Setenv("INFLUX_BUCKET", "testbucket")

	influxConfig, err := LoadInfluxConfig()
	assert.NoError(t, err)
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
	assert.Error(t, err)
	assert.Nil(t, influxConfig)
}
