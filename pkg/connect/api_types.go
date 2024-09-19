package connect

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Language  string `json:"language"`
	ID        string `json:"id"`
}

type LoginResponse struct {
	Token          string `json:"token"`
	IsInitialLogIn bool   `json:"isInitialLogIn"`
	User           User   `json:"user"`
}

type MeterValues struct {
	Date                       []string  `json:"date"`
	CurrentImport              []float64 `json:"currentImport"`
	CurrentOffered             []float64 `json:"currentOffered"`
	EnergyActiveImportRegister []int     `json:"energyActiveImportRegister"`
	PowerActiveImport          []float64 `json:"powerActiveImport"`
	SoC                        []int     `json:"SoC"`
	Temperature                []float64 `json:"temperature"`
	Voltage                    []int     `json:"voltage"`
}

type Transaction struct {
	ID               string  `json:"_id"`
	User             string  `json:"user"`
	Station          string  `json:"station"`
	IdTag            string  `json:"idTag"`
	ConnectorId      int     `json:"connectorId"`
	StartAt          string  `json:"startAt"`
	Duration         int     `json:"duration"`
	Energy           int     `json:"energy"`
	Status           int     `json:"status"`
	Power            int     `json:"power"`
	Currency         string  `json:"currency"`
	PriceKW          float64 `json:"priceKW"`
	PriceTotal       float64 `json:"priceTotal"`
	MeterStart       int     `json:"meterStart"`
	MeterStop        int     `json:"meterStop"`
	StopAt           string  `json:"stopAt"`
	StopReason       string  `json:"stopReason"`
	AverageCurrent   float64 `json:"averageCurrent"`
	ChargingDuration int     `json:"chargingDuration"`
}

type Connector struct {
	ID        int     `json:"id"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Power     float64 `json:"power"`
	MaxPower  float64 `json:"maxPower"`
	ErrorCode string  `json:"errorCode"`
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
