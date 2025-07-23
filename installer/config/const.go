package config

import (
	"path/filepath"
	"time"
)

const INSTALLER_VERSION = "v11.0.0-alpha.1"
const REPLACE = "0aXq879sPfvy2Zrc"
const PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxoZ5m/bsj4NOlSblXPg3
eGbx2nCg5jMoVsCevxX64+MkOpwHBNiA0VPCZfOrdv2atWbOnJ7KfmZJWWZFSlTf
tJ6VA0jaODnSlQeoTn/XIMKkfyzxKgLN+miG89M4ysidZgHPYlwz8+R1gIPouxYa
BUQ8mUgRrkW3JCpYoGvf6k0Od9k8NXdR52rFf0Ryl6oGwedOWh/tiYE0he0pWkB7
zPAveFqvxJnte4aN1Xjv2Qp1OmQVvc37RMZsh8oNpfYTMxrGFWZpBmF61NLXOYjn
YtvKtxUQPR/TOf9p55H9cFeRf2LzSAM4L3NQ3Xdss6eWndjywL5Giqd2gtrD/1Tt
+wIDAQAB
-----END PUBLIC KEY-----
`

const (
	RegisterInstanceEndpoint   = "/api/v1/instances/register"
	GetInstanceDetailsEndpoint = "/api/v1/instances"
	GetUpdatesInfoEndpoint     = "/api/v1/updates"
	SetUpdateSentEndpoint      = "/api/v1/updates/sent"
	GetLicenseEndpoint         = "/api/v1/licenses"
	HealthEndpoint             = "/api/v1/health"
	LogCollectorEndpoint       = "/api/v1/logcollectors/upload"

	ImagesPath = "/utmstack/images"

	CMServer = "https://customermanager.utmstack.com"

	RequiredMinCPUCores  = 2
	RequiredMinDiskSpace = 30
	RequiredDistroUbuntu = "ubuntu"
	RequiredDistroRHEL   = "redhat"
)

var (
	BackendConfigEndpoint  = "https://127.0.0.1/api/utm-configuration-parameters?page=0&size=10000&sectionId.equals=%d&sort=id,asc"
	ConfigPath             = filepath.Join("/root", "utmstack.yml")
	InstanceConfigPath     = filepath.Join(GetConfig().UpdatesFolder, "instance-config.yml")
	ServiceLogPath         = filepath.Join(GetConfig().UpdatesFolder, "logs", "utmstack-updater.log")
	VersionFilePath        = filepath.Join(GetConfig().UpdatesFolder, "version.json")
	LicenseFilePath        = filepath.Join(GetConfig().UpdatesFolder, "LICENSE")
	EventProcessorLogsPath = filepath.Join(GetConfig().DataDir, "events-engine-workdir", "logs")
	CheckUpdatesEvery      = 5 * time.Minute
	SyncSystemLogsEvery    = 5 * time.Minute
	ConnectedToInternet    = false
	Updating               = false
)

func GetCMServer() string {
	cnf := GetConfig()
	return CMServer + "/" + cnf.Branch
}
