package connect

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestConstruct(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "myHost")

	assert.Equal(t, "myUser", c.Username, "Username should be set")
	assert.Equal(t, "myPassword", c.Password, "Password should be set")
}

func TestLogin(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	// Response data
	loginResp := LoginResponse{
		Token:          "myToken",
		IsInitialLogIn: false,
	}

	httpmock.RegisterResponder("POST", "https://example.com/client/auth/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(201, loginResp)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	c.Login()

	assert.Equal(t, loginResp.Token, c.Token, "Token should be set")
}

func TestGetStations(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	// Fake token
	c.Token = "deadbeef"

	stationResp := []Station{
		Station{
			ID:     "station1",
			Status: "online",
		},
	}

	httpmock.RegisterResponder("GET", "https://example.com/client/stations",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, stationResp)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	resp, err := c.GetStations()

	assert.Nil(t, err, "Error should be nil")
	assert.Len(t, resp, 1, "Response should have 1 station")
	assert.Equal(t, stationResp[0].ID, resp[0].ID, "Station ID should match")
	assert.Equal(t, stationResp[0].Status, resp[0].Status, "Station Status should match")
}

func TestGetStation(t *testing.T) {
	c := NewConnectAPI("myUser", "myPassword", "https://example.com")
	httpmock.ActivateNonDefault(c.Client.GetClient())

	// Fake token
	c.Token = "deadbeef"

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

	resp, err := c.GetStation("station1")

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, stationResp.ID, resp.ID, "Station ID should match")
}
