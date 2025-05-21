package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/utils"

	"github.com/utmstack/config-client-go/types"
)

const authURL string = "https://id.sophos.com/api/v2/oauth2/token"
const whoamiURL = "https://api.central.sophos.com/whoami/v1"

type SophosCentralProcessor struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	DataRegion   string
	AccessToken  string
	ExpiresAt    time.Time
}

func getSophosCentralProcessor(group types.ModuleGroup) SophosCentralProcessor {
	sophosProcessor := SophosCentralProcessor{}

	for _, cnf := range group.Configurations {
		switch cnf.ConfName {
		case "Client Id":
			sophosProcessor.ClientID = cnf.ConfValue
		case "Client Secret":
			sophosProcessor.ClientSecret = cnf.ConfValue
		}
	}
	return sophosProcessor
}

func (p *SophosCentralProcessor) getAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", p.ClientID)
	data.Set("client_secret", p.ClientSecret)
	data.Set("scope", "token")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Retry logic for getting access token
	maxRetries := 3
	retryDelay := 2 * time.Second

	var response map[string]any
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		response, _, err = utils.DoReq[map[string]any](authURL, []byte(data.Encode()), http.MethodPost, headers)
		if err == nil {
			accessToken, ok := response["access_token"].(string)
			if ok && accessToken != "" {
				expiresIn, ok := response["expires_in"].(float64)
				if ok {
					p.AccessToken = accessToken
					p.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
					return accessToken, nil
				}
			}
		}

		_ = catcher.Error("error getting access token, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		}
	}

	if err != nil {
		return "", catcher.Error("all retries failed when getting access token", err, nil)
	}

	accessToken, ok := response["access_token"].(string)
	if !ok || accessToken == "" {
		return "", catcher.Error("access_token not found in response after all retries", nil, map[string]any{
			"response": response,
		})
	}

	expiresIn, ok := response["expires_in"].(float64)
	if !ok {
		return "", catcher.Error("expires_in not found in response after all retries", nil, map[string]any{
			"response": response,
		})
	}

	p.AccessToken = accessToken
	p.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return accessToken, nil
}

type WhoamiResponse struct {
	ID       string   `json:"id"`
	ApiHosts ApiHosts `json:"apiHosts"`
}
type ApiHosts struct {
	Global     string `json:"global"`
	DataRegion string `json:"dataRegion"`
}

func (p *SophosCentralProcessor) getTenantInfo(accessToken string) error {
	headers := map[string]string{
		"accept":        "application/json",
		"Authorization": "Bearer " + accessToken,
	}

	// Retry logic for getting tenant info
	maxRetries := 3
	retryDelay := 2 * time.Second

	var response WhoamiResponse
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		response, _, err = utils.DoReq[WhoamiResponse](whoamiURL, nil, http.MethodGet, headers)
		if err == nil {
			if response.ID != "" && response.ApiHosts.DataRegion != "" {
				p.TenantID = response.ID
				p.DataRegion = response.ApiHosts.DataRegion
				return nil
			}
		}

		_ = catcher.Error("error getting tenant info, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		}
	}

	if err != nil {
		return catcher.Error("all retries failed when getting tenant info", err, nil)
	}

	if response.ID == "" {
		return catcher.Error("tenant ID not found in whoami response after all retries", nil, map[string]any{
			"response": response,
		})
	}
	p.TenantID = response.ID

	if response.ApiHosts.DataRegion == "" {
		return catcher.Error("dataRegion not found in whoami response after all retries", nil, map[string]any{
			"response": response,
		})
	}
	p.DataRegion = response.ApiHosts.DataRegion

	return nil
}

func (p *SophosCentralProcessor) getValidAccessToken() (string, error) {
	if p.AccessToken != "" && time.Now().Before(p.ExpiresAt) {
		return p.AccessToken, nil
	}
	return p.getAccessToken()
}

type EventAggregate struct {
	Pages Pages            `json:"pages"`
	Items []map[string]any `json:"items"`
}

type Pages struct {
	FromKey string `json:"fromKey"`
	NextKey string `json:"nextKey"`
	Size    int64  `json:"size"`
	MaxSize int64  `json:"maxSize"`
}

func (p *SophosCentralProcessor) getLogs(fromTime int, nextKey string) ([]string, string, error) {
	// Retry logic for getting access token
	maxRetries := 3
	retryDelay := 2 * time.Second

	var accessToken string
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		accessToken, err = p.getValidAccessToken()
		if err == nil {
			break
		}

		_ = catcher.Error("error getting access token, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		}
	}

	if err != nil {
		return nil, "", catcher.Error("all retries failed when getting access token", err, nil)
	}

	if p.TenantID == "" || p.DataRegion == "" {
		// Retry logic for getting tenant info
		for retry := 0; retry < maxRetries; retry++ {
			err = p.getTenantInfo(accessToken)
			if err == nil {
				break
			}

			_ = catcher.Error("error getting tenant info, retrying", err, map[string]any{
				"retry":      retry + 1,
				"maxRetries": maxRetries,
			})

			if retry < maxRetries-1 {
				time.Sleep(retryDelay)
				// Increase delay for next retry
				retryDelay *= 2
			}
		}

		if err != nil {
			return nil, "", catcher.Error("all retries failed when getting tenant info", err, nil)
		}
	}

	logs := make([]string, 0, 1000)

	for {
		u, err := p.buildURL(fromTime, nextKey)
		if err != nil {
			return nil, "", err
		}

		headers := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + accessToken,
			"X-Tenant-ID":   p.TenantID,
		}

		// Retry logic for getting logs
		var response EventAggregate
		for retry := 0; retry < maxRetries; retry++ {
			response, _, err = utils.DoReq[EventAggregate](u.String(), nil, http.MethodGet, headers)
			if err == nil {
				break
			}

			_ = catcher.Error("error getting logs, retrying", err, map[string]any{
				"retry":      retry + 1,
				"maxRetries": maxRetries,
			})

			if retry < maxRetries-1 {
				time.Sleep(retryDelay)
				// Increase delay for next retry
				retryDelay *= 2
			}
		}

		if err != nil {
			return nil, "", catcher.Error("all retries failed when getting logs", err, nil)
		}

		for _, item := range response.Items {
			jsonItem, err := json.Marshal(item)
			if err != nil {
				_ = catcher.Error("error marshalling content details", err, map[string]any{})
				continue
			}
			logs = append(logs, string(jsonItem))
		}

		if response.Pages.NextKey == "" {
			break
		}
		nextKey = response.Pages.NextKey
	}

	return logs, nextKey, nil
}

func (p *SophosCentralProcessor) buildURL(fromTime int, nextKey string) (*url.URL, error) {
	baseURL := p.DataRegion + "/siem/v1/events"
	u, parseErr := url.Parse(baseURL)
	if parseErr != nil {
		return nil, catcher.Error("error parsing url", parseErr, map[string]any{
			"url": baseURL,
		})
	}

	params := url.Values{}
	if nextKey != "" {
		params.Set("pageFromKey", nextKey)
	} else {
		params.Set("from_date", fmt.Sprintf("%d", fromTime))
	}

	u.RawQuery = params.Encode()
	return u, nil
}
