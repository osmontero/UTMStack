package main

import (
	"context"
	"net"
	"os"

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

	socketsFolder, err := twutil.MkdirJoin(plugins.WorkDir, "sockets")
	if err != nil {
		utils.Logger.ErrorF("cannot create socket directory: %v", err)
		_ = catcher.Error("cannot create socket directory", err, nil)
		os.Exit(1)
	}

	socketFile := socketsFolder.FileJoin("com.utmstack.soc_ai_correlation.sock")
	_ = os.Remove(socketFile)

	unixAddress, err := net.ResolveUnixAddr("unix", socketFile)
	if err != nil {
		utils.Logger.ErrorF("cannot resolve unix address: %v", err)
		_ = catcher.Error("cannot resolve unix address", err, nil)
		os.Exit(1)
	}

	listener, err := net.ListenUnix("unix", unixAddress)
	if err != nil {
		utils.Logger.ErrorF("cannot listen to unix socket: %v", err)
		_ = catcher.Error("cannot listen to unix socket", err, nil)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	plugins.RegisterCorrelationServer(grpcServer, &socAiServer{})

	if err := grpcServer.Serve(listener); err != nil {
		utils.Logger.ErrorF("cannot serve grpc: %v", err)
		_ = catcher.Error("cannot serve grpc", err, nil)
		os.Exit(1)
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
