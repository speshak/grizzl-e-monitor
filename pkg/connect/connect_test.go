package connect

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// Create a JWT for testing
// If expired is true, the token will be expired
func CreateToken(expired bool) string {
	key := []byte("asb1234")
	t := jwt.New(jwt.SigningMethodHS256)
	s, err := t.SignedString(key)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func SetupHTTPMock() {
	loginRespSuccess := LoginResponse{
		Token:          CreateToken(false),
		IsInitialLogIn: false,
		User: User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jdoe@example.com",
			Phone:     "+1234567890",
			ID:        "123456",
		},
	}

	loginRespError := ApiError{
		StatusCode: 400,
		Timestamp:  "2022-10-10T10:10:10.000Z",
		Path:       "/client/auth/login",
		Message: ApiMessage{
			StatusCode: 400,
			Message:    "Bad username or password",
			Error:      "Bad request",
		},
	}

	httpmock.RegisterResponder("POST", "https://example.com/client/auth/login",
		func(req *http.Request) (*http.Response, error) {
			buf := new(strings.Builder)
			io.Copy(buf, req.Body)

			var err error
			var resp *http.Response

			if strings.Contains(buf.String(), "bad") {
				resp, err = httpmock.NewJsonResponse(400, loginRespError)
			} else {
				resp, err = httpmock.NewJsonResponse(201, loginRespSuccess)
			}

			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	stationListResp := GetStationsResponse{
		Stations: []Station{
			{
				ID:     "station1",
				Status: "online",
			},
		},
	}

	httpmock.RegisterResponder("GET", "https://example.com/client/stations",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, stationListResp)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	stationResp := Station{
		ID:     "station1",
		Status: "online",
	}

	httpmock.RegisterResponder("GET", "https://example.com/client/stations/station1",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, stationResp)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	// An station that returns an bad request
	httpmock.RegisterResponder("GET", "https://example.com/client/stations/errstation",
		httpmock.NewStringResponder(400, ""),
	)

	httpmock.RegisterResponder("GET", "https://example.com/client/stations/missing",
		httpmock.NewStringResponder(404, ""),
	)
}

func TestConstruct(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "myHost")

	assert.Equal(t, "myUser", c.Username, "Username should be set")
	assert.Equal(t, "myPassword", c.Password, "Password should be set")
}

func TestLogin(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	err := c.Login()

	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, c.Token, "Token should not be empty")
}

func TestBadLogin(t *testing.T) {
	c := NewConnectAPI("badUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	err := c.Login()

	assert.Error(t, err, "Error should not be nil")
	assert.Equal(t, "", c.Token, "Token should be empty")
}

func TestGetStations(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())
	SetupHTTPMock()

	// Fake token
	c.Token = "deadbeef"

	resp, err := c.GetStations()

	assert.Nil(t, err, "Error should be nil")
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
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "station1", resp.ID, "Station ID should match")

	resp, err = c.GetStation("missing")
	assert.Nil(t, err, "Error should be nil")

	resp, err = c.GetStation("errstation")
	assert.Nil(t, err, "Error should be nil")
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

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, statsResp.Sessions, resp.Sessions, "Sessions should match")
	assert.Equal(t, statsResp.Duration, resp.Duration, "Duration should match")
}
