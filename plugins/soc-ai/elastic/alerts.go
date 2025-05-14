package elastic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/utmstack/UTMStack/plugins/soc-ai/configurations"
	"github.com/utmstack/UTMStack/plugins/soc-ai/schema"
	"github.com/utmstack/UTMStack/plugins/soc-ai/utils"
)

func ChangeAlertStatus(id string, status int, observations string) error {
	url := configurations.GetConfig().Backend + configurations.API_ALERT_STATUS_ENDPOINT
	headers := map[string]string{
		"Content-Type":     "application/json",
		"Utm-Internal-Key": configurations.GetConfig().InternalKey,
	}

	body := schema.ChangeAlertStatus{AlertIDs: []string{id}, Status: status, StatusObservation: observations}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}

	resp, statusCode, err := utils.DoReq(url, bodyBytes, "POST", headers, configurations.HTTP_TIMEOUT)
	if err != nil || statusCode != http.StatusOK {
		return fmt.Errorf("error while doing request: %v, status: %d, response: %v", err, statusCode, string(resp))
	}

	return nil
}
