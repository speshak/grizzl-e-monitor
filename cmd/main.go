package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/speshak/grizzl-e-prom/internal/influx"
	"github.com/speshak/grizzl-e-prom/internal/monitor"
	"github.com/speshak/grizzl-e-prom/internal/prometheus"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*monitor.Config, error) {
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

	return &monitor.Config{
		APIHost:     apiHost,
		Username:    username,
		Password:    password,
		InfluxHost:  os.Getenv("INFLUX_HOST"),
		InfluxToken: os.Getenv("INFLUX_TOKEN"),
		Debug:       debug == "true",
	}, nil
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Start monitoring stations
	monitor := monitor.NewStationMonitor(config)

	prom := prometheus.NewPrometheusPublisher()
	influx := influx.NewInfluxPublisher(config.InfluxHost, config.InfluxToken)

	monitor.TransactionHistoryPublisher = influx
	monitor.TransactionStatsPublisher = prom
	monitor.StationStatusPublisher = prom
	go monitor.MonitorStations()

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
