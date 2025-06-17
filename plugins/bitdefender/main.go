package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"

	"github.com/utmstack/UTMStack/plugins/bitdefender/configuration"
	"github.com/utmstack/UTMStack/plugins/bitdefender/processor"
	"github.com/utmstack/UTMStack/plugins/bitdefender/server"
	"github.com/utmstack/config-client-go/types"
)

var (
	mutex        = &sync.Mutex{}
	moduleConfig = types.ConfigurationSection{}
)

func main() {
	// Recover from panics to ensure the main function doesn't terminate
	defer func() {
		if r := recover(); r != nil {
			_ = catcher.Error("recovered from panic in main function", nil, map[string]any{
				"panic": r,
			})
			// Restart the main function after a brief delay
			time.Sleep(5 * time.Second)
			go main()
		}
	}()

	mode := plugins.GetCfg().Env.Mode
	if mode != "manager" {
		// Don't exit, just return
		return
	}

	// Retry logic for loading certificates
	maxRetries := 3
	retryDelay := 2 * time.Second
	var cert, key string
	var err error

	for retry := 0; retry < maxRetries; retry++ {
		cert, key, err = loadCerts()
		if err == nil {
			break
		}

		_ = catcher.Error("cannot load certificates, retrying", err, map[string]any{
			"retry":      retry + 1,
			"maxRetries": maxRetries,
		})

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			// Increase delay for next retry
			retryDelay *= 2
		} else {
			// If all retries failed, log the error and return
			_ = catcher.Error("all retries failed when loading certificates", err, nil)
			return
		}
	}

	server.StartServer(&moduleConfig, cert, key)

	go configuration.ConfigureModules(&moduleConfig, mutex)

	go processor.ProcessLogs()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	<-sigCh
}

func loadCerts() (string, string, error) {
	utmConfig := plugins.PluginCfg("com.utmstack", false)
	certsFolder, err := utils.MkdirJoin(utmConfig.Get("certsFolder").String())
	if err != nil {
		return "", "", fmt.Errorf("cannot create certificates directory: %v", err)
	}

	certPath := certsFolder.FileJoin(configuration.UtmCertFileName)
	keyPath := certsFolder.FileJoin(configuration.UtmCertFileKey)

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("certificate file does not exist: %s", certPath)
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("key file does not exist: %s", keyPath)
	}

	return certPath, keyPath, nil
}
