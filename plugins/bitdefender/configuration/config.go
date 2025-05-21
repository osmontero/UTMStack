package configuration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"

	"github.com/utmstack/UTMStack/plugins/bitdefender/utils"
	"github.com/utmstack/config-client-go/enum"
	"github.com/utmstack/config-client-go/types"

	UTMStackConfigurationClient "github.com/utmstack/config-client-go"
)

type ModuleConfig struct {
	ConnectionKey string
	AccessUrl     string
	MasterIp      string
	CompaniesIDs  []string
}

const delayCheckConfig = 30 * time.Second

var configsSent = make(map[string]ModuleConfig)

func ConfigureModules(cnf *types.ConfigurationSection, mutex *sync.Mutex) {
	// Recover from panics to ensure the function doesn't terminate
	defer func() {
		if r := recover(); r != nil {
			_ = catcher.Error("recovered from panic in ConfigureModules", nil, map[string]any{
				"panic": r,
			})
			// Restart the function after a brief delay
			time.Sleep(5 * time.Second)
			go ConfigureModules(cnf, mutex)
		}
	}()

	for {
		if err := utils.ConnectionChecker(UrlCheckConnection); err != nil {
			_ = catcher.Error("External connection failure detected: %v", err, nil)
		}
		utmConfig := plugins.PluginCfg("com.utmstack", false)
		internalKey := utmConfig.Get("internalKey").String()
		backendUrl := utmConfig.Get("backend").String()
		client := UTMStackConfigurationClient.NewUTMClient(internalKey, backendUrl)

		time.Sleep(delayCheckConfig)

		// Retry logic for getting module configuration
		maxRetries := 3
		retryDelay := 2 * time.Second
		var tempModuleConfig *types.ConfigurationSection
		var err error

		for retry := 0; retry < maxRetries; retry++ {
			tempModuleConfig, err = client.GetUTMConfig(enum.BITDEFENDER)
			if err == nil {
				break
			}

			if strings.Contains(err.Error(), "invalid character '<'") {
				break // This is a special case, just continue with the main loop
			}

			if (err.Error() != "") && (err.Error() != " ") {
				_ = catcher.Error("error getting configuration of the Bitdefender module, retrying", err, map[string]any{
					"retry":      retry + 1,
					"maxRetries": maxRetries,
				})
			}

			if retry < maxRetries-1 {
				time.Sleep(retryDelay)
				// Increase delay for next retry
				retryDelay *= 2
			}
		}

		if err != nil {
			// If there was an error, continue with the next iteration
			continue
		}

		mutex.Lock()
		*cnf = *tempModuleConfig
		mutex.Unlock()

		for _, group := range (*cnf).ConfigurationGroups {
			isNecessaryConfig := compareConfigs(configsSent, group)
			if isNecessaryConfig {
				if !areAnyEmpty(group.Configurations[0].ConfValue, group.Configurations[1].ConfValue, group.Configurations[2].ConfValue, group.Configurations[3].ConfValue) {
					// Retry logic for API operations
					func() {
						// Recover from panics in this block
						defer func() {
							if r := recover(); r != nil {
								_ = catcher.Error("recovered from panic in API operations", nil, map[string]any{
									"panic": r,
									"group": group.GroupName,
								})
							}
						}()

						if err := apiPush(group, "sendConf"); err != nil {
							_ = catcher.Error("error sending configuration", err, map[string]any{
								"group": group.GroupName,
							})
							return
						}
						time.Sleep(15 * time.Second)
						if err := apiPush(group, "getConf"); err != nil {
							_ = catcher.Error("error getting configuration", err, map[string]any{
								"group": group.GroupName,
							})
							return
						}
						if err := apiPush(group, "sendTest"); err != nil {
							_ = catcher.Error("error sending test event", err, map[string]any{
								"group": group.GroupName,
							})
							return
						}

						configsSent[group.GroupName] = ModuleConfig{
							ConnectionKey: group.Configurations[0].ConfValue,
							AccessUrl:     group.Configurations[1].ConfValue,
							MasterIp:      group.Configurations[2].ConfValue,
							CompaniesIDs:  strings.Split(group.Configurations[3].ConfValue, ","),
						}
					}()
				}
			}
		}
	}
}

func apiPush(config types.ModuleGroup, operation string) error {
	operationFunc := map[string]func(types.ModuleGroup) (*http.Response, error){
		"sendConf": sendPushEventSettings,
		"getConf":  getPushEventSettings,
		"sendTest": sendTestPushEvent,
	}

	fn, ok := operationFunc[operation]
	if !ok {
		return catcher.Error("wrong operation", nil, map[string]any{})
	}

	for i := 0; i < 5; i++ {
		response, err := fn(config)
		if err != nil {
			_ = catcher.Error(fmt.Sprintf("%v", err), err, map[string]any{})
			time.Sleep(1 * time.Minute)
			continue
		}

		func() { _ = response.Body.Close() }()

		return nil
	}

	return catcher.Error("error sending configuration after 5 retries", nil, map[string]any{})
}

func sendPushEventSettings(config types.ModuleGroup) (*http.Response, error) {
	byteTemplate := getTemplateSetPush(config)
	body, err := json.Marshal(byteTemplate)
	if err != nil {
		return nil, catcher.Error("error when marshaling the request body to send the configuration", err, map[string]any{})
	}
	return sendRequest(body, config)
}

func getPushEventSettings(config types.ModuleGroup) (*http.Response, error) {
	byteTemplate := getTemplateGet()
	body, err := json.Marshal(byteTemplate)
	if err != nil {
		return nil, catcher.Error("error when marshaling the request body to get the configuration", err, map[string]any{})
	}
	return sendRequest(body, config)
}

func sendTestPushEvent(config types.ModuleGroup) (*http.Response, error) {
	byteTemplate := getTemplateTest()
	body, err := json.Marshal(byteTemplate)
	if err != nil {
		return nil, catcher.Error("error when marshaling the request body to send the test event", err, map[string]any{})
	}
	return sendRequest(body, config)
}

func areAnyEmpty(strings ...string) bool {
	for _, s := range strings {
		if s == "" || s == " " {
			return true
		}
	}
	return false
}
