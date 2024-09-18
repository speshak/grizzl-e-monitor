# Grizzl-E Connect Prometheus Scraper

This is a Prometheus scraper for the Grizzl-E Connect API. It is designed to be
run as a service and will scrape the Grizzl-E Connect API for the current
status of your Grizzl-E charger and expose it as Prometheus metrics.

## Configuration

The following environment variables can be used to configure the scraper:

- `GRIZZLE_CONNECT_API_URL`: The URL of the Grizzl-E Connect API. Defaults to
  `connect-api.unitedchargers.com`.
- `GRIZZLE_CONNECT_API_USERNAME`: The username to use when authenticating with
  the Grizzl-E Connect API.
- `GRIZZLE_CONNECT_API_PASSWORD`: The password to use when authenticating with
  the Grizzl-E Connect API.

