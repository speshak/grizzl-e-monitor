package monitor

import (
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
	PublishTransactionHistory(stationId string, transaction connect.Transaction) error
	TransactionPublished(transaction connect.Transaction) bool
	Close() error
}
