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
		os.Exit(0)
	}

	pluginConfig := plugins.PluginCfg("com.utmstack", false)
	internalKey := pluginConfig.Get("internalKey").String()
	backendUrl := pluginConfig.Get("backend").String()
	config.Backend = backendUrl
	config.InternalKey = internalKey

	client := moduleConf.NewUTMClient(internalKey, backendUrl)

	for {
		if err := utils.ConnectionChecker(GPT_API_ENDPOINT); err != nil {
			_ = catcher.Error("Failed to establish internet connection: %v", err, nil)
		}

		moduleConfig, err := client.GetUTMConfig(enum.SOCAI)
		if err != nil && err.Error() != "" && err.Error() != " " {
			_ = catcher.Error("Error while getting GPT configuration: %v", err, nil)
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
