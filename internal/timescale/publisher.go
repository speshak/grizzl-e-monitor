package timescale

import (
	"database/sql"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	pgmig "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/speshak/grizzl-e-monitor/pkg/connect"
)

type Config struct {
	Url string
}

type TimescalePublisher struct {
	DbClient *sql.DB
	Config   *Config
}

//go:embed migrations/*.sql
var migrationsBox embed.FS

func (t *TimescalePublisher) RunMigrations() error {
	log.Printf("Running database migrations")
	migrations, err := iofs.New(migrationsBox, "migrations")
	if err != nil {
		return err
	}

	driver, err := pgmig.WithInstance(t.DbClient, &pgmig.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		migrations,
		t.Config.Url,
		driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, _, _ := m.Version()
	log.Printf("Database version: %v", version)
	return nil
}

func NewTimescalePublisher(config *Config) *TimescalePublisher {
	db, err := sql.Open("postgres", config.Url)

	if err != nil {
		log.Fatalf("Error opening database connection: %v\n", err)
	}

	// Create connection and run migrations
	publisher := &TimescalePublisher{
		DbClient: db,
		Config:   config,
	}

	err = publisher.RunMigrations()

	if err != nil {
		log.Fatalf("Error running migrations: %v\n", err)
	}

	return publisher
}

func (t *TimescalePublisher) Close() error {
	return t.DbClient.Close()
}

func (t *TimescalePublisher) PublishTransactionHistory(stationId string, transaction connect.Transaction) error {
	log.Printf("Logging Transaction '%s' starting", transaction.ID)
	log.Printf("%d data points to log", len(transaction.MeterValues.Date))

	// Insert the transaction
	_, err := t.DbClient.Exec(`
		INSERT INTO transactions (
			id, duration, station, startAt, stopAt, status, power, currency, priceKW,
			priceTotal, meterStart, meterStop, stopReason, averageCurrent, chargingDuration
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (id) DO UPDATE SET
			duration = EXCLUDED.duration,
			station = EXCLUDED.station,
			startAt = EXCLUDED.startAt,
			stopAt = EXCLUDED.stopAt,
			status = EXCLUDED.status,
			power = EXCLUDED.power,
			currency = EXCLUDED.currency,
			priceKW = EXCLUDED.priceKW,
			priceTotal = EXCLUDED.priceTotal,
			meterStart = EXCLUDED.meterStart,
			meterStop = EXCLUDED.meterStop,
			stopReason = EXCLUDED.stopReason,
			averageCurrent = EXCLUDED.averageCurrent,
			chargingDuration = EXCLUDED.chargingDuration
		`,
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
	)

	log.Printf("Transaction '%s' inserted", transaction.ID)

	if err != nil {
		log.Fatalf("Error inserting transaction: %v", err)
		return err
	}

	// Prepare the meter statement (we do this a bunch, so it makes sense to prepare it)
	meter_stmt, err := t.DbClient.Prepare(`
		INSERT INTO meter_values (
			date, transaction_id, currentImport, currentOffered,
			energyActiveImportRegister, powerActiveImport, soC, temperature, voltage
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (date, transaction_id) DO UPDATE SET
			currentImport = EXCLUDED.currentImport,
			currentOffered = EXCLUDED.currentOffered,
			energyActiveImportRegister = EXCLUDED.energyActiveImportRegister,
			powerActiveImport = EXCLUDED.powerActiveImport,
			soC = EXCLUDED.soC,
			temperature = EXCLUDED.temperature,
			voltage = EXCLUDED.voltage
	`)
	if err != nil {
		log.Fatalf("Error preparing meter statement: %v", err)
		return err
	}

	for index, metricDate := range transaction.MeterValues.Date {
		_, err := meter_stmt.Exec(
			metricDate,
			transaction.ID,
			transaction.MeterValues.CurrentImport[index],
			transaction.MeterValues.CurrentOffered[index],
			transaction.MeterValues.Voltage[index],
			transaction.MeterValues.PowerActiveImport[index],
			transaction.MeterValues.SoC[index],
			transaction.MeterValues.Temperature[index],
			transaction.MeterValues.Voltage[index],
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TimescalePublisher) TransactionPublished(transaction connect.Transaction) bool {
	var count int
	// An in-progress transaction will not have a stop time, so we will consider
	// the transaction unpublished until the stop time is set
	err := t.DbClient.QueryRow(`
		SELECT COUNT(1) FROM transactions
		WHERE id = $1
		AND stopat IS NOT NULL`,
		transaction.ID).Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	return count > 0
}
