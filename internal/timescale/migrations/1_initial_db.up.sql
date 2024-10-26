CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE transactions (
    id VARCHAR(25) PRIMARY KEY,
	duration INTERVAL,
	station VARCHAR(26),
	startAt TIMESTAMPTZ, 
	stopAt TIMESTAMPTZ,
	status INTEGER, 
	power INTEGER,
	currency VARCHAR(3),
	priceKW DOUBLE PRECISION,
	priceTotal DOUBLE PRECISION,
	meterStart INTEGER,
	meterStop INTEGER,
	stopReason VARCHAR(255), 
	averageCurrent DOUBLE PRECISION,
	chargingDuration INTEGER
);

CREATE TABLE meter_values (
    date TIMESTAMPTZ NOT NULL,
    transaction_id VARCHAR(25) REFERENCES transactions, 
	currentImport DOUBLE PRECISION,
	currentOffered DOUBLE PRECISION,
	energyActiveImportRegister INTEGER,
	powerActiveImport DOUBLE PRECISION,
	soC INTEGER,
	temperature DOUBLE PRECISION,
	voltage INTEGER,
	PRIMARY KEY (date, transaction_id)
);

SELECT create_hypertable('meter_values', 'date');