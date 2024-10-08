package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/speshak/grizzl-e-monitor/internal/influx"
	"github.com/speshak/grizzl-e-monitor/internal/monitor"
	"github.com/speshak/grizzl-e-monitor/internal/prometheus"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*monitor.Config, *influx.InfluxConfig, error) {
	apiHost := os.Getenv("GRIZZLE_CONNECT_API_URL")
	if apiHost == "" {
		apiHost = "https://connect-api.unitedchargers.com"
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

	influxConfig, err := LoadInfluxConfig()
	if err != nil {
		log.Printf("Error loading InfluxDB config: %v\n", err)
		log.Println("InfluxDB will not be used")
		influxConfig = nil
	}

	return &monitor.Config{
			APIHost:  apiHost,
			Username: username,
			Password: password,
			Debug:    debug == "true",
		},
		influxConfig, nil
}

func LoadInfluxConfig() (*influx.InfluxConfig, error) {
	influxHost := os.Getenv("INFLUX_HOST")
	if influxHost == "" {
		influxHost = "http://localhost:8086"
	}

	influxToken := os.Getenv("INFLUX_TOKEN")
	if influxToken == "" {
		return nil, fmt.Errorf("INFLUX_TOKEN environment variable is required")
	}

	org := os.Getenv("INFLUX_ORG")
	if org == "" {
		org = "default"
	}

	bucket := os.Getenv("INFLUX_BUCKET")
	if bucket == "" {
		bucket = "default"
	}

	return &influx.InfluxConfig{
		Host:   influxHost,
		Token:  influxToken,
		Org:    org,
		Bucket: bucket,
	}, nil
}

func main() {
	config, influxConfig, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Start monitoring stations
	monitor := monitor.NewStationMonitor(config)

	prom := prometheus.NewPrometheusPublisher()

	if influxConfig != nil {
		influx := influx.NewInfluxPublisher(influxConfig)
		monitor.TransactionHistoryPublisher = influx
	}

	monitor.TransactionStatsPublisher = prom
	monitor.StationStatusPublisher = prom

	ctx := context.Background()
	go monitor.MonitorStations(ctx)

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
