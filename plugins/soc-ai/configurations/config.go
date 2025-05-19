package configurations

import (
	"os"
	"sync"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/utmstack/UTMStack/plugins/soc-ai/utils"
	moduleConf "github.com/utmstack/config-client-go"
	"github.com/utmstack/config-client-go/enum"
)

var (
	config     Config
	configOnce sync.Once
)

type Config struct {
	Backend                   string
	InternalKey               string
	Openseach                 string
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

func UpdateGPTConfigurations() {
	GetConfig()

	mode := plugins.GetCfg().Env.Mode
	if mode != "manager" {
		utils.Logger.ErrorF("Plugin is not running in manager mode, exiting...")
		os.Exit(0)
	}

	pluginConfig := plugins.PluginCfg("com.utmstack", false)
	if !pluginConfig.Exists() {
		utils.Logger.ErrorF("Plugin configuration not found")
		_ = catcher.Error("plugin configuration not found", nil, nil)
		os.Exit(1)
	}

	config.Backend = pluginConfig.Get("backend").String()
	config.InternalKey = pluginConfig.Get("internalKey").String()
	config.Openseach = pluginConfig.Get("opensearch").String()

	client := moduleConf.NewUTMClient(config.InternalKey, config.Backend)

	for {
		if err := utils.ConnectionChecker(GPT_API_ENDPOINT); err != nil {
			utils.Logger.ErrorF("Failed to establish internet connection: %v", err)
			_ = catcher.Error("Failed to establish internet connection: %v", err, nil)
		}

		moduleConfig, err := client.GetUTMConfig(enum.SOCAI)
		if err != nil && err.Error() != "" && err.Error() != " " {
			utils.Logger.LogF(100, "Error while getting module configuration: %v", err)
			time.Sleep(TIME_FOR_GET_CONFIG * time.Second)
			continue
		}
		if moduleConfig == nil {
			time.Sleep(TIME_FOR_GET_CONFIG * time.Second)
			continue
		}

		config.ModuleActive = moduleConfig.ModuleActive

		if config.ModuleActive && len(moduleConfig.ConfigurationGroups) > 0 {
			for _, c := range moduleConfig.ConfigurationGroups[0].Configurations {
				switch c.ConfKey {
				case "utmstack.socai.key":
					if c.ConfValue != "" && c.ConfValue != " " {
						config.APIKey = c.ConfValue
					}
				case "utmstack.socai.incidentCreation":
					if c.ConfValue != "" && c.ConfValue != " " {
						config.AutomaticIncidentCreation = c.ConfValue == "true"
					}
				case "utmstack.socai.changeAlertStatus":
					if c.ConfValue != "" && c.ConfValue != " " {
						config.ChangeAlertStatus = c.ConfValue == "true"
					}
				case "utmstack.socai.model":
					if c.ConfValue != "" && c.ConfValue != " " {
						config.Model = c.ConfValue
					}
				}
			}
		}

		time.Sleep(TIME_FOR_GET_CONFIG * time.Second)
	}
}
