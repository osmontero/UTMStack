package validations

import (
	"net/http"
	"net/url"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/utils"
	"github.com/utmstack/UTMStack/plugins/modules-config/config"
)

const (
	sophosAuthURL = "https://id.sophos.com/api/v2/oauth2/token"
)

func ValidateSophosConfig(config *config.ModuleGroup) error {
	var clientID, clientSecret string
	if config == nil {
		return catcher.Error("Sophos configuration is nil", nil, nil)
	}

	for _, cnf := range config.ModuleGroupConfigurations {
		switch cnf.ConfName {
		case "Client Id":
			clientID = cnf.ConfValue
		case "Client Secret":
			clientSecret = cnf.ConfValue
		}
	}

	if clientID == "" {
		return catcher.Error("Client ID is required in Sophos configuration", nil, nil)
	}
	if clientSecret == "" {
		return catcher.Error("Client Secret is required in Sophos configuration", nil, nil)
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "token")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	response, status, err := utils.DoReq[map[string]any](sophosAuthURL, []byte(data.Encode()), http.MethodPost, headers)
	if err != nil {
		return catcher.Error("error validating Sophos credentials", err, map[string]any{
			"status": status,
		})
	}

	if status != http.StatusOK {
		return catcher.Error("Sophos authentication failed", nil, map[string]any{
			"status":   status,
			"response": response,
		})
	}

	accessToken, ok := response["access_token"].(string)
	if !ok || accessToken == "" {
		return catcher.Error("Sophos credentials are invalid - no access token received", nil, map[string]any{
			"response": response,
			"status":   status,
		})
	}

	return nil
}
