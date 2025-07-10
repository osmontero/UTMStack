package config

import (
	"context"
	"log"
	"os"
	"strings"
	sync "sync"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/utmstack/UTMStack/plugins/soc-ai/utils"
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
	grpcConfig *ConfigurationSection
	grpcMutex  sync.RWMutex

	config      Config
	configMutex sync.RWMutex
	configOnce  sync.Once
)

type Config struct {
	Backend                   string
	InternalKey               string
	Opensearch                string
	APIKey                    string
	ChangeAlertStatus         bool
	AutomaticIncidentCreation bool
	Model                     string
	ModuleActive              bool
}

func GetConfig() *Config {
	configOnce.Do(func() {
		config = Config{}
	})
	return &config
}

func StartConfigurationSystem() {
	GetConfig()

	pluginConfig := plugins.PluginCfg("com.utmstack", false)
	if !pluginConfig.Exists() {
		_ = catcher.Error("plugin configuration not found", nil, nil)
		os.Exit(1)
	}

	configMutex.Lock()
	config.Backend = pluginConfig.Get("backend").String()
	config.InternalKey = pluginConfig.Get("internalKey").String()
	config.Opensearch = pluginConfig.Get("opensearch").String()
	configMutex.Unlock()

	utils.Logger.Info("Starting gRPC configuration client...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ConnectAndStreamConfig("localhost:9003", config.InternalKey)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if GetAPIKey() != "" {
				if err := utils.ConnectionChecker(GPT_API_ENDPOINT); err != nil {
					_ = catcher.Error("Failed to establish internet connection", err, nil)
				}
			}

			currentGRPCConfig := GetGRPCConfig()
			if currentGRPCConfig != nil && currentGRPCConfig.Id != 0 {
				updateConfigFromGRPC(currentGRPCConfig)
			} else {
				utils.Logger.LogF(100, "No gRPC configuration available yet...")
			}
		case <-ctx.Done():
			return
		}
	}
}

func ConnectAndStreamConfig(serverAddress, internalKey string) {
	for {
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctx = metadata.AppendToOutgoingContext(ctx, "internal-key", internalKey)
			conn, err := grpc.NewClient(
				serverAddress,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)),
			)

			if err != nil {
				catcher.Error("Failed to connect to server", err, nil)
				return
			}

			state := conn.GetState()
			if state == connectivity.Shutdown || state == connectivity.TransientFailure {
				catcher.Error("Connection is in shutdown or transient failure state", nil, nil)
				return
			}

			client := NewConfigServiceClient(conn)
			stream, err := client.StreamConfig(ctx)
			if err != nil {
				catcher.Error("Failed to create stream", err, nil)
				return
			}

			err = stream.Send(&BiDirectionalMessage{
				Payload: &BiDirectionalMessage_PluginInit{
					PluginInit: &PluginInit{Type: PluginType_SOC_AI},
				},
			})
			if err != nil {
				catcher.Error("Failed to send PluginInit", err, nil)
				return
			}

			for {
				in, err := stream.Recv()
				if err != nil {
					if strings.Contains(err.Error(), "EOF") {
						catcher.Info("Stream closed by server, reconnecting...", nil)
						conn.Close()
						time.Sleep(reconnectDelay)
						break
					}
					st, ok := status.FromError(err)
					if ok && (st.Code() == codes.Unavailable || st.Code() == codes.Canceled) {
						catcher.Error("Stream error: "+st.Message(), err, nil)
						conn.Close()
						time.Sleep(reconnectDelay)
						break
					} else {
						catcher.Error("Stream receive error", err, nil)
						time.Sleep(reconnectDelay)
						continue
					}
				}

				switch message := in.Payload.(type) {
				case *BiDirectionalMessage_Config:
					log.Printf("Received configuration update: %v", message.Config)
					grpcMutex.Lock()
					grpcConfig = message.Config
					grpcMutex.Unlock()
				}
			}
		}()

		time.Sleep(reconnectDelay)
	}
}

func GetGRPCConfig() *ConfigurationSection {
	grpcMutex.RLock()
	defer grpcMutex.RUnlock()
	if grpcConfig == nil {
		return &ConfigurationSection{}
	}
	return grpcConfig
}

func updateConfigFromGRPC(grpcConf *ConfigurationSection) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if config.ModuleActive != grpcConf.ModuleActive {
		utils.Logger.LogF(100, "Module active status changed: %v -> %v",
			config.ModuleActive, grpcConf.ModuleActive)
	}

	config.ModuleActive = grpcConf.ModuleActive

	if !config.ModuleActive {
		utils.Logger.Info("SOC-AI module is disabled")
		return
	}

	if len(grpcConf.ModuleGroups) == 0 {
		utils.Logger.LogF(100, "No module groups found in gRPC configuration")
		return
	}

	for _, group := range grpcConf.ModuleGroups {
		utils.Logger.LogF(100, "Processing configuration group: %s", group.GroupName)

		for _, c := range group.ModuleGroupConfigurations {
			oldValue := GetConfigValue(c.ConfKey)

			switch c.ConfKey {
			case "utmstack.socai.key":
				if c.ConfValue != "" && c.ConfValue != " " && c.ConfValue != config.APIKey {
					config.APIKey = c.ConfValue
					utils.Logger.LogF(100, "Updated API Key from gRPC config")
				}
			case "utmstack.socai.incidentCreation":
				if c.ConfValue != "" && c.ConfValue != " " {
					newValue := c.ConfValue == "true"
					if newValue != config.AutomaticIncidentCreation {
						config.AutomaticIncidentCreation = newValue
						utils.Logger.LogF(100, "Updated incident creation setting: %v -> %v",
							oldValue, newValue)
					}
				}
			case "utmstack.socai.changeAlertStatus":
				if c.ConfValue != "" && c.ConfValue != " " {
					newValue := c.ConfValue == "true"
					if newValue != config.ChangeAlertStatus {
						config.ChangeAlertStatus = newValue
						utils.Logger.LogF(100, "Updated alert status change setting: %v -> %v",
							oldValue, newValue)
					}
				}
			case "utmstack.socai.model":
				if c.ConfValue != "" && c.ConfValue != " " && c.ConfValue != config.Model {
					config.Model = c.ConfValue
					utils.Logger.LogF(100, "Updated GPT model: %s -> %s", oldValue, c.ConfValue)
				}
			default:
				utils.Logger.LogF(100, "Unknown configuration key: %s", c.ConfKey)
			}
		}
	}

	utils.Logger.LogF(100, "Configuration updated from gRPC - Active: %v, API Key: %s, Model: %s, IncidentCreation: %v, ChangeAlertStatus: %v",
		config.ModuleActive,
		maskAPIKey(config.APIKey),
		config.Model,
		config.AutomaticIncidentCreation,
		config.ChangeAlertStatus)
}

func GetConfigValue(key string) interface{} {
	switch key {
	case "utmstack.socai.key":
		return maskAPIKey(config.APIKey)
	case "utmstack.socai.incidentCreation":
		return config.AutomaticIncidentCreation
	case "utmstack.socai.changeAlertStatus":
		return config.ChangeAlertStatus
	case "utmstack.socai.model":
		return config.Model
	default:
		return nil
	}
}

func maskAPIKey(apiKey string) string {
	if apiKey == "" {
		return "not set"
	}
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}

func GetAPIKey() string {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config.APIKey
}
