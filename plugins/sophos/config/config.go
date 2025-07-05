package config

import (
	"context"
	"strings"
	sync "sync"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	reconnectDelay = 5 * time.Second
	maxMessageSize = 1024 * 1024 * 1024
)

var (
	cnf *ConfigurationSection
	mu  sync.Mutex
)

func ConnectAndStreamConfig(serverAddress, internalKey string) {
	for {
		if err := establishConnection(serverAddress, internalKey); err != nil {
			catcher.Error("Connection attempt failed", err, nil)
			time.Sleep(reconnectDelay)
		}
	}
}

func establishConnection(serverAddress, internalKey string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := createGRPCClient(serverAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	return handleConfigStream(ctx, conn, internalKey)
}

func createGRPCClient(serverAddress string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)),
	)
	if err != nil {
		catcher.Error("Failed to connect to server", err, nil)
		return nil, err
	}

	state := conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		conn.Close()
		catcher.Error("Connection is in shutdown or transient failure state", nil, nil)
		return nil, err
	}

	return conn, nil
}

func handleConfigStream(ctx context.Context, conn *grpc.ClientConn, internalKey string) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "internal-key", internalKey)

	client := NewConfigServiceClient(conn)
	stream, err := client.StreamConfig(ctx)
	if err != nil {
		catcher.Error("Failed to create stream", err, nil)
		return err
	}

	err = stream.Send(&BiDirectionalMessage{
		Payload: &BiDirectionalMessage_PluginInit{
			PluginInit: &PluginInit{Type: PluginType_SOPHOS},
		},
	})
	if err != nil {
		catcher.Error("Failed to send PluginInit", err, nil)
		return err
	}

	for {
		in, err := stream.Recv()
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				catcher.Info("Stream closed by server, reconnecting...", nil)
				return err
			}
			st, ok := status.FromError(err)
			if ok && (st.Code() == codes.Unavailable || st.Code() == codes.Canceled) {
				catcher.Error("Stream error: "+st.Message(), err, nil)
				return err
			} else {
				catcher.Error("Stream receive error", err, nil)
				time.Sleep(reconnectDelay)
				continue
			}
		}

		processConfigMessage(in)
	}
}

func processConfigMessage(message *BiDirectionalMessage) {
	switch msg := message.Payload.(type) {
	case *BiDirectionalMessage_Config:
		cnf = msg.Config
	}
}

func GetConfig() *ConfigurationSection {
	mu.Lock()
	defer mu.Unlock()
	if cnf == nil {
		return &ConfigurationSection{}
	}
	return cnf
}
