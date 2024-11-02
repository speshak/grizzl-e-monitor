package monitor

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
)

type StationMonitor struct {
	Config    *Config
	Connect   connect.ConnectAPI
	Scheduler gocron.Scheduler

	TransactionHistoryPublisher TransactionHistoryPublisher
	TransactionStatsPublisher   TransactionStatsPublisher
	StationStatusPublisher      StationStatusPublisher
}

func NewStationMonitor(config *Config) *StationMonitor {
	connect := connect.NewConnectAPI(config.Username, config.Password, config.APIHost)

	if config.Debug {
		connect.SetDebug()
	}

	s, err := gocron.NewScheduler()

	if err != nil {
		log.Fatalf("Error creating scheduler %v", err)
	}

	ret := StationMonitor{
		Config:    config,
		Connect:   connect,
		Scheduler: s,
	}

	return &ret
}

func (m *StationMonitor) MonitorStations(ctx context.Context) error {
	log.Printf("Monitoring stations")

	// I've decided that station creation/removal is a rare event, so we'll just create the jobs once
	// If someone changes the station configuration, they'll need to restart the monitor

	err := m.CreateJobsForStations()
	if err != nil {
		return err
	}

	m.Scheduler.Start()

	defer func() {
		err := m.Scheduler.Shutdown()
		log.Printf("Error shutting down scheduler: %v", err)
	}()

	// Check if the context has been cancelled
	<-ctx.Done()
	log.Println("Context cancelled, stopping monitoring")
	return ctx.Err()
}

func (m *StationMonitor) CreateJobsForStations() error {
	// Get the list of stations
	stations, err := m.Connect.GetStations()
	if err != nil {
		return err
	}

	log.Print("Creating monitor jobs for stations")

	// Iterate over the stations
	for _, station := range stations {
		log.Printf("Creating monitor jobs for station %s", station.ID)

		// Stats
		_, err := m.Scheduler.NewJob(
			gocron.DurationRandomJob(time.Minute*5, time.Minute*20),
			gocron.NewTask(
				func() {
					m.stationStats(station)
					m.transactionStats(station)
				},
			),
			gocron.WithTags("station_stats"),
		)

		if err != nil {
			return err
		}

		// Transactions
		_, err = m.Scheduler.NewJob(
			gocron.DurationRandomJob(time.Minute*60, time.Minute*90),
			gocron.NewTask(
				func() {
					m.transactionHistory(station)
				},
			),
			gocron.WithTags("transaction"),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// Get the station's transaction stats
func (m *StationMonitor) transactionStats(station connect.Station) {
	// Get the transaction statistics for the station
	stats, err := m.Connect.GetTransactionStatistics(station.ID)
	if err != nil {
		log.Printf("Error getting transaction statistics for station %s: %v", station.ID, err)
		return
	}

	log.Printf("Station %s statistics: %+v", station.ID, stats)
	m.TransactionStatsPublisher.PublishTransactionStats(station.ID, stats)
}

// Get the station's stats
func (m *StationMonitor) stationStats(station connect.Station) {
	station, err := m.Connect.GetStation(station.ID)
	if err != nil {
		log.Printf("Error getting station %s: %v", station.ID, err)
		return
	}

	m.StationStatusPublisher.PublishStationStatus(station)
}

func (m *StationMonitor) transactionHistory(station connect.Station) {
	// Get all transactions for the station
	transactions, err := m.Connect.GetAllTransactions(station.ID)
	if err != nil {
		log.Printf("Error getting all transactions for station %s: %v", station.ID, err)
		return
	}

	for _, transaction := range transactions {
		// If we've already published the history, don't do it again
		// This is up to the implementation of the TransactionHistoryPublisher to check.
		if !m.TransactionHistoryPublisher.TransactionPublished(transaction) {
			log.Printf("Publishing transaction history for transaction %s", transaction.ID)
			// The all transactions endpoint gets a subset of the transaction data, so we need to get the full transaction
			fullTrans, err := m.Connect.GetTransaction(transaction.ID)

			if err != nil {
				log.Printf("Error getting full transaction %s: %v", transaction.ID, err)
				continue
			}
			err = m.TransactionHistoryPublisher.PublishTransactionHistory(station.ID, fullTrans)

			if err != nil {
				log.Printf("Error publishing transaction history for transaction %s: %v", transaction.ID, err)
			}
		} else {
			log.Printf("Transaction %s already published", transaction.ID)
		}
	}
}
