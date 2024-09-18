package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/speshak/grizzl-e-prom/pkg/connect"
)

// TotalEnergy:142754 Sessions:5 AverageEnergy:28550.8 Duration:165483161 TopSession:56857 Currency:USD}
var (
	lastUpdate = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "station_last_update",
		Help: "The last time the station was polled",
	})
	stationSessions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "station_sessions_total",
		Help: "The total number of charging sessions",
	})
	totalEnergy = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "station_total_energy",
		Help: "The total amount of energy consumed",
	})
	totalDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "station_total_duration",
		Help: "The total duration of charging sessions",
	})
	topSession = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "station_top_session",
		Help: "The 'top' session",
	})
	aveEnergy = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "station_ave_energy",
		Help: "The average energy consumed in a session",
	})
)

func MonitorStations(config *Config) error {
	log.Printf("Monitoring stations")

	// Create a new ConnectAPI client
	c := connect.NewConnectAPI(config.Username, config.Password, config.APIHost)

	// Enable debug mode if needed
	if config.Debug {
		c.SetDebug()
	}

	for {
		// Get the list of stations
		stations, err := c.GetStations()
		if err != nil {
			log.Printf("Error getting stations: %v", err)
			continue
		}

		// Iterate over the stations
		for _, station := range stations {
			// Get the transaction statistics for the station
			stats, err := c.GetTransactionStatistics(station.ID)
			if err != nil {
				log.Printf("Error getting transaction statistics for station %s: %v", station.ID, err)
				continue
			}

			log.Printf("Station %s statistics: %+v", station.ID, stats)

			lastUpdate.SetToCurrentTime()
			stationSessions.Set(float64(stats.Sessions))
			totalEnergy.Set(stats.AverageEnergy)
			totalDuration.Set(float64(stats.Duration))
			aveEnergy.Set(stats.AverageEnergy)
			topSession.Set(float64(stats.TopSession))
		}

		// Sleep for 5 minutes
		time.Sleep(5 * time.Minute)
	}
}

func MonitorStation() {

}
