package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var LogsChan = make(chan *plugins.Log)

func ProcessLogs() {
	// Recover from panics to ensure the function doesn't terminate
	defer func() {
		if r := recover(); r != nil {
			_ = catcher.Error("recovered from panic in ProcessLogs", nil, map[string]any{
				"panic": r,
			})
			// Restart the function after a brief delay
			time.Sleep(5 * time.Second)
			go ProcessLogs()
		}
	}()

	// Create a context with cancel for controlling the sender goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is canceled when function returns

	// Initialize with retry logic instead of exiting
	var conn *grpc.ClientConn
	var client plugins.EngineClient
	var inputClient plugins.Engine_InputClient
	var err error
	var socketFile string

	// Retry loop for initialization
	for {
		socketsPath, err := utils.MkdirJoin(plugins.WorkDir, "sockets")
		if err != nil {
			_ = catcher.Error("cannot create socket directory", err, nil)
			time.Sleep(5 * time.Second)
			continue
		}

		socketFile = socketsPath.FileJoin("engine_server.sock")

		conn, err = grpc.NewClient(
			fmt.Sprintf("unix://%s", socketFile),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			_ = catcher.Error("failed to connect to engine server", err, map[string]any{})
			time.Sleep(5 * time.Second)
			continue
		}

		client = plugins.NewEngineClient(conn)

		inputClient, err = client.Input(context.Background())
		if err != nil {
			_ = catcher.Error("failed to create input client", err, map[string]any{})
			if conn != nil {
				_ = conn.Close()
			}
			time.Sleep(5 * time.Second)
			continue
		}

		// If we got here, initialization was successful
		break
	}

	// Sender goroutine with its own panic recovery
	go func() {
		// Recover from panics to ensure goroutine doesn't terminate
		defer func() {
			if r := recover(); r != nil {
				_ = catcher.Error("recovered from panic in log sender", nil, map[string]any{
					"panic": r,
				})
				// Let the main goroutine handle reconnection
				panic(r) // Re-panic to be caught by the outer recover
			}
		}()

		for {
			select {
			case <-ctx.Done():
				// Context was canceled, exit the goroutine
				return
			case log := <-LogsChan:
				err := inputClient.Send(log)
				if err != nil {
					_ = catcher.Error("failed to send log", err, map[string]any{})
					// If there's a send error, we might need to reconnect
					panic(fmt.Errorf("send error, need to reconnect: %v", err))
				}
			}
		}
	}()

	// Receiver loop with error handling
	for {
		_, err = inputClient.Recv()
		if err != nil {
			_ = catcher.Error("failed to receive ack", err, map[string]any{})
			// If there's a receive error, we might need to reconnect
			if conn != nil {
				_ = conn.Close()
			}
			// Cancel the context to stop the sender goroutine
			cancel()
			// Restart the entire process
			time.Sleep(5 * time.Second)
			go ProcessLogs()
			return
		}
	}
}
