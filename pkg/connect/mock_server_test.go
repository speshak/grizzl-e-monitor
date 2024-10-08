package connect

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jarcoal/httpmock"
)

// Create a JWT for testing
// If expired is true, the token will be expired
func CreateToken(expired bool) string {
	key := []byte("asb1234")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp":           1729274813,
			"iat":           1726682813,
			"userId":        "deadbeef",
			"userSessionId": "cafecafe",
		})
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

	appVersion, _ := json.Marshal(ApplicationVersion{
		Id:                    "65c5fc9d6664ff3bb4de7ada",
		ApplicationName:       "Grizzl-E Connect",
		IosLatestVersion:      "0.7.5",
		IosMinimalVersion:     "0.7.0",
		AndroidLatestVersion:  "0.7.5",
		AndroidMinimalVersion: "0.7.0",
	})

	versionHeader := http.Header{
		"X-Application-Version": []string{string(appVersion)},
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

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/client/stations",
		httpmock.NewJsonResponderOrPanic(200, stationListResp).HeaderAdd(versionHeader),
	)

	stationResp := Station{
		ID:     "station1",
		Status: "online",
	}

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/client/stations/station1",
		httpmock.NewJsonResponderOrPanic(200, stationResp).HeaderAdd(versionHeader),
	)

	// An station that returns an bad request
	httpmock.RegisterResponder("GET", "https://example.com/client/stations/errstation",
		httpmock.NewStringResponder(400, "").HeaderAdd(versionHeader),
	)

	httpmock.RegisterResponder("GET", "https://example.com/client/stations/missing",
		httpmock.NewStringResponder(404, "").HeaderAdd(versionHeader),
	)

	transactionPage1 := GetTransactionsResponse{
		Transactions: []Transaction{
			{
				ID:      "transaction1",
				Station: "station1",
			},
			{
				ID:      "transaction2",
				Station: "station1",
			},
		},
	}

	transactionPage2 := GetTransactionsResponse{
		Transactions: []Transaction{
			{
				ID:      "transaction3",
				Station: "station1",
			},
			{
				ID:      "transaction4",
				Station: "station1",
			},
		},
	}

	transactionPage3 := GetTransactionsResponse{
		Transactions: []Transaction{},
	}

	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://example.com/client/transactions",
		map[string]string{
			"stationId": "station1",
			"limit":     "2",
			"offset":    "0",
		},
		httpmock.NewJsonResponderOrPanic(200, transactionPage1).HeaderAdd(versionHeader),
	)

	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://example.com/client/transactions",
		map[string]string{
			"stationId": "station1",
			"limit":     "2",
			"offset":    "1",
		},
		httpmock.NewJsonResponderOrPanic(200, transactionPage2).HeaderAdd(versionHeader),
	)

	httpmock.RegisterResponderWithQuery(
		"GET",
		"https://example.com/client/transactions",
		map[string]string{
			"stationId": "station1",
			"limit":     "2",
			"offset":    "2",
		},
		httpmock.NewJsonResponderOrPanic(200, transactionPage3).HeaderAdd(versionHeader),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/client/transactions/transaction1",
		httpmock.NewJsonResponderOrPanic(200, GetTransactionResponse{
			Transaction: Transaction{
				ID: "transaction1",
			}}).HeaderAdd(versionHeader),
	)
}
