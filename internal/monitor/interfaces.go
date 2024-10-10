package monitor

import (
	"context"

	"github.com/speshak/grizzl-e-monitor/pkg/connect"
)

type TransactionStatsPublisher interface {
	PublishTransactionStats(stationId string, stats connect.TransactionStats)
}

type StationStatusPublisher interface {
	PublishStationStatus(station connect.Station)
	Close() error
}

type TransactionHistoryPublisher interface {
	PublishTransactionHistory(stationId string, transaction connect.Transaction)
	TransactionPublished(transaction connect.Transaction) bool
	Close() error
}

type SingleStationMonitor interface {
	MonitorStation(ctx context.Context, station connect.Station)
}
