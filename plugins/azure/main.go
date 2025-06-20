package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
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

	for range ticker.C {
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
						pull(group)
					}
				}(grp)
			}

			wg.Wait()
		}

	}
}

func pull(group types.ModuleGroup) {
	agent := getAzureProcessor(group)

	if agent.EventHubConnection == "" || agent.ConsumerGroup == "" ||
		agent.StorageContainer == "" || agent.StorageConnection == "" {
		_ = catcher.Error("missing required configuration for Event Hub", nil, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	eventHubParts := strings.Split(agent.EventHubConnection, ";EntityPath=")
	if len(eventHubParts) != 2 {
		_ = catcher.Error("invalid Event Hub connection string format", nil, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	eventHubConnection := eventHubParts[0]
	eventHubName := eventHubParts[1]

	blobClient, err := azblob.NewClientFromConnectionString(agent.StorageConnection, nil)
	if err != nil {
		_ = catcher.Error("cannot create blob client", err, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	checkpointStore, err := checkpoints.NewBlobStore(
		blobClient.ServiceClient().NewContainerClient(agent.StorageContainer), nil)
	if err != nil {
		_ = catcher.Error("cannot create checkpoint store", err, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	maxRetries := 3
	retryDelay := 2 * time.Second
	var client *azeventhubs.ConsumerClient

	for retry := 0; retry < maxRetries; retry++ {
		client, err = azeventhubs.NewConsumerClientFromConnectionString(
			eventHubConnection, eventHubName, agent.ConsumerGroup, nil)
		if err == nil {
			break
		}

		_ = catcher.Error("cannot create Event Hub consumer client, retrying", err, map[string]any{
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
		_ = catcher.Error("all retries failed when creating Event Hub consumer client", err, map[string]any{
			"group": agent.GroupName,
		})
		return
	}
	defer client.Close(context.Background())

	processor, err := azeventhubs.NewProcessor(client, checkpointStore, nil)
	if err != nil {
		_ = catcher.Error("cannot create Event Hub processor", err, map[string]any{
			"group": agent.GroupName,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	go func() {
		for {
			pc := processor.NextPartitionClient(ctx)
			if pc == nil {
				return
			}
			go processPartition(pc, agent.GroupName)
		}
	}()

	if err := processor.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		_ = catcher.Error("error running Event Hub processor", err, map[string]any{
			"group": agent.GroupName,
		})
	}
}

func processPartition(pc *azeventhubs.ProcessorPartitionClient, groupName string) {
	defer pc.Close(context.Background())

	for {
		recvCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		events, err := pc.ReceiveEvents(recvCtx, 100, nil)
		cancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			_ = catcher.Error("error receiving events", err, map[string]any{
				"group":       groupName,
				"partitionID": pc.PartitionID(),
			})
			return
		}

		if len(events) == 0 {
			continue
		}

		for _, event := range events {
			var logData map[string]any
			if err := json.Unmarshal(event.Body, &logData); err != nil {
				_ = catcher.Error("cannot parse event body", err, map[string]any{
					"group":       groupName,
					"partitionID": pc.PartitionID(),
				})
				continue
			}

			jsonLog, err := json.Marshal(logData)
			if err != nil {
				_ = catcher.Error("cannot encode log to JSON", err, map[string]any{
					"group":       groupName,
					"partitionID": pc.PartitionID(),
				})
				continue
			}

			plugins.EnqueueLog(&plugins.Log{
				Id:         uuid.New().String(),
				TenantId:   defaultTenant,
				DataType:   "azure",
				DataSource: groupName,
				Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
				Raw:        string(jsonLog),
			})
		}

		if err := pc.UpdateCheckpoint(context.Background(), events[len(events)-1], nil); err != nil {
			_ = catcher.Error("checkpoint error", err, map[string]any{
				"group":       groupName,
				"partitionID": pc.PartitionID(),
			})
		}
	}
}

type AzureConfig struct {
	GroupName          string
	EventHubConnection string
	ConsumerGroup      string
	StorageContainer   string
	StorageConnection  string
}

func getAzureProcessor(group types.ModuleGroup) AzureConfig {
	azurePro := AzureConfig{}
	azurePro.GroupName = group.GroupName
	for _, cnf := range group.Configurations {
		switch cnf.ConfKey {
		case "eventHubConnection":
			azurePro.EventHubConnection = cnf.ConfValue
		case "consumerGroup":
			azurePro.ConsumerGroup = cnf.ConfValue
		case "storageContainer":
			azurePro.StorageContainer = cnf.ConfValue
		case "storageConnection":
			azurePro.StorageConnection = cnf.ConfValue
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
