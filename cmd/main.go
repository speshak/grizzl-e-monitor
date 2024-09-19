package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Start monitoring stations
	go monitorStations()

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
