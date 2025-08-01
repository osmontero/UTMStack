package configuration

import (
	"fmt"

	"github.com/utmstack/UTMStack/office365/utils"
)

const (
	LoginUrl                  = "https://login.microsoftonline.com/"
	GRANTTYPE                 = "client_credentials"
	SCOPE                     = "https://manage.office.com/.default"
	endPointLogin             = "/oauth2/v2.0/token"
	endPointStartSubscription = "/activity/feed/subscriptions/start"
	endPointContent           = "/activity/feed/subscriptions/content"
	BASEURL                   = "https://manage.office.com/api/v1.0/"
	LogstashEndpoint          = "http://%s:%s"
	UTMLogSeparator           = "<utm-log-separator>"
)

func GetInternalKey() string {
	return utils.Getenv("INTERNAL_KEY")
}

func GetPanelServiceName() string {
	return utils.Getenv("PANEL_SERV_NAME")
}

func GetMicrosoftLoginLink(tenant string) string {
	return fmt.Sprintf("%s%s%s", LoginUrl, tenant, endPointLogin)
}

func GetStartSubscriptionLink(tenant string) string {
	return fmt.Sprintf("%s%s%s", BASEURL, tenant, endPointStartSubscription)
}

func GetContentLink(tenant string) string {
	return fmt.Sprintf("%s%s%s", BASEURL, tenant, endPointContent)
}

func GetLogstashHost() string {
	return utils.Getenv("UTM_LOGSTASH_HOST")
}

func GetLogstashPort() string {
	return utils.Getenv("UTM_LOGSTASH_PORT")
}
