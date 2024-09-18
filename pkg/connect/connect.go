package connect

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

// A Connect API client
type ConnectAPI struct {
	Username string
	Password string
	Token    string
	Client   *resty.Client
}

func NewConnectAPI(username, password, host string) *ConnectAPI {
	client := resty.New()
	client.
		EnableTrace().
		SetBaseURL(host).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "GrizzlEConnect/87 CFNetwork/1568.100.1 Darwin/24.0.0").
		SetHeader("x-app-client", "Apple, iPad14,3, iPadOS 18.0").
		SetHeader("x-app-version", "v0.8.0 (87)").
		SetHeader("x-application-name", "Grizzl-E Connect")

	return &ConnectAPI{
		Username: username,
		Password: password,
		Client:   client,
	}
}

// Enable debug mode in the underlying resty client
func (c *ConnectAPI) SetDebug() {
	c.Client = c.Client.SetDebug(true)
}

// Ensure the login token is valid
func (c *ConnectAPI) AssertValidToken() error {
	jwtToken, err := jwt.Parse(c.Token, func(token *jwt.Token) (interface{}, error) {
		// Since we are not verifying the token, return nil
		return nil, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %s", err)
	}

	fmt.Println(jwtToken)
	// TODO: Check if the token is expired and re-login if needed

	if jwtToken == nil || !jwtToken.Valid {
		log.Println("No token set, logging in")
		err := c.Login()
		if err != nil {
			log.Fatalf("Error logging in: %s", err)
		}
	}

	return nil
}

// Get a resty client with the auth token set
// This will log in if no token is set or the token is expired
func (c *ConnectAPI) client() *resty.Client {
	c.AssertValidToken()

	return c.Client.
		SetAuthToken(c.Token)
}

func (c *ConnectAPI) Login() error {
	// Get login token for future requests
	log.Printf("Logging in as %s", c.Username)
	result := LoginResponse{}
	errorResult := ApiError{}

	resp, err := c.Client.R().
		SetBody(map[string]interface{}{
			"emailOrPhone": c.Username,
			"password":     c.Password,
		}).
		SetResult(&result).
		SetError(&errorResult).
		Post("/client/auth/login")

	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		c.Token = result.Token
		return nil
	}

	return fmt.Errorf("error logging in: %s", errorResult.Message.Message)
}

func (c *ConnectAPI) Logout() {
	// TODO: Call the logout endpoint to invalidate the tokens
	c.Token = ""
}

func (c *ConnectAPI) GetStations() ([]Station, error) {
	log.Println("Getting stations")
	client := c.client()
	result := GetStationsResponse{}

	_, err := client.R().
		SetResult(&result).
		Get("/client/stations")

	if err != nil {
		return nil, err
	}

	return result.Stations, nil
}

func (c *ConnectAPI) GetStation(id string) (Station, error) {
	log.Printf("Getting station %s", id)
	client := c.client()
	result := Station{}

	_, err := client.R().
		SetResult(&result).
		Get("/client/stations/" + id)

	if err != nil {
		return Station{}, err
	}

	return result, nil
}

func (c *ConnectAPI) GetTransactionStatistics(stationId string) (TransactionStats, error) {
	log.Printf("Getting transaction statistics for station %s", stationId)
	client := c.client()
	result := TransactionStats{}

	_, err := client.R().
		SetResult(&result).
		SetQueryString("stationId=" + stationId).
		SetResult(&result).
		Get("/client/transactions/statistics")

	if err != nil {
		return TransactionStats{}, err
	}

	return result, nil
}
