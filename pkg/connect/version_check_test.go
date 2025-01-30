package connect

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const versionHeaderTemplate = "{\"_id\":\"65c5fc9d6664ff3bb4de7ada\",\"applicationName\":\"Grizzl-E Connect\",\"iosLatestVersion\":\"0.9.1\",\"iosMinimalVersion\":\"%s\",\"androidLatestVersion\":\"0.9.1\",\"androidMinimalVersion\":\"0.9.0\",\"__v\":0}"

func TestCheckVersion(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		expectedResult bool
	}{
		{"SameVersion", "v0.9.1", true},
		{"SameVersionNoV", "0.9.1", true},
		{"NewerVersion", "v1.0.0", false},
		{"NewerVersionNoV", "1.0.0", false},
		{"OlderVersion", "v0.7.0", true},
		{"OlderVersionNoV", "0.7.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ApiVersionSupported(tt.currentVersion)

			require.NoError(t, err, "Error should be nil")
			if result != tt.expectedResult {
				t.Errorf("ApiVersionSupported(%s) = %v; want %v", tt.currentVersion, result, tt.expectedResult)
			}
		})
	}
}

func TestHeaderParse(t *testing.T) {
	// A real X-Application-Version header from the Connect API
	appVersion, err := ParseAppVersionHeader(fmt.Sprintf(versionHeaderTemplate, "0.9.0"))

	require.NoError(t, err, "Error should be nil")
	assert.Equal(t, "Grizzl-E Connect", appVersion.ApplicationName, "ApplicationName should match")
	assert.Equal(t, "0.9.1", appVersion.IosLatestVersion, "IosLatestVersion should match")
	assert.Equal(t, "0.9.0", appVersion.IosMinimalVersion, "IosMinimalVersion should match")
	assert.Equal(t, "0.9.1", appVersion.AndroidLatestVersion, "AndroidLatestVersion should match")
	assert.Equal(t, "0.9.0", appVersion.AndroidMinimalVersion, "AndroidMinimalVersion should match")
}

func TestAssertApiSupporte(t *testing.T) {
	t.Run("API Requires Older", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		err := AssertApiSupported(fmt.Sprintf(versionHeaderTemplate, "0.5.0"))

		require.NoError(t, err, "Error should be nil")
	})

	t.Run("API Requires Newer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		err := AssertApiSupported(fmt.Sprintf(versionHeaderTemplate, "10.9.0"))

		require.Error(t, err, "Error should not be nil")
	})

	t.Run("API Requires Equal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		err := AssertApiSupported(fmt.Sprintf(versionHeaderTemplate, "0.9.1"))

		require.NoError(t, err, "Error should be nil")
	})
}
