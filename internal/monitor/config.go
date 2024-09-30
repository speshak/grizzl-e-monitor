package monitor

// Config holds the configuration values
type Config struct {
	APIHost     string
	Username    string
	Password    string
	Debug       bool
	InfluxHost  string
	InfluxToken string
}
