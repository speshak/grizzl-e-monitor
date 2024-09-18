package main

import (
	"fmt"
	"os"
)

// Config holds the configuration values
type Config struct {
	APIHost  string
	Username string
	Password string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	apiHost := os.Getenv("GRIZZLE_CONNECT_API_URL")
	if apiHost == "" {
		apiHost = "https://connect-api.unitedchargers.com"
	}

	username := os.Getenv("GRIZZLE_CONNECT_API_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("GRIZZLE_CONNECT_API_USERNAME environment variable is required")
	}

	password := os.Getenv("GRIZZLE_CONNECT_API_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("GRIZZLE_CONNECT_API_PASSWORD environment variable is required")
	}

	return &Config{
		APIHost:  apiHost,
		Username: username,
		Password: password,
	}, nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config loaded: %+v\n", config)
}
