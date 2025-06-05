package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/google/uuid"
	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"

	utmconf "github.com/utmstack/config-client-go"
	"github.com/utmstack/config-client-go/enum"
	"github.com/utmstack/config-client-go/types"
)

const (
	defaultTenant      string = "ce66672c-e36d-4761-a8c8-90058fee1a24"
	urlCheckConnection        = "https://login.microsoftonline.com/"
	wait                      = 1 * time.Second
)

func main() {
	mode := plugins.GetCfg().Env.Mode
	if mode != "manager" {
		return
	}

	for t := 0; t < 2*runtime.NumCPU(); t++ {
		go func() {
			plugins.SendLogsFromChannel()
		}()
	}

	delay := 5 * time.Minute
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	startTime := time.Now().UTC().Add(-delay)

	for range ticker.C {
		endTime := time.Now().UTC()

		if err := connectionChecker(urlCheckConnection); err != nil {
			_ = catcher.Error("External connection failure detected: %v", err, nil)
		}

		utmConfig := plugins.PluginCfg("com.utmstack", false)
		internalKey := utmConfig.Get("internalKey").String()
		backendUrl := utmConfig.Get("backend").String()
		if internalKey == "" || backendUrl == "" {
			continue
		}

		client := utmconf.NewUTMClient(internalKey, backendUrl)
		moduleConfig, err := client.GetUTMConfig(enum.AZURE)
		if err == nil && moduleConfig.ModuleActive {
			var wg sync.WaitGroup
			wg.Add(len(moduleConfig.ConfigurationGroups))

			for _, grp := range moduleConfig.ConfigurationGroups {
				go func(group types.ModuleGroup) {
					defer wg.Done()
					var invalid bool
					for _, cnf := range group.Configurations {
						if strings.TrimSpace(cnf.ConfValue) == "" {
							invalid = true
							break
						}
					}
					if !invalid {
						pull(startTime, endTime, group)
					}
				}(grp)
			}

			wg.Wait()
		}

		startTime = endTime.Add(1 * time.Nanosecond)
	}
}

func pull(startTime time.Time, endTime time.Time, group types.ModuleGroup) {
	agent := getAzureProcessor(group)

	// Retry logic for Azure credential creation
	maxRetries := 3
	retryDelay := 2 * time.Second
	var cred *azidentity.ClientSecretCredential
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		cred, err = azidentity.NewClientSecretCredential(agent.TenantID, agent.ClientID, agent.ClientSecretValue, nil)
		if err == nil {
			break
		}

		_ = catcher.Error("cannot obtain Azure credentials, retrying", err, map[string]any{
			"group":      agent.GroupName,
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
		_ = catcher.Error("all retries failed when obtaining Azure credentials", err, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	// Retry logic for Azure client creation
	retryDelay = 2 * time.Second
	var client *azquery.LogsClient

	for retry := 0; retry < maxRetries; retry++ {
		client, err = azquery.NewLogsClient(cred, nil)
		if err == nil {
			break
		}

		_ = catcher.Error("cannot create Logs client, retrying", err, map[string]any{
			"group":      agent.GroupName,
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
		_ = catcher.Error("all retries failed when creating Logs client", err, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	// Create a context with timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Retry logic for Azure query
	retryDelay = 2 * time.Second
	var tables []*azquery.Table
	var queryErr error
	query := fmt.Sprintf(`
		union *
		| where TimeGenerated >= datetime(%s) and TimeGenerated < datetime(%s)
		| order by TimeGenerated desc`,
		startTime.Format(time.RFC3339Nano),
		endTime.Format(time.RFC3339Nano),
	)

	for retry := 0; retry < maxRetries; retry++ {
		resp, err := client.QueryWorkspace(
			ctx,
			agent.WorkspaceID,
			azquery.Body{
				Query: to.Ptr(query),
			},
			nil,
		)

		// Determine the actual error
		if resp.Error != nil {
			queryErr = resp.Error
		} else {
			queryErr = err
		}

		if queryErr == nil {
			tables = resp.Tables
			break
		}

		_ = catcher.Error("cannot query Logs, retrying", queryErr, map[string]any{
			"group":      agent.GroupName,
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		}
	}

	if queryErr != nil {
		_ = catcher.Error("all retries failed when querying Logs", queryErr, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	var logs []map[string]any
	for _, table := range tables {
		for _, row := range table.Rows {
			rowMap := make(map[string]any)
			for i, column := range table.Columns {
				if row[i] != nil {
					if str, ok := row[i].(string); ok && str == "" {
						continue
					}
					rowMap[*column.Name] = row[i]
				}
			}
			logs = append(logs, rowMap)
		}
	}

	if len(logs) > 0 {
		for _, log := range logs {
			jsonLog, err := json.Marshal(log)
			if err != nil {
				_ = catcher.Error("cannot encode log to JSON", err, map[string]any{
					"group": agent.GroupName,
				})
				continue
			}
			plugins.EnqueueLog(&plugins.Log{
				Id:         uuid.New().String(),
				TenantId:   defaultTenant,
				DataType:   "azure",
				DataSource: agent.GroupName,
				Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
				Raw:        string(jsonLog),
			})
		}
	}
}

type AzureConfig struct {
	GroupName         string
	TenantID          string
	ClientID          string
	ClientSecretValue string
	WorkspaceID       string
}

func getAzureProcessor(group types.ModuleGroup) AzureConfig {
	azurePro := AzureConfig{}
	azurePro.GroupName = group.GroupName
	for _, cnf := range group.Configurations {
		switch cnf.ConfKey {
		case "tenantId":
			azurePro.TenantID = cnf.ConfValue
		case "clientId":
			azurePro.ClientID = cnf.ConfValue
		case "clientSecret":
			azurePro.ClientSecretValue = cnf.ConfValue
		case "workspaceId":
			azurePro.WorkspaceID = cnf.ConfValue
		}
	}
	return azurePro
}

func connectionChecker(url string) error {
	checkConn := func() error {
		if err := checkConnection(url); err != nil {
			return fmt.Errorf("connection failed: %v", err)
		}
		return nil
	}

	if err := infiniteRetryIfXError(checkConn, "connection failed"); err != nil {
		return err
	}

	return nil
}

func checkConnection(url string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			_ = catcher.Error("cannot close response body", err, nil)
		}
	}()

	return nil
}

func infiniteRetryIfXError(f func() error, exception string) error {
	var xErrorWasLogged bool

	for {
		err := f()
		if err != nil && is(err, exception) {
			if !xErrorWasLogged {
				_ = catcher.Error("An error occurred (%s), will keep retrying indefinitely...", err, nil)
				xErrorWasLogged = true
			}
			time.Sleep(wait)
			continue
		}

		return err
	}
}

func is(e error, args ...string) bool {
	for _, arg := range args {
		if strings.Contains(e.Error(), arg) {
			return true
		}
	}
	return false
}
