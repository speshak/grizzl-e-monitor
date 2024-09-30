package monitor

import "github.com/speshak/grizzl-e-prom/pkg/connect"

type TransactionStatsPublisher interface {
	PublishTransactionStats(stationId string, stats connect.TransactionStats)
}

type StationStatusPublisher interface {
	PublishStationStatus(station connect.Station)
}

type TransactionHistoryPublisher interface {
	PublishTransactionHistory(stationId string, transaction connect.Transaction)
	TransactionPublished(transactionId string) bool
}
