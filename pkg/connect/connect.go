package connect

import (
	"github.com/go-resty/resty/v2"
)

type ConnectAPI struct {
	Username string
	Password string
	Token    string
	Client   *resty.Client
}

func NewConnectAPI(username, password, host string) *ConnectAPI {
	client := resty.New()
	client.
		SetBaseURL(host).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "GrizzlEConnect/66 CFNetwork/1568.100.1 Darwin/24.0.0").
		SetHeader("x-app-client", "Apple, iPad14,3, iPadOS 18.0")

	return &ConnectAPI{
		Username: username,
		Password: password,
		Client:   client,
	}
}

func (c *ConnectAPI) client() *resty.Client {
	return c.Client.
		SetAuthToken(c.Token)
}

func (c *ConnectAPI) Login() error {
	// Get login token for future requests
	result := LoginResponse{}

	_, err := c.Client.R().
		SetBody(map[string]interface{}{
			"emailOrPhone": c.Username,
			"password":     c.Password,
		}).
		SetResult(&result).
		Post("/client/auth/login")

	if err != nil {
		return err
	}
	c.Token = result.Token
	return nil
}

func (c *ConnectAPI) Logout() {
	// TODO: Call the logout endpoint to invalidate the tokens
	c.Token = ""
}

func (c *ConnectAPI) GetStations() ([]Station, error) {
	client := c.client()
	result := []Station{}

	_, err := client.R().
		SetResult(&result).
		Get("/client/stations")

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ConnectAPI) GetStation(id string) (Station, error) {
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
