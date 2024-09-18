package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds the configuration values
type Config struct {
	APIHost  string
	Username string
	Password string
	Debug    bool
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

	debug := os.Getenv("GRIZZLE_CONNECT_DEBUG")
	if debug == "" {
		debug = "false"
	}

	return &Config{
		APIHost:  apiHost,
		Username: username,
		Password: password,
		Debug:    debug == "true",
	}, nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Start monitoring stations
	go MonitorStations(config)

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
