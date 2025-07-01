package connect

import (
	"encoding/json"
	"fmt"

	"golang.org/x/mod/semver"
)

/**
* Functions that handle parsing & comparing the Connect API version information.
*
* Connect API responses include an X-Application-Version header that contains a
* JSON object with the latest and minimal versions of the iOS and Android apps.
* We can use that information to compare the versions of the app we're emulating
* to try to detect breaking API changes or other issues.
 */

// The version of the app that was used when capturing traffic to reverse
// engineer the API.
const EmulatedAppVersion = "v0.9.2"

/* JSON structure of the X-Application-Version header */
type AppVersion struct {
	ID                    string `json:"_id"`
	ApplicationName       string `json:"applicationName"`
	IosLatestVersion      string `json:"iosLatestVersion"`
	IosMinimalVersion     string `json:"iosMinimalVersion"`
	AndroidLatestVersion  string `json:"androidLatestVersion"`
	AndroidMinimalVersion string `json:"androidMinimalVersion"`
}

func ParseAppVersionHeader(header string) (AppVersion, error) {
	var appVersion AppVersion
	err := json.Unmarshal([]byte(header), &appVersion)
	return appVersion, err
}

func ApiVersionSupported(apiVersion string) (bool, error) {
	// The API version is a semver string, but the API header doesn't include the v prefix
	// so we need to add it to compare it correctly.
	if apiVersion[0] != 'v' {
		apiVersion = "v" + apiVersion
	}

	if !semver.IsValid(apiVersion) {
		return false, fmt.Errorf("invalid API version: %s", apiVersion)
	}

	res := semver.Compare(EmulatedAppVersion, apiVersion)
	return res > -1, nil
}

/**
 * Given the X-Application-Version header from the Connect API, check if we can work with the API.
 */
func AssertApiSupported(header string) error {
	appVersion, err := ParseAppVersionHeader(header)
	if err != nil {
		return fmt.Errorf("failed to parse X-Application-Version header: %w", err)
	}

	// Check if the API version is supported
	// We use captures of the app on iOS, so compare against the iOS minimal version
	supported, err := ApiVersionSupported(appVersion.IosMinimalVersion)
	if err != nil {
		return fmt.Errorf("failed to compare API versions: %w", err)
	}

	if !supported {
		return fmt.Errorf("API version %s is not supported, minimal version is v%s", EmulatedAppVersion, appVersion.IosMinimalVersion)
	}

	return nil
}
