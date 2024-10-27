package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
)

type MockPrometheusPublisher struct {
	TransactionCounter *prometheus.CounterVec
}

func NewMockPrometheusPublisher() *MockPrometheusPublisher {
	return &MockPrometheusPublisher{
		TransactionCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "transaction_total",
				Help: "Total number of transactions",
			},
			[]string{"station_id"},
		),
	}
}

func (m *MockPrometheusPublisher) PublishTransactionHistory(stationId string, transaction connect.Transaction) error {
	m.TransactionCounter.WithLabelValues(stationId).Inc()
	return nil
}

func TestNewPrometheusPublisher(t *testing.T) {
	publisher := NewMockPrometheusPublisher()

	if publisher.TransactionCounter == nil {
		t.Fatalf("Expected TransactionCounter to be initialized")
	}
}

func TestPublishTransactionHistory(t *testing.T) {
	publisher := NewMockPrometheusPublisher()

	transaction := connect.Transaction{
		ID: "test-transaction",
	}

	err := publisher.PublishTransactionHistory("station-1", transaction)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := 1.0
	actual := testutil.ToFloat64(publisher.TransactionCounter.WithLabelValues("station-1"))
	if actual != expected {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
