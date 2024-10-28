package timescale

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMigrations(t *testing.T) {
	t.Skip()
	// Is this worth testing?
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	config := &Config{Url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"}
	publisher := &TimescalePublisher{
		DbClient: db,
		Config:   config,
	}

	mock.ExpectQuery("SELECT CURRENT_DATABASE()")

	err = publisher.RunMigrations()
	require.NoError(t, err)
}

func TestPublishTransactionHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	config := &Config{Url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"}
	publisher := &TimescalePublisher{
		DbClient: db,
		Config:   config,
	}

	transaction := connect.Transaction{
		ID:               "tx1",
		Duration:         3600,
		Station:          "station1",
		StartAt:          "2021-01-01T00:00:00Z",
		StopAt:           "2021-01-01T01:00:00Z",
		Status:           1,
		Power:            50.0,
		Currency:         "USD",
		PriceKW:          0.15,
		PriceTotal:       7.5,
		MeterStart:       1000,
		MeterStop:        1500,
		StopReason:       "user",
		AverageCurrent:   10.0,
		ChargingDuration: 3600,
		MeterValues: connect.MeterValues{
			Date:                       []time.Time{time.Now()},
			CurrentImport:              []float64{10.0},
			CurrentOffered:             []float64{10.0},
			EnergyActiveImportRegister: []int{1000},
			PowerActiveImport:          []float64{50.0},
			SoC:                        []int{80},
			Temperature:                []float64{25.0},
			Voltage:                    []int{230},
		},
	}

	mock.ExpectExec("INSERT INTO transactions").WithArgs(
		transaction.ID,
		transaction.Duration,
		transaction.Station,
		transaction.StartAt,
		transaction.StopAt,
		transaction.Status,
		transaction.Power,
		transaction.Currency,
		transaction.PriceKW,
		transaction.PriceTotal,
		transaction.MeterStart,
		transaction.MeterStop,
		transaction.StopReason,
		transaction.AverageCurrent,
		transaction.ChargingDuration,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectPrepare("INSERT INTO meter_values")
	mock.ExpectExec("INSERT INTO meter_values").WithArgs(
		sqlmock.AnyArg(),
		transaction.ID,
		transaction.MeterValues.CurrentImport[0],
		transaction.MeterValues.CurrentOffered[0],
		transaction.MeterValues.Voltage[0],
		transaction.MeterValues.PowerActiveImport[0],
		transaction.MeterValues.SoC[0],
		transaction.MeterValues.Temperature[0],
		transaction.MeterValues.Voltage[0],
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = publisher.PublishTransactionHistory("station1", transaction)
	require.NoError(t, err)
}

func TestTransactionPublished(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	config := &Config{Url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"}
	publisher := &TimescalePublisher{
		DbClient: db,
		Config:   config,
	}

	transaction := connect.Transaction{
		ID: "tx1",
	}

	mock.ExpectQuery("SELECT COUNT").WithArgs(transaction.ID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	published := publisher.TransactionPublished(transaction)
	assert.True(t, published)
}

func TestClose(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	config := &Config{Url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"}
	publisher := &TimescalePublisher{
		DbClient: db,
		Config:   config,
	}

	mock.ExpectClose()

	err = publisher.Close()
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
