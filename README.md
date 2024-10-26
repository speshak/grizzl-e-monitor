# Grizzl-E Connect Monitor

This is a tool to monitor the Grizzl-E Connect API. It is designed to be run as
a service and will scrape the Grizzl-E Connect API for the current status of
your Grizzl-E charger and expose Prometheus metrics and push to InfluxDB.

## Configuration

The following environment variables can be used to configure the scraper:

- `GRIZZLE_CONNECT_API_URL`: The URL of the Grizzl-E Connect API. Defaults to
  `connect-api.unitedchargers.com`.
- `GRIZZLE_CONNECT_API_USERNAME`: The username to use when authenticating with
  the Grizzl-E Connect API.
- `GRIZZLE_CONNECT_API_PASSWORD`: The password to use when authenticating with
  the Grizzl-E Connect API.

InfluxDB output can be enabled to defining the following
- `INFLUX_HOST` - Hostname of Influx
- `INFLUX_TOKEN` - API key for Influx
- `INFLUX_BUCKET` - InfluxDB bucket name. Defaults to `default`
- `INFLUX_ORG` - InfluxDB Organization name

TimescaleDB output for transaction metrics can be enabled by defining:
- `TIMESCALE_URL` - A DB URL for the PostgreSQL database.

The database should be empty, go-migrate will be used to create the required
tables during startup.

## Running

The easiest way to run the scraper is to use the docker image. Make sure to set the environment variables as needed.

```bash
docker run -d -p 8080:8080 -e GRIZZLE_CONNECT_API_USERNAME=your-username -e GRIZZLE_CONNECT_API_PASSWORD=your-password ghcr.io/speshak/grizzl-e-monitor:main
```

## API Client

There is an implementation of a grizzl-e connect API client in `pkg/connect`. The
Connect service does not publish any API information so this client was built by
capturing traffic from the iPadOS client app. It certainly is feature incomplete
and could break in the future.
