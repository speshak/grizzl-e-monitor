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
- `INFLUX_DB_HOST` - Hostname of Influx
- `INFLUX_DB_KEY` - API key for Influx

## Running

The easiest way to run the scraper is to use the docker image. Make sure to set the environment variables as needed.

```bash
docker run -d -p 8080:8080 -e GRIZZLE_CONNECT_API_USERNAME=your-username -e GRIZZLE_CONNECT_API_PASSWORD=your-password ghcr.io/speshak/grizzl-e-prom:main
```

## API Client

There is an implementation of a grizzl-e connect API client in `pkg/connect`. The
Connect service does not publish any API information so this client was built by
capturing traffic from the iPadOS client app. It certainly is feature incomplete
and could break in the future.
