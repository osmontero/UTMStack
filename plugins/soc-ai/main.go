package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	twutil "github.com/threatwinds/go-sdk/utils"
	"github.com/utmstack/UTMStack/plugins/soc-ai/configurations"
	"github.com/utmstack/UTMStack/plugins/soc-ai/elastic"
	"github.com/utmstack/UTMStack/plugins/soc-ai/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type socAiServer struct {
	plugins.UnimplementedCorrelationServer
}

func main() {
	utils.Logger.Info("Starting soc-ai plugin...")

	go configurations.UpdateGPTConfigurations()

	// Retry logic for creating socket directory
	maxRetries := 3
	retryDelay := 2 * time.Second
	var socketsFolder twutil.Folder
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		socketsFolder, err = twutil.MkdirJoin(plugins.WorkDir, "sockets")
		if err == nil {
			utils.Logger.LogF(100, "Socket directory %s created", socketsFolder)
			break
		}

		_ = catcher.Error("cannot create socket directory, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		} else {
			// If all retries failed, log the error and return
			_ = catcher.Error("all retries failed when creating socket directory", err, nil)
			return
		}
	}

	socketFile := socketsFolder.FileJoin("com.utmstack.soc-ai_correlation.sock")
	_ = os.Remove(socketFile)

	// Retry logic for resolving unix address
	retryDelay = 2 * time.Second
	var unixAddress *net.UnixAddr

	for retry := 0; retry < maxRetries; retry++ {
		unixAddress, err = net.ResolveUnixAddr("unix", socketFile)
		if err == nil {
			utils.Logger.LogF(100, "Socket file %s created", socketFile)
			break
		}

		_ = catcher.Error("cannot resolve unix address, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		} else {
			// If all retries failed, log the error and return
			_ = catcher.Error("all retries failed when resolving unix address", err, nil)
			return
		}
	}

	// Retry logic for listening to unix socket
	retryDelay = 2 * time.Second
	var listener *net.UnixListener

	for retry := 0; retry < maxRetries; retry++ {
		listener, err = net.ListenUnix("unix", unixAddress)
		if err == nil {
			utils.Logger.LogF(100, "Listening on %s", socketFile)
			break
		}

		_ = catcher.Error("cannot listen to unix socket, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		} else {
			// If all retries failed, log the error and return
			_ = catcher.Error("all retries failed when listening to unix socket", err, nil)
			return
		}
	}

	grpcServer := grpc.NewServer()
	plugins.RegisterCorrelationServer(grpcServer, &socAiServer{})

	// Serve with error handling and retry logic
	retryDelay = 2 * time.Second
	for retry := 0; retry < maxRetries; retry++ {
		err := grpcServer.Serve(listener)
		if err == nil {
			break
		}

		_ = catcher.Error("cannot serve grpc, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		} else {
			// If all retries failed, log the error and return
			_ = catcher.Error("all retries failed when listening to unix socket", err, nil)
			return
		}
	}
}

func (p *socAiServer) Correlate(_ context.Context,
	alert *plugins.Alert) (*emptypb.Empty, error) {
	utils.Logger.LogF(100, "Received alert: %s", alert.Id)

	alertFields := cleanAlerts(alertToAlertFields(alert))

	err := sendRequestToOpenAI(&alertFields)
	if err != nil {
		elastic.RegisterError(err.Error(), alertFields.ID)
		return nil, err
	}

	err = processAlertToElastic(&alertFields)
	if err != nil {
		elastic.RegisterError(err.Error(), alertFields.ID)
		return nil, err
	}

	return nil, nil
}
