package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/speshak/grizzl-e-monitor/internal/monitor"
	"github.com/speshak/grizzl-e-monitor/internal/prometheus"
	"github.com/speshak/grizzl-e-monitor/internal/timescale"
)

// Default values for paramters
const DefaultConnectApiHost = "https://connect-api.unitedchargers.com"
const DefaultInfluxOrg = "default"
const DefaultInfluxBucket = "default"

// LoadConfig loads configuration from environment variables
func LoadConfig() (*monitor.Config, *timescale.Config, error) {
	apiHost := os.Getenv("GRIZZLE_CONNECT_API_URL")
	if apiHost == "" {
		apiHost = DefaultConnectApiHost
	}

	username := os.Getenv("GRIZZLE_CONNECT_API_USERNAME")
	if username == "" {
		return nil, nil, fmt.Errorf("GRIZZLE_CONNECT_API_USERNAME environment variable is required")
	}

	password := os.Getenv("GRIZZLE_CONNECT_API_PASSWORD")
	if password == "" {
		return nil, nil, fmt.Errorf("GRIZZLE_CONNECT_API_PASSWORD environment variable is required")
	}

	debug := os.Getenv("GRIZZLE_CONNECT_DEBUG")
	if debug == "" {
		debug = "false"
	}

	timescaleConfig, err := LoadTimescaleConfig()
	if err != nil {
		log.Printf("Error loading TimescaleDB config: %v\n", err)
		log.Println("TimescaleDB will not be used")
		timescaleConfig = nil
	}

	return &monitor.Config{
			APIHost:  apiHost,
			Username: username,
			Password: password,
			Debug:    debug == "true",
		},
		timescaleConfig, nil
}

func LoadTimescaleConfig() (*timescale.Config, error) {
	url := os.Getenv("TIMESCALE_URL")

	if url == "" {
		return nil, fmt.Errorf("TIMESCALE_URL environment variable is required")
	}

	return &timescale.Config{
		Url: url,
	}, nil
}

func main() {
	config, timescaleConfig, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Start monitoring stations
	monitor := monitor.NewStationMonitor(config)

	prom := prometheus.NewPrometheusPublisher()

	if timescaleConfig != nil {
		timescale := timescale.NewTimescalePublisher(timescaleConfig)
		monitor.TransactionHistoryPublisher = timescale
	}

	monitor.TransactionStatsPublisher = prom
	monitor.StationStatusPublisher = prom

	ctx := context.Background()
	errs := make(chan error, 1)
	go func() {
		errs <- monitor.MonitorStations(ctx)
	}()

	// Handle any errors
	if err := <-errs; err != nil {
		log.Fatal(err)
	}

}
