package influx

import (
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/speshak/grizzl-e-prom/pkg/connect"
)

type InfluxPublisher struct {
	// InfluxDB client
	InfluxClient influxdb2.Client
}

// NewInfluxPublisher creates a new InfluxPublisher
func NewInfluxPublisher(host, token string) *InfluxPublisher {
	return &InfluxPublisher{
		InfluxClient: influxdb2.NewClient(host, token),
	}
}

func (p *InfluxPublisher) TransactionPublished(transactionId string) bool {
	return false
}

func (p *InfluxPublisher) PublishTransactionHistory(stationId string, transaction connect.Transaction) {
	log.Printf("Logging Transaction '%s' starting", transaction.ID)
	log.Printf("%d data points to log", len(transaction.MeterValues.Date))

	writeApi := p.InfluxClient.WriteAPI("grizzl_e", "station")

	// Loop over the metrics and add measurements for each one
	for index, metricDate := range transaction.MeterValues.Date {
		p := influxdb2.NewPointWithMeasurement("meter").
			AddTag("transaction", transaction.ID).
			AddTag("station", transaction.Station).
			SetTime(metricDate).
			AddField("current", transaction.MeterValues.CurrentImport[index]).
			AddField("current_offered", transaction.MeterValues.CurrentOffered[index]).
			AddField("voltage", transaction.MeterValues.Voltage[index]).
			AddField("power_import", transaction.MeterValues.PowerActiveImport[index])

		writeApi.WritePoint(p)
	}
	log.Printf("Logging Transaction '%s' finished", transaction.ID)
}
