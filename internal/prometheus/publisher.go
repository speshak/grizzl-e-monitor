package prometheus

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
)

type PrometheusPublisher struct {
	// Prometheus metrics
	LastUpdate      prometheus.Gauge
	StationSessions *prometheus.GaugeVec
	TotalEnergy     *prometheus.GaugeVec
	TotalDuration   *prometheus.GaugeVec
	TopSession      *prometheus.GaugeVec
	AveEnergy       *prometheus.GaugeVec

	EnergyCost     *prometheus.GaugeVec
	AvaliablePower *prometheus.GaugeVec
	MaxPower       *prometheus.GaugeVec
}

func NewPrometheusPublisher() *PrometheusPublisher {
	stationLabels := []string{"station_id"}
	connectorLabels := []string{"station_id", "connector"}

	return &PrometheusPublisher{
		LastUpdate: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "last_poll_timestamp_seconds",
			Help:      "The last time the station was polled",
		}),
		StationSessions: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "sessions_total",
			Help:      "The total number of charging sessions",
		}, stationLabels),
		TotalEnergy: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "total_energy",
			Help:      "The total amount of energy consumed",
		}, stationLabels),
		TotalDuration: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "total_duration_seconds",
			Help:      "The total duration of charging sessions",
		}, stationLabels),
		TopSession: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "top_session_duration_seconds",
			Help:      "The 'top' session",
		}, stationLabels),
		AveEnergy: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "ave_energy",
			Help:      "The average energy consumed in a session",
		}, stationLabels),
		EnergyCost: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "energy_cost_dollars",
			Help:      "The configured cost of electrical power",
		}, stationLabels),
		AvaliablePower: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "avaliable_power_kw",
			Help:      "The amount of power avaliable to the station",
		}, connectorLabels),
		MaxPower: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "grizzl_e",
			Subsystem: "station",
			Name:      "max_power_kw",
			//TODO: Add help text when we know what this is
		}, connectorLabels),
	}
}

func (p *PrometheusPublisher) PublishStationStatus(station connect.Station) {
	p.LastUpdate.SetToCurrentTime()
	p.EnergyCost.With(prometheus.Labels{"station_id": station.ID}).Set(station.PriceKW)

	for _, connector := range station.Connectors {
		labels := prometheus.Labels{"station_id": station.ID, "connector": strconv.Itoa(connector.ID)}
		p.AvaliablePower.With(labels).Set(connector.Power)
		p.MaxPower.With(labels).Set(connector.MaxPower)
	}
}

func (p *PrometheusPublisher) PublishTransactionStats(stationId string, stats connect.TransactionStats) {
	labels := prometheus.Labels{"station_id": stationId}

	// TODO: Prometheus metrics best practices suggests that energy
	// should be expressed in Joules, and that power should be a counter
	// of energy.
	p.StationSessions.With(labels).Set(float64(stats.Sessions))
	p.TotalEnergy.With(labels).Set(stats.AverageEnergy)
	p.TotalDuration.With(labels).Set(float64(stats.Duration))
	p.AveEnergy.With(labels).Set(stats.AverageEnergy)
	p.TopSession.With(labels).Set(float64(stats.TopSession))
}

func (p *PrometheusPublisher) Close() error {
	// Nothing to close for Prometheus
	return nil
}
