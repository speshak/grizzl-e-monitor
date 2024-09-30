package connect

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
)

// A Connect API client
type ConnectAPI interface {
	SetDebug()
	AssertValidToken() error
	Login() error
	Logout() error
	GetStations() ([]Station, error)
	GetStation(id string) (Station, error)
	GetTransactionStatistics(stationId string) (TransactionStats, error)
	GetAllTransactions(stationId string) ([]Transaction, error)
	GetTransactions(stationId string, limit int, offset int) ([]Transaction, error)
	GetTransaction(id string) (Transaction, error)
}

type ConnectAPIClient struct {
	Username string
	Password string
	Token    string
	Client   *resty.Client
	PageSize int
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Iat           int64  `json:"iat"`
	UserId        string `json:"userId"`
	UserSessionId string `json:"userSessionId"`
}

func NewConnectAPI(username, password, host string) *ConnectAPIClient {
	client := resty.New()
	client.
		EnableTrace().
		SetBaseURL(host).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "GrizzlEConnect/87 CFNetwork/1568.100.1 Darwin/24.0.0").
		SetHeader("x-app-client", "Apple, iPad14,3, iPadOS 18.0").
		SetHeader("x-app-version", "v0.8.0 (87)").
		SetHeader("x-application-name", "Grizzl-E Connect")

	return &ConnectAPIClient{
		Username: username,
		Password: password,
		Client:   client,
		PageSize: 10,
	}
}

// Enable debug mode in the underlying resty client
func (c *ConnectAPIClient) SetDebug() {
	c.Client = c.Client.SetDebug(true)
}

// Ensure the login token is valid
func (c *ConnectAPIClient) AssertValidToken() error {
	parser := jwt.NewParser()
	claims := TokenClaims{}

	// We don't need to verify the token, just parse it
	jwtToken, _, err := parser.ParseUnverified(c.Token, &claims)

	// We only care about the error if we have a token
	if err != nil && c.Token != "" {
		log.Printf("Error parsing token: %s", err)
		log.Printf("Token: %s", c.Token)
	}

	if jwtToken == nil || !jwtToken.Valid || IsExpired(claims.ExpiresAt) {
		log.Println("No valid token, logging in")
		err := c.Login()
		if err != nil {
			log.Fatalf("Error logging in: %s", err)
		}
	}

	return nil
}

// Check if a token is expired
func IsExpired(expires *jwt.NumericDate) bool {
	// Include a 30 second buffer to account for clock skew
	return time.Until(expires.Time) < 30*time.Second
}

// Get a resty client with the auth token set
// This will log in if no token is set or the token is expired
func (c *ConnectAPIClient) client() *resty.Client {
	c.AssertValidToken()

	return c.Client.
		SetAuthToken(c.Token)
}

func (c *ConnectAPIClient) Login() error {
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

func (c *ConnectAPIClient) Logout() error {
	// TODO: Call the logout endpoint to invalidate the tokens
	c.Token = ""
	return nil
}

func (c *ConnectAPIClient) GetStations() ([]Station, error) {
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

func (c *ConnectAPIClient) GetStation(id string) (Station, error) {
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

func (c *ConnectAPIClient) GetTransactionStatistics(stationId string) (TransactionStats, error) {
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

func (c *ConnectAPIClient) GetAllTransactions(stationId string) ([]Transaction, error) {
	log.Printf("Getting all transactions for station %s", stationId)
	var transactions []Transaction

	offset := 0

	// Get pages of transactions until we get back less than the limit
	for {
		page, err := c.GetTransactions(stationId, c.PageSize, offset)

		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, page...)

		if len(page) < c.PageSize {
			return transactions, nil
		}
		offset++
	}
}

// Get a single page of transactions, defined by the limit and offset
func (c *ConnectAPIClient) GetTransactions(stationId string, limit int, offset int) ([]Transaction, error) {
	log.Printf("Getting transactions for station %s", stationId)
	client := c.client()
	result := GetTransactionsResponse{}

	_, err := client.R().
		SetQueryParams(map[string]string{
			"stationId": stationId,
			"limit":     strconv.Itoa(limit),
			"offset":    strconv.Itoa(offset),
		}).
		SetResult(&result).
		Get("/client/transactions")

	if err != nil {
		return nil, err
	}

	return result.Transactions, nil
}

func (c *ConnectAPIClient) GetTransaction(id string) (Transaction, error) {
	log.Printf("Getting transaction %s", id)
	client := c.client()
	result := Transaction{}

	_, err := client.R().
		SetResult(&result).
		Get("/client/transactions/" + id)

	if err != nil {
		return Transaction{}, err
	}

	return result, nil
}
