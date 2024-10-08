package influx

import (
	"context"
	"fmt"
	"log"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
)

type InfluxConfig struct {
	Host   string
	Token  string
	Org    string
	Bucket string
}

type InfluxPublisher struct {
	// InfluxDB client
	InfluxClient influxdb2.Client
	Config       *InfluxConfig
}

// NewInfluxPublisher creates a new InfluxPublisher
func NewInfluxPublisher(config *InfluxConfig) *InfluxPublisher {
	return &InfluxPublisher{
		InfluxClient: influxdb2.NewClient(config.Host, config.Token),
		Config:       config,
	}
}

func (p *InfluxPublisher) TransactionPublished(transaction connect.Transaction) bool {
	log.Printf("Checking if transaction '%s' has been published", transaction.ID)
	queryAPI := p.InfluxClient.QueryAPI(p.Config.Org)

	// Build a Flux query for the transaction
	// We'll consider a transaction published if there are any data points for
	// it. If that becomes a problem, we can add some sort of flag to force
	// publishing

	// I'm targeting the local InfuxDB instance, so parameterized queries aren't
	// supported.
	var buf strings.Builder
	fmt.Fprintf(&buf, `from(bucket: "%s")`, p.Config.Bucket)
	fmt.Fprintf(&buf, ` |> range(start: %s, stop: %s)`, transaction.StartAt, transaction.StopAt)
	fmt.Fprintf(&buf, ` |> filter(fn: (r) => r["transaction"] == "%s")`, transaction.ID)

	result, err := queryAPI.Query(context.Background(), buf.String())

	if err == nil {
		for result.Next() {
			// If we get any results, the transaction has been published
			return true
		}
	} else {
		log.Printf("Error querying InfluxDB: %v", err)
	}

	return false
}

func (p *InfluxPublisher) PublishTransactionHistory(stationId string, transaction connect.Transaction) {
	log.Printf("Logging Transaction '%s' starting", transaction.ID)
	log.Printf("%d data points to log", len(transaction.MeterValues.Date))

	writeApi := p.InfluxClient.WriteAPI(p.Config.Org, p.Config.Bucket)

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

func (p *InfluxPublisher) Close() error {
	p.InfluxClient.Close()
	return nil
}
