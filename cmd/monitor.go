package main

import {
	"log"
}

func monitorStations(&Config config) {
	log.Printf("Monitoring stations")

	// Create a new ConnectAPI client
	c := connect.NewConnectAPI(config.Username, config.Password, config.APIHost)

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
		}

		// Sleep for 5 minutes
		time.Sleep(5 * time.Minute)
	}
}
