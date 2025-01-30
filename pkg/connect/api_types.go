package connect

import "time"

// A message from the Connect API. This is used in error responses.
type ApiMessage struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

// An error response from the Connect API
type ApiError struct {
	StatusCode int        `json:"statusCode"`
	Timestamp  string     `json:"timestamp"`
	Path       string     `json:"path"`
	Message    ApiMessage `json:"message"`
}

// A user object from the Connect API
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Language  string `json:"language"`
	ID        string `json:"id"`
}

// A login response from the Connect API
type LoginResponse struct {
	Token          string `json:"token"`
	IsInitialLogIn bool   `json:"isInitialLogIn"`
	User           User   `json:"user"`
}

type MeterValues struct {
	Date                       []time.Time `json:"date"`
	CurrentImport              []float64   `json:"currentImport"`
	CurrentOffered             []float64   `json:"currentOffered"`
	EnergyActiveImportRegister []int       `json:"energyActiveImportRegister"`
	PowerActiveImport          []float64   `json:"powerActiveImport"`
	SoC                        []int       `json:"SoC"`
	Temperature                []float64   `json:"temperature"`
	Voltage                    []int       `json:"voltage"`
}

type Transaction struct {
	ID               string      `json:"_id"`
	User             string      `json:"user"`
	Station          string      `json:"station"`
	IdTag            string      `json:"idTag"`
	ConnectorId      int         `json:"connectorId"`
	StartAt          string      `json:"startAt"`
	Duration         float64     `json:"duration"`
	Energy           int         `json:"energy"`
	Status           int         `json:"status"`
	Power            int         `json:"power"`
	Currency         string      `json:"currency"`
	PriceKW          float64     `json:"priceKW"`
	PriceTotal       float64     `json:"priceTotal"`
	MeterStart       int         `json:"meterStart"`
	MeterStop        int         `json:"meterStop"`
	StopAt           string      `json:"stopAt"`
	StopReason       string      `json:"stopReason"`
	AverageCurrent   float64     `json:"averageCurrent"`
	ChargingDuration float64     `json:"chargingDuration"`
	MeterValues      MeterValues `json:"meterValues"`
}

type Connector struct {
	ID        int     `json:"id"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Power     float64 `json:"power"`
	MaxPower  float64 `json:"maxPower"`
	ErrorCode string  `json:"errorCode"`
}

// Response type of the stations endipoint
type GetStationsResponse struct {
	Stations []Station `json:"stations"`
}

type Station struct {
	ID           string      `json:"id"`
	Identity     string      `json:"identity"`
	SerialNumber string      `json:"serialNumber"`
	Online       bool        `json:"online"`
	Mode         string      `json:"mode"`
	Status       string      `json:"status"`
	ErrorCode    string      `json:"errorCode"`
	Connectors   []Connector `json:"connectors"`
	Currency     string      `json:"currency"`
	PriceKW      float64     `json:"priceKW"`
	// There are a load more fields in the response, but we only need these for now
}

type TransactionStats struct {
	TotalEnergy   int     `json:"totalEnergy"`
	Sessions      int     `json:"sessions"`
	AverageEnergy float64 `json:"averageEnergy"`
	Duration      int     `json:"duration"`
	TopSession    int     `json:"topSession"`
	Currency      string  `json:"currency"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"data"`
}

type GetTransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}

/**
 * Contents of the X-Application-Version header
 */
type ApplicationVersion struct {
	Id                    string `json:"_id"`
	ApplicationName       string `json:"applicationName"`
	IosLatestVersion      string `json:"iosLatestVersion"`
	IosMinimalVersion     string `json:"iosMinimalVersion"`
	AndroidLatestVersion  string `json:"androidLatestVersion"`
	AndroidMinimalVersion string `json:"androidMinimalVersion"`
}
