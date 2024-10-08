package connect

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstruct(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "myHost")

	assert.Equal(t, "myUser", c.Username, "Username should be set")
	assert.Equal(t, "myPassword", c.Password, "Password should be set")
}

func TestParseToken(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "myHost")
	c.Token = CreateToken(false)
	token, claims, err := c.ParseToken()

	require.NoError(t, err, "Error should be nil")
	assert.NotNil(t, token, "Token should not be nil")
	assert.Equal(t, "deadbeef", claims.UserId, "UserID should match")
}

func TestLogin(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	err := c.Login()

	require.NoError(t, err, "Error should be nil")
	assert.NotEmpty(t, c.Token, "Token should not be empty")
}

func TestBadLogin(t *testing.T) {
	c := NewConnectAPI("badUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	err := c.Login()

	require.Error(t, err, "Error should not be nil")
	assert.Equal(t, "", c.Token, "Token should be empty")
}

func TestGetStations(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	// Fake token
	c.Token = "deadbeef"

	resp, err := c.GetStations()

	require.NoError(t, err, "Error should be nil")
	assert.Len(t, resp, 1, "Response should have 1 station")
	assert.Equal(t, "station1", resp[0].ID, "Station ID should match")
	assert.Equal(t, "online", resp[0].Status, "Station Status should match")
}

func TestGetStation(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	// Fake token
	c.Token = "deadbeef"

	resp, err := c.GetStation("station1")
	require.NoError(t, err, "Error should be nil")
	assert.Equal(t, "station1", resp.ID, "Station ID should match")

	resp, err = c.GetStation("missing")
	require.NoError(t, err, "Error should be nil")

	resp, err = c.GetStation("errstation")
	require.NoError(t, err, "Error should be nil")
}

func TestTransactionStats(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	// Fake token
	c.Token = "deadbeef"

	statsResp := TransactionStats{
		Sessions:      10,
		AverageEnergy: 19384.19,
		Duration:      10394,
		TopSession:    19482,
		Currency:      "USD",
	}

	httpmock.RegisterResponder("GET", "https://example.com/client/transactions/statistics",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, statsResp)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	resp, err := c.GetTransactionStatistics("station1")

	require.NoError(t, err, "Error should be nil")
	assert.Equal(t, statsResp.Sessions, resp.Sessions, "Sessions should match")
	assert.Equal(t, statsResp.Duration, resp.Duration, "Duration should match")
}

func TestIsExipred(t *testing.T) {
	expiredDate := jwt.NewNumericDate(time.Now().Add(time.Hour * -2))
	futureDate := jwt.NewNumericDate(time.Now().Add(time.Hour))

	assert.False(t, IsExpired(futureDate), "Future date should not be expired")
	assert.True(t, IsExpired(expiredDate), "Expired date should be expired")

	// Test that tokens near expiration are considered expired
	soonExpire := jwt.NewNumericDate(time.Now().Add(time.Second * 10))
	assert.True(t, IsExpired(soonExpire), "Tokens near expiration should be considered expired")
}

func TestTransactionPage(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	resp, err := c.GetTransactions("station1", 2, 0)

	require.NoError(t, err, "Error should be nil")
	assert.Len(t, resp, 2, "Response should have 2 transactions")
}

func TestAllTransactions(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	c.PageSize = 2
	resp, err := c.GetAllTransactions("station1")

	require.NoError(t, err, "Error should be nil")
	assert.Len(t, resp, 4, "Response should have 4 transactions")
}

func TestSingleTransaction(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	resp, err := c.GetTransaction("transaction1")

	require.NoError(t, err, "Error should be nil")
	assert.Equal(t, "transaction1", resp.ID, "Transaction should have requested ID")
}

func TestSingleBadTransaction(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	_, err := c.GetTransaction("bogusId")

	assert.Error(t, err, "Error should not be nil")
}

func TestLogout(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	// Login first
	c.Login()
	assert.NotEmpty(t, c.Token, "Token should not be empty")

	c.Logout()
	assert.Empty(t, c.Token, "Token should be empty after logout")
}

func TestSetDebug(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")

	assert.False(t, c.Client.Debug, "Debug should be false")
	c.SetDebug()
	assert.True(t, c.Client.Debug, "Debug should be true")
}
