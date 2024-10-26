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
	"github.com/speshak/grizzl-e-monitor/internal/timescale"
)

// Default values for paramters
const DefaultConnectApiHost = "https://connect-api.unitedchargers.com"
const DefaultInfluxOrg = "default"
const DefaultInfluxBucket = "default"

// LoadConfig loads configuration from environment variables
func LoadConfig() (*monitor.Config, *influx.InfluxConfig, *timescale.Config, error) {
	apiHost := os.Getenv("GRIZZLE_CONNECT_API_URL")
	if apiHost == "" {
		apiHost = DefaultConnectApiHost
	}

	username := os.Getenv("GRIZZLE_CONNECT_API_USERNAME")
	if username == "" {
		return nil, nil, nil, fmt.Errorf("GRIZZLE_CONNECT_API_USERNAME environment variable is required")
	}

	password := os.Getenv("GRIZZLE_CONNECT_API_PASSWORD")
	if password == "" {
		return nil, nil, nil, fmt.Errorf("GRIZZLE_CONNECT_API_PASSWORD environment variable is required")
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
		influxConfig, timescaleConfig, nil
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
		org = DefaultInfluxOrg
	}

	bucket := os.Getenv("INFLUX_BUCKET")
	if bucket == "" {
		bucket = DefaultInfluxBucket
	}

	return &influx.InfluxConfig{
		Host:   influxHost,
		Token:  influxToken,
		Org:    org,
		Bucket: bucket,
	}, nil
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
	config, influxConfig, timescaleConfig, err := LoadConfig()
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

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
